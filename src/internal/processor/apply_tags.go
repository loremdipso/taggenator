package processor

import (
	"errors"
	"fmt"
	"internal/database"
	"internal/searcher"

	"github.com/loremdipso/go_utils"
)

func apply_tags(self *QueryProcessor, args []string, db *database.Database) error {
	// TODO: speed up. This takes ~5ms right now for just 1000K entries
	// fmt.Printf("DURATION: %v\n", time.Since(start))

	tags := make([]string, 0)
	for i, arg := range args {
		switch arg {
		case "-tag", "--tag":
			tags = append(tags, args[i+1])

			// TODO: a little messy, but we can handle it
			args[i] = ""
			args[i+1] = ""
			break
		}
	}

	if len(tags) == 0 {
		return errors.New("no tags to add")
	} else {
		fmt.Printf("Tags to add: %s\n", go_utils.StringArrayToString(tags))
	}

	args = go_utils.RemoveEmpty(args)

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
		// TODO: what do here?
		fmt.Printf("\n\n%d / %d\n", i+1, len(entries))

		for _, tag := range tags {
			addTagStringToEntry(db, entry, tag, nil, nil)
		}
		db.UpdateEntry(entry)
	}

	return nil
}
