package processor

import (
	"fmt"
	"internal/database"
	"internal/searcher"
	"path/filepath"

	"github.com/fatih/color"
)

func dump(self *QueryProcessor, args []string, db *database.Database) error {
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
		color.HiRed("No entries")
		return nil
	} else {
		color.HiBlue("Found %d entries", len(entries))
	}

	for _, entry := range entries {
		abspath, _ := filepath.Abs(entry.Location)
		fmt.Println(abspath)
	}

	return nil
}
