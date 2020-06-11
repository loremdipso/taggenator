package processor

import (
	"fmt"
	"internal/database"
	"internal/opener"
	"internal/searcher"
	"strconv"
	"strings"

	"github.com/loremdipso/liner"

	"github.com/fatih/color"
)

func open(self *QueryProcessor, args []string, db *database.Database) error {
	// TODO: speed up. This takes ~5ms right now for just 1000K entries
	// fmt.Printf("DURATION: %v\n", time.Since(start))

	opener := opener.New(db.GetOpenerConfig())
	defer opener.Close()
	// TODO: ugly, but seemingly necessary if we need to do things after ctrl+c
	self.myOpener = opener

	search := searcher.New(db)
	err := search.Parse(args)
	if err != nil {
		return err
	}

	entries, err := search.Execute()
	if err != nil {
		return err
	}

	reader := liner.NewLiner()
	defer reader.Close()
	reader.SetCtrlCAborts(true)
	reader.SetTabCompletionStyle(liner.TabPrints)
	reader.SetBeep(false)

	var sorted_tags []string
	if tags, err := getAllTags(db); err == nil {
		// TODO: is this sorting even necessary?
		// sort.Strings(tags)
		sorted_tags = tags
	} else {
		return err
	}

	completer := NewCompleter(sorted_tags)
	reader.SetCompleter(completer.Complete)
	for i := 0; i < len(entries); {
		if i < 0 {
			i = 0
		}
		entry := entries[i]

		// TODO: what do here?
		fmt.Printf("\n\n%d / %d\n", i+1, len(entries))

		response, err := interactiveAddTags(db, entry, reader, opener, completer.GetSortedTags, completer.Append)

		// TODO: should we only do this sometimes? Or all the time?
		db.UpdateEntry(entry)

		if err != nil {
			return err
		}

		if len(response) > 0 {
			var amount int
			switch response {
			case "q":
				fmt.Println("Quitting...")
				return nil
			default:
				amount = calculateDelta(response)
				i += amount
			}
		} else {
			i++
		}
	}

	return nil
}

func calculateDelta(response string) int {
	var direction int
	if response[0] == '<' {
		direction = -1
		response = response[1:]
	} else if response[0] == '>' {
		direction = 1
		response = response[1:]
	} else if response[len(response)-1] == '<' {
		direction = -1
		response = response[0 : len(response)-1]
	} else if response[len(response)-1] == '>' {
		direction = 1
		response = response[0 : len(response)-1]
	} else {
		return 0
	}

	remainder := strings.Trim(response, " ")
	amount := 1
	var err error
	if len(remainder) > 0 {
		amount, err = strconv.Atoi(remainder)
		if err != nil {
			color.HiRed("Error: %s is an invalid amount", remainder)
			amount = 1
		}
	}
	return amount * direction
}
