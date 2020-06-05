package database

import (
	"fmt"
	. "internal/data"
	"sort"
	"strconv"
	"strings"

	"github.com/loremdipso/go_utils"

	"github.com/fatih/color"
)

const (
	rename_prefix = "rename:"
)

func (db *Database) PreProcessTag(tag string, entry *Entry) (string, bool) {
	if strings.HasPrefix(tag, rename_prefix) {
		color.HiBlue("Renaming...")
		err := db.RenameEntry(entry, tag[len(rename_prefix):])
		if err != nil {
			// TODO: handle this better, maybe?
			color.HiRed("%v", err)
		}
		return "", false
	} else if tag == "tempnew" || tag == "newtemp" {
		tags, err := db.GetAllTags()
		if err != nil {
			// TODO: handle better
			color.HiRed("%v", err)
			return "", false
		}

		return createNewTemp(tags), true
	} else if tag == "temp" {
		tags, err := db.GetAllTags()
		if err != nil {
			// TODO: handle better
			color.HiRed("%v", err)
			return "", false
		}

		return findNewestTemp(tags), true
	} else if tag == "u" || tag == "uns" {
		entry.Times_Opened = 0
		entry.HaveManuallyTouched = false
		fmt.Println("Putting back on the to sort pile")
		return "", false
	} else {
		return db.getReplacementTag(tag)
	}
}

func (db *Database) getReplacementTag(tag string) (string, bool) {
	if newTag, ok := db.settings.Synonyms[tag]; ok {
		return newTag, true
	}

	if command, ok := db.settings.Commands[tag]; ok {
		go_utils.ExecuteCommand(command, false)
		return "", false
	}

	return tag, true
}

func createNewTemp(tags []string) string {
	tag := findNewestTemp(tags)
	numb := getNumbFromTemp(tag)
	numb++
	return fmt.Sprintf("temp%d", numb)
}

func getNumbFromTemp(temp string) int {
	numb, _ := strconv.Atoi(temp[len("temp"):])
	return numb
}

func (self *Database) FindNewestTemp() string {
	tags, err := self.GetAllTags()
	if err != nil {
		// TODO: handle better
		color.HiRed("%v", err)
		return ""
	}

	return findNewestTemp(tags)
}

func findNewestTemp(tags []string) string {
	temps := go_utils.FindWithRegex(tags, "^temp[0-9]*$")
	if len(temps) == 0 {
		return "temp0"
	}

	// TODO: pretty inefficient
	sort.Slice(temps, func(i, j int) bool {
		I := getNumbFromTemp(temps[i])
		J := getNumbFromTemp(temps[j])
		return I < J
	})

	return temps[len(temps)-1]
}
