package processor

import (
	"errors"
	"fmt"
	"internal/data"
	"internal/database"
	"internal/opener"

	"internal/searcher"

	"github.com/loremdipso/go_utils"

	"github.com/fatih/color"
)

type actionStruct struct {
	name        string
	synonyms    []string
	compute     func(*QueryProcessor, []string, *database.Database) error
	description descriptionStruct
}

type descriptionStruct struct {
	text  string
	color func(string, ...interface{}) string
}

var actions = []actionStruct{
	{"help", []string{"help", "?", "h", "-help", "-h"}, TODO, descriptionStruct{"[action] get help", color.HiBlueString}},
	{"open", nil, open, descriptionStruct{"", color.HiBlueString}},
	{"open_read_only", nil, open_read_only, descriptionStruct{"Open read only", color.HiBlueString}},
	{"open_all", nil, open_all, descriptionStruct{"Open all", color.HiBlueString}},
	{"apply_tags", nil, apply_tags, descriptionStruct{"Apply tags (--tag) to the search results. Optionally can use --threads \"#\" to spawn extra workers", color.HiBlueString}},
	{"move", nil, move, descriptionStruct{"move results to -destination", color.HiYellowString}},
	{"delete", nil, delete, descriptionStruct{"delete results", color.HiYellowString}},
	{"dump_tags", nil, dump_tags, descriptionStruct{"Dump all tags", color.HiBlueString}},
	{"dump", nil, dump, descriptionStruct{"Dump paths to all entries", color.HiBlueString}},
	{"fix", nil, fix, descriptionStruct{"try and fix innacuracies in database", color.HiGreenString}},

	// TODO: remove
	{"combine", nil, absorb_old_database, descriptionStruct{"[filename] Combine old database", color.HiBlueString}},
}

type QueryProcessor struct {
	myOpener *opener.Opener
}

// New create a new opener
func New() *QueryProcessor {
	return &QueryProcessor{}
}

// Close closes the opener
func (self *QueryProcessor) Close() {
	if self.myOpener != nil {
		self.myOpener.Close()
	}
}

func (self *QueryProcessor) ProcessQuery(args []string, db *database.Database, projectName string) error {
	var arg string
	if len(args) > 0 {
		arg, args = args[0], args[1:]
	}
	switch arg {
	case "-h", "help", "":
		return help(args, projectName)
	default:
		var foundAction *actionStruct
		for _, action := range actions {
			if arg == action.name {
				foundAction = &action
				break
			}
			for _, synonym := range action.synonyms {
				if arg == synonym {
					foundAction = &action
					break
				}
			}
		}
		if foundAction == nil {
			fmt.Printf("Error: argument %s is invalid\n", color.HiRedString(arg))
			return help(nil, projectName)
		} else {
			return foundAction.compute(self, args, db)
		}
	}
	/* Query Format:
	-h {options}
	*/
	// err := db.View(func(tx *buntdb.Tx) error {
	// 	err := tx.Ascend("", func(key, value string) bool {
	// 		fmt.Printf("key: %s, value: %s\n", key, value)
	// 	})
	// 	return err
	// })
}

//#region Help
func help(args []string, projectName string) error {
	if len(args) > 0 {
		TODO(nil, nil, nil)
	} else {
		names := GetActionNames(actions)
		prefixSpacing := 4
		longestKey := go_utils.FindLongest(names) + prefixSpacing
		color.HiGreen("%s %s\n", projectName, go_utils.Join(names, " | "))
		for index, key := range names {
			action := actions[index]
			fmt.Printf("%*s", longestKey, key)
			fmt.Printf(": %s\n", action.description.color("%s", action.description.text))
		}
	}

	return nil
}

func GetActionNames(actions []actionStruct) []string {
	var names []string
	for _, action := range actions {
		names = append(names, action.name)
	}
	return names
}

//#endregion Help

func TODO(*QueryProcessor, []string, *database.Database) error {
	return errors.New("ERROR: not yet implemented")
}

func getAllEntries(db *database.Database) (data.Entries, error) {
	search := searcher.New(db)
	search.AppendSimpleFilter(searcher.FilterAny)
	entries, err := search.Execute()
	// tags, err := searcher.BatchSearchAndFilter(db, &data.SearchAndFilter{SimpleSearch: searcher.Any})
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func getAllTags(db *database.Database) ([]string, error) {
	// TODO: redo performance testing once we have more tags
	return getAllTagsOptimized(db)
}

func getAllTagsOptimized(db *database.Database) ([]string, error) {
	return db.GetAllTags()
}

func getAllTagsNaive(db *database.Database) ([]string, error) {
	search := searcher.New(db)
	search.AppendSimpleFilter(searcher.FilterAny)
	entries, err := search.Execute()
	if err != nil {
		return nil, err
	}
	return getTagsForEntries(entries), nil
}

// TODO: maybe this should be a map?
func getTagsForEntries(entries data.Entries) []string {
	mapping := make(map[string]bool)
	for _, entry := range entries {
		for _, tag := range entry.Tags {
			mapping[tag] = true
		}
	}
	return go_utils.MapToArray(mapping)
}
