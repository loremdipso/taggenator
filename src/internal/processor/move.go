package processor

import (
	"errors"
	"fmt"
	"internal/database"
	"internal/searcher"
	"os"
	"path/filepath"

	"github.com/loremdipso/go_utils"
)

func move(self *QueryProcessor, args []string, db *database.Database) error {
	var destination string

	for i, arg := range args {
		switch arg {
		case "destination", "-destination", "--destination":
			destination = args[i+1]

			// TODO: a little messy, but we can handle it
			args[i] = ""
			args[i+1] = ""
			break
		}
	}

	if destination == "" {
		return errors.New("ERROR: need to specify -destination")
	}

	search := searcher.New(db)
	err := search.Parse(args)
	if err != nil {
		return err
	}

	entries, err := search.Execute()
	if err != nil {
		return err
	}

	// TODO: make sure this is valid
	err = os.MkdirAll(destination, os.ModePerm)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		oldLocation := entry.Location
		newLocation := filepath.Join(destination, filepath.Base(oldLocation))

		if !go_utils.FileExists(oldLocation) {
			fmt.Printf("%s already missing\n", oldLocation)
		} else {
			if go_utils.FileExists(newLocation) {
				fmt.Println("yo")
				return errors.New(fmt.Sprintf("ERROR: %s already exists", newLocation))
			}

			// TODO: support other disks
			fmt.Printf("Moving %s => %s\n", oldLocation, newLocation)
			err := go_utils.MoveFile(oldLocation, newLocation)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}

	return nil
}
