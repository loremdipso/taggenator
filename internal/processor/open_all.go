package processor

import (
	"fmt"
	"internal/data"
	"internal/database"
	"internal/searcher"
	"io/ioutil"
	"path/filepath"

	"github.com/loremdipso/go_utils"

	"github.com/fatih/color"
)

func open_all(self *QueryProcessor, args []string, db *database.Database) error {
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

	switch filepath.Ext(entries[0].Location) {
	case ".jpg", ".png", ".jpeg":
		return openAllImages(entries)
	case ".cbz":
		return openAllComics(entries)
	case ".mp4", ".wmv", ".avi", ".flv": // TODO: refactor
		return openAllVideos(entries)
	default:
		return openAllGeneric(entries)
	}
}

func openAllImages(entries data.Entries) error {
	// TODO: make generic. Assumes geeqie
	var entryString = ""
	for _, entry := range entries {
		entryString += fmt.Sprintf(" \"%s\"", entry.Location)
	}

	go_utils.ExecuteCommand(fmt.Sprintf("geeqie -r %s", entryString), false)

	return nil
}

func openAllComics(entries data.Entries) error {
	// TODO: make generic. Assumes mcomix
	var entryString = ""
	for _, entry := range entries {
		entryString += fmt.Sprintf(" \"%s\"", entry.Location)
	}

	go_utils.ExecuteCommand(fmt.Sprintf("mcomix %s", entryString), false)

	return nil
}

func openAllGeneric(entries data.Entries) error {
	// TODO: make generic. Assumes mcomix
	var entryString = ""
	for _, entry := range entries {
		entryString += fmt.Sprintf(" \"%s\"", entry.Location)
	}

	go_utils.ExecuteCommand(fmt.Sprintf("gnome-open %s", entryString), false)

	return nil
}

func openAllVideos(entries data.Entries) error {
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
