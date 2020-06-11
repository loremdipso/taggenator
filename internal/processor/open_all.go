package processor

import (
	"fmt"
	"internal/database"
	"internal/searcher"
	"io/ioutil"
	"path/filepath"

	"github.com/loremdipso/go_utils"

	"github.com/fatih/color"
)

func open_all(self *QueryProcessor, args []string, db *database.Database) error {
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

	tempfile, err := ioutil.TempFile("/tmp", "taggy_")
	defer tempfile.Close()
	if err != nil {
		return err
	}

	fmt.Fprintln(tempfile, "#EXTM3U")
	for _, entry := range entries {
		abspath, _ := filepath.Abs(entry.Location)
		fmt.Fprintln(tempfile, abspath)
	}

	go_utils.ExecuteCommand(fmt.Sprintf("vlc \"%s\"", tempfile.Name()), false)

	return nil
}
