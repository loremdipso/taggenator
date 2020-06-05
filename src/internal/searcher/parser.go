package searcher

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/loremdipso/go_utils"

	"github.com/fatih/color"
)

func (s *Searcher) Parse(args []string) error {
	// NOTE: can we just not do this? Or is it necessary?
	s.AppendSimpleFilter(FilterAny)

	sortArgs := make([]string, 0, 10) // TODO: how long should this be
	for _, arg := range args {
		if len(arg) > 0 {
			arg := strings.ToLower(arg)
			if arg == "-sort" {
				err := s.handleSortArgs(sortArgs)
				if err != nil {
					fmt.Println("Error: ", err)
					return err
				}

				// NOTE: turns out this wasn't the right way :)
				// sortArgs = sortArgs[:0] // TODO: is this the right way to clear the array?
				sortArgs = make([]string, 0, 10) // TODO: how long should this be?
			} else {
				sortArgs = append(sortArgs, s.fixArg(arg))
			}
		}
	}

	return s.handleSortArgs(sortArgs)
}

func (s *Searcher) fixArg(arg string) string {
	// Special case: open up newest arg
	if arg == "tempnew" || arg == "newtemp" {
		newTag := s.db.FindNewestTemp()
		fmt.Printf("Swapping %s out for %s\n", color.HiBlueString(arg), color.HiBlueString(newTag))
		return newTag
	}
	return arg
}

func (s *Searcher) handleSortArgs(args []string) error {
	if len(args) == 0 {
		return nil
	}

	arg := args[0]

	remainder := make([]string, 0)
	if len(args) > 1 {
		remainder = args[1:]
	}

	arg, remainder = splitArgIntoRemainder(arg, remainder)

	switch arg {
	// sorters
	case "newest":
		s.AppendSimpleSorter(SortNewestAdded)
	case "oldest":
		s.AppendSimpleSorter(SortOldestAdded)
	case "largest", "biggest":
		s.AppendSimpleSorter(SortLargest)
	case "smallest":
		s.AppendSimpleSorter(SortSmallest)
	case "reverse":
		s.AppendCustomSorter(SortReverse)
	case "random":
		s.AppendCustomSorter(SortRandom)
	case "alpha", "alphabetical":
		s.AppendSimpleSorter(SortName)

	case "most_recently_opened":
		s.AppendSimpleSorter(SortMostRecentlyOpened)
	case "least_recently_opened":
		s.AppendSimpleSorter(SortLeastRecentlyOpened)

	// filters
	case "seen":
		s.AppendSimpleFilter(FilterSeen)

	case "unseen":
		s.AppendSimpleFilter(FilterUnseen)

	case "touched":
		s.AppendSimpleFilter(FilterTouched)
	case "untouched":
		s.AppendSimpleFilter(FilterUnTouched)

	case "name_inclusive", "name_includes":
		s.AppendSimpleFilter(FilterNameInclusive(remainder))

	case "name", "name_exclusive":
		s.AppendSimpleFilter(FilterNameInclusive(remainder))

	case "path_includes", "path_inclusive":
		s.AppendSimpleFilter(FilterPathInclusive(remainder))

	case "path", "path_exclusive":
		s.AppendSimpleFilter(FilterPathExclusive(remainder))

	case "tags_exclusive":
		s.AppendSimpleFilter(FilterTagsExclusive(remainder))

	case "tags_inclusive":
		s.AppendSimpleFilter(FilterTagsInclusive(remainder))

	case "search":
		s.AppendSimpleFilter(FilterSearchExclusive(remainder))

	case "search_inclusive":
		s.AppendSimpleFilter(FilterSearchInclusive(remainder))

	default:
		max, err := strconv.Atoi(arg)
		if err != nil {
			return errors.New(fmt.Sprintf("%s is an invalid search option", arg))
		} else {
			s.AppendCustomFilter(FilterMaxCapacity(max))
		}
	}

	return nil
}

// Special case: this argument might actually be ':' (or ",") delimited, so in that case we'll use that
// NOTE: kind of slow, but it's so small we don't really care
// Also, mostly kept for hostoric reasons
func splitArgIntoRemainder(arg string, remainder []string) (string, []string) {
	// supports : and , as delims. Whichever comes first is the one we'll split on
	delims := []string{":", ","}
	var firstDelim string
	var firstDelimIndex int = len(arg)
	for _, delim := range delims {
		index := strings.Index(arg, delim)
		if index > -1 && index < firstDelimIndex {
			firstDelim = delim
			firstDelimIndex = index
		}
	}

	if firstDelim != "" {
		pieces := strings.Split(arg, firstDelim)
		arg = pieces[0]
		if len(pieces) > 1 {
			remainder = append(pieces[1:], remainder...)
		}
	}

	return arg, go_utils.RemoveEmpty(remainder)
}
