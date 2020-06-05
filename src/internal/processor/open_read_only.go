package processor

import (
	"fmt"
	"internal/database"
	"internal/opener"
	"internal/searcher"

	"github.com/loremdipso/go_utils"

	"github.com/loremdipso/fancy_printer"

	"github.com/fatih/color"
)

func open_read_only(self *QueryProcessor, args []string, db *database.Database) error {
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

	for i, entry := range entries {
		opener.Open(entry.Location)

		fmt.Printf("\n\n\n")
		fmt.Printf("\n\n%d / %d\n", i+1, len(entries))

		prefix := "Trying "
		if truncated_line, postfix, err := fancy_printer.GetTruncatedLine(prefix, entry.Location); err != nil {
			return err
		} else {
			fmt.Printf("%s%s%s\n", prefix, color.HiRedString("%s", truncated_line), postfix)
		}

		fmt.Printf("Times opened: %d\n", entry.Times_Opened)

		fancy_printer.PrintArrayAsGrid(entry.Tags, false, true)
		fmt.Printf("Size: %s\n", go_utils.HumanReadableSize(entry.Size))
		go_utils.Readline("<Read Only. Enter to Advance>")
	}

	return nil
}
