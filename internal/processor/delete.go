package processor

import (
	"fmt"
	"internal/database"
	"internal/searcher"
	"os"

	"github.com/loremdipso/go_utils"
)

func delete(self *QueryProcessor, args []string, db *database.Database) error {
	search := searcher.New(db)
	err := search.Parse(args)
	if err != nil {
		return err
	}

	entries, err := search.Execute()
	if err != nil {
		return err
	}

	fmt.Println(entries)
	fmt.Println(args)

	for _, entry := range entries {
		location := entry.Location

		if !go_utils.FileExists(location) {
			fmt.Printf("ERROR: can't delete, %s since it doesn't exist\n", location)
			return fmt.Errorf("ERROR: can't delete, %s since it doesn't exist\n", location)
		} else {
			// TODO: also remove from db, or just wait until next startup?
			fmt.Printf("Deleting %s...\n", location)
			err := os.Remove(location)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}

	return nil
}
