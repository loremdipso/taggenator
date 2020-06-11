package processor

import (
	"fmt"
	"internal/database"
	"internal/searcher"

	"github.com/loremdipso/go_utils"

	"github.com/fatih/color"
)

func dump_tags(self *QueryProcessor, args []string, db *database.Database) error {
	if len(args) == 0 {
		// Special Case: no search arguments, just print all tags
		tags, err := getAllTags(db)
		if err != nil {
			return err
		}

		// fmt.Println("Tag Dump:", go_utils.StringArrayToString(tags))
		for _, tag := range tags {
			fmt.Println(tag)
		}
	} else {
		// NOTE: this assumes vlc
		// TODO: make more generic
		search := searcher.New(db)
		err := search.Parse(args)
		if err != nil {
			return err
		}

		entries, err := search.Execute()
		if err != nil {
			return err
		}

		if len(entries) == 0 {
			color.HiRed("# No entries")
			return nil
		} else {
			fmt.Printf("# Found %s entries\n", color.HiBlueString("%d", len(entries)))
		}

		sortedTags := make([]string, 0, 1000) // TODO: length?
		for _, entry := range entries {
			for _, tag := range entry.Tags {
				_, sortedTags = go_utils.InsertIntoSortedListIfNotThereAlready(sortedTags, tag)
			}
		}

		for _, tag := range sortedTags {
			fmt.Println(tag)
		}
	}

	return nil
}
