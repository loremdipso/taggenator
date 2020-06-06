package processor

import (
	"fmt"
	"internal/data"
	"internal/database"
	"internal/opener"
	"sort"
	"time"

	"github.com/loremdipso/liner"

	"github.com/loremdipso/fancy_printer"

	"strings"

	"github.com/loremdipso/go_utils"

	"github.com/fatih/color"
)

func interactiveAddTags(db *database.Database, entry *data.Entry, reader *liner.State, opener *opener.Opener, getSortedTags func() []string, onTagAdd func(string) bool) (string, error) {
	// var wasUpdated bool
	entry.Date_Last_Accessed = time.Now()
	entry.Times_Opened += 1

	// defer func() {
	// 	// Special case: we're returning, and we haven't added/removed any tags
	// 	// Update the entry before moving on, so at least the Time_Opened count is right
	// 	if !wasUpdated {
	// 		db.UpdateEntry(entry)
	// 	}
	// }()
	opener.Open(entry.Location)
	firstOpening := true

	var lastResponse []string
	for {
		prefix := "Trying "
		if truncated_line, postfix, err := fancy_printer.GetTruncatedLine(prefix, entry.Location); err != nil {
			return "", err
		} else {
			fmt.Printf("%s%s%s\n", prefix, color.HiRedString("%s", truncated_line), postfix)
		}

		fmt.Printf("Times opened: %d\n", entry.Times_Opened)

		// TODO: add back, maybe?
		if firstOpening && entry.Times_Opened == 1 {
			firstOpening = false
			// if entry.Last_Auto_Updated
			lastResponse = addAutoTagsToEntry(db, entry, getSortedTags, onTagAdd)
		}

		fancy_printer.PrintArrayAsGrid(entry.Tags, false, true)
		fmt.Printf("Size: %s\n", go_utils.HumanReadableSize(entry.Size))

		// response, err := go_utils.Readline("") // simple version
		var response string
		fmt.Printf("Tags for %s?\n", color.HiGreenString("%s", entry.Name))
		// prompt := fmt.Sprintf("Tags for %s:", color.HiGreenString("%s", entry.Name))
		prompt := fmt.Sprintf("Tags: ")
		response, err := reader.Prompt(prompt)
		if err != nil {
			reader.Close()
			fmt.Println("ERROR", err)
			return "q", err
		}
		fmt.Printf("\n\n\n")

		response = strings.ToLower(response)
		if len(response) > 0 {
			if response == "q" {
				return response, nil
			}

			switch response[0] {
			case '<', '>':
				return response, nil
			}

			switch response[len(response)-1] {
			case '<', '>':
				return response, nil
			}
		} else {
			return "", nil
		}

		var numAdded, numRemoved int
		numAdded, numRemoved, lastResponse = addTagStringToEntry(db, entry, response, lastResponse, onTagAdd)

		if response != "-" {
			if numAdded > 0 || numRemoved > 0 {
				// attempt to avoid immediate duplicates
				if go_utils.LastElement(reader.GetHistory()) != response {
					reader.AppendHistory(response)
				}
			}

			entry.HaveManuallyTouched = true
		}
	}
}

// TODO: use this for helping out with readline
// unless... maybe GO already has such a library?
func findSubstringMatches(arr []string, searchStr string) []string {
	var matches []string = make([]string, 0) // TODO: capacity of matches?
	for _, tag := range arr {
		if strings.Contains(tag, searchStr) {
			matches = append(matches, tag)
		}
	}
	return matches
}

// Returns num added, num removed, previous tags
// TODO: refactor. This is kind of complicated now
func addTagStringToEntry(db *database.Database, entry *data.Entry, s_tags string, previousTags []string, onTagAdd func(string) bool) (int, int, []string) {
	// NOTE: assumes lowercase

	// Special Case: if the response is just '-', flip all previous tags.
	// That is, remove those we added and add those we removed
	var tags []string
	if s_tags == "-" {
		tags = previousTags
		for i, tag := range tags {
			tag, shouldRemove := shouldRemoveTag(tag)
			if shouldRemove {
				tags[i] = tag
			} else {
				tags[i] = tag + "-"
			}
		}
	} else {
		tags = stringToTags(s_tags)
		tags = preProcessTags(db, entry, tags)
	}

	numAdded, numRemoved := addTagsToEntry(db, entry, tags, onTagAdd)
	return numAdded, numRemoved, tags
}

func stringToTags(s_tags string) []string {
	tags := strings.Split(s_tags, ",")
	for i, tag := range tags {
		tags[i] = strings.Trim(tag, " ")
	}
	return tags
}

func preProcessTags(db *database.Database, entry *data.Entry, tags []string) []string {
	originalTagLength := len(tags) // since we'll be appending to tags as we go along
	//for i, tag := range tags {
	for i := 0; i < originalTagLength; i++ {
		tag := tags[i]
		tag, shouldRemove := shouldRemoveTag(tag)
		newTag, extraTags := db.PreProcessTag(tag, entry)
		if len(newTag) > 0 {
			if shouldRemove {
				tags[i] = newTag + "-"
			} else {
				tags[i] = newTag
			}
		} else {
			tags[i] = "" // not very clever, but it'll do
			if len(extraTags) > 0 {
				tags = append(tags, extraTags...)
			}
		}
	}
	return tags
}

func addTagsToEntry(db *database.Database, entry *data.Entry, tags []string, onTagAdd func(string) bool) (int, int) {
	var numAdded, numRemoved int

	// NOTE: this is kind of slow, but it's
	sort.Strings(entry.Tags)

	for _, newTag := range tags {
		if len(newTag) > 0 {
			newTag, shouldRemove := shouldRemoveTag(newTag)

			if shouldRemove {
				fmt.Println("Removing Tag: ", color.HiBlueString("%s", newTag))
				didRemove, newArr := go_utils.RemoveStringArrayElementIfExists(entry.Tags, newTag)
				if didRemove {
					numRemoved++
					entry.Tags = newArr
				}
			} else {
				fmt.Println("Adding Tag: ", color.HiBlueString("%s", newTag))
				didAdd, newArr := go_utils.InsertIntoSortedListIfNotThereAlready(entry.Tags, newTag)
				if didAdd {
					numAdded++
					if onTagAdd != nil {
						onTagAdd(newTag)
					}
					entry.Tags = newArr
				}
			}
		}
	}

	return numAdded, numRemoved
}

func shouldRemoveTag(tag string) (string, bool) {
	if tag[len(tag)-1] == '-' {
		tag = tag[:len(tag)-1]
		return tag, true
	} else if tag[0] == '-' {
		tag = tag[1:]
		return tag, true
	}
	return tag, false
}
