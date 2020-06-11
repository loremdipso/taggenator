package database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/karrick/godirwalk"
	"github.com/loremdipso/go_utils"
	"github.com/tidwall/buntdb"
)

type Database struct {
	bdb      *buntdb.DB
	settings *Settings

	// true if the database was just created this session
	isNew bool
}

func New(settingsFilename string) (*Database, error) {
	log.Println("initializing database...")

	settings, err := getSettings(settingsFilename)
	if err != nil {
		return nil, err
	}

	return createBuntDB(settings)
}

func createBuntDB(settings *Settings) (*Database, error) {
	var isNew bool
	if !go_utils.FileExists(database_filename) {
		answer, _ := go_utils.ReadChar(fmt.Sprintf("%s doesn't exist. Create it? (y/n)> ", database_filename))
		answer = strings.ToLower(answer)
		if len(answer) == 0 || answer[0] != 'y' {
			return nil, errors.New("Oops, didn't want to make a database")
		}
		isNew = true
	}

	buntDB, err := buntdb.Open(database_filename)
	if err != nil {
		return nil, err
	}

	setupIndexes(buntDB)

	return &Database{buntDB, settings, isNew}, nil
}

func setupIndexes(db *buntdb.DB) {
	//db.CreateIndex("tags", "*", buntdb.IndexString)

	// NOTE: these are actually very slow, since they have to be re-created when loading :(
	// taking out for now
	// Actually, do they stick around? Well, whatever, it doesn't super matter I guess

	// db.CreateIndex(IndexModificationDate, KeyEntries, buntdb.IndexJSON("Modification_Date"))
	// db.CreateIndex(IndexCreationDate, KeyEntries, buntdb.IndexJSON("Creation_Date"))
	// db.CreateIndex(index_tags, id_entry_prefix+"*", buntdb.IndexJSON("Tags"), buntdb.IndexJSON("Tags"))
	// db.CreateIndex(index_creation_date, id_entry_prefix+"*", buntdb.IndexJSON("Tags"), buntdb.IndexJSON("Tags"))

	// db.CreateIndex(index_name, id_entry_prefix+"*", buntdb.IndexJSON("Name"), buntdb.IndexJSON("Times_Opened"))
}

// returns added, modified, removed
func (db *Database) UpdateFiles() (int, int, int, error) {
	addedArr, filePaths, namesToPaths, err := db.addNewFiles()
	if err != nil {
		return len(addedArr), -1, -1, err
	}

	// TODO: if added == 0, do we still need to try and merge?
	// might be a slight optimization
	// db.getFileMaps() // TODO: should this be a reference?
	// filePaths, namesToPaths, err := db.getFileMaps() // TODO: should this be a reference?
	// if err != nil {
	// 	return len(addedArr), 0, 0, err
	// }

	err = findDuplicateFiles(namesToPaths)
	if err != nil {
		return len(addedArr), -1, -1, err
	}

	// _, err = db.DealWithTrimmedFiles(namesToPaths, addedArr) // TODO: this?
	_, err = db.DealWithTrimmedFiles(addedArr) // TODO: this?
	if err != nil {
		return len(addedArr), -1, -1, err
	}

	removed, merged, err := db.RemoveMissingEntries(filePaths, addedArr) // TODO: this?
	if err != nil {
		return len(addedArr), merged, removed, err
	}

	return len(addedArr), merged, removed, nil
}

// // func getAllFileNamesAndPathsMap(db *Database) (map[string][]string, error) {
// func (db *Database) getFileMaps() (map[string]bool, map[string][]string, error) {
// 	filePaths := make(map[string]bool)
// 	namesToPaths := make(map[string][]string)

// 	options := &godirwalk.Options{
// 		Callback: func(path string, dirent *godirwalk.Dirent) error {
// 			if dirent.ModeType().IsRegular() {
// 				if go_utils.ContainsString(db.settings.Extensions, filepath.Ext(path)) {
// 					filePaths[getEntryId(path)] = true

// 					key := filepath.Base(path)
// 					if arr, ok := namesToPaths[key]; ok {
// 						arr = append(arr, path)
// 						namesToPaths[key] = arr
// 					} else {
// 						namesToPaths[key] = []string{path}
// 					}
// 					// if arr, ok := filePaths[pathToName(path)] = path
// 					// filePaths[pathToName(path)] = path
// 				}
// 			}
// 			return nil
// 		},

// 		// faster traversal by not sorting ^_^
// 		Unsorted: true,
// 	}

// 	err := godirwalk.Walk(".", options)
// 	return filePaths, namesToPaths, err
// }

func (db *Database) addNewFiles() ([]string, map[string]bool, map[string][]string, error) {
	jobs := makeJobs(db)
	go generateFilePaths(db, jobs)

	added := make([]string, 0, 100) // TODO: initial size?
	numberOfEntries, _ := db.Len()  // TODO: initial size?
	filePaths := make(map[string]bool, numberOfEntries)
	namesToPaths := make(map[string][]string, numberOfEntries)

	for {
		path, more := <-jobs
		if more {
			didAdd, err := dealWithFilePath(db, path)
			if err != nil {
				// TODO: how do we cancel that coroutine? Do we need to?
				close(jobs)
				return nil, nil, nil, err
			}

			if didAdd {
				if len(added) < 10 {
					fmt.Printf("Adding %s...\n", color.HiBlueString(path))
				} else if len(added) == 10 {
					fmt.Printf("(adding a lot, not listing them all)...\n")
				}
				added = append(added, path)
			}

			if go_utils.ContainsString(db.settings.Extensions, filepath.Ext(path)) {
				filePaths[getEntryId(path)] = true

				key := filepath.Base(path)
				if arr, ok := namesToPaths[key]; ok {
					arr = append(arr, path)
					namesToPaths[key] = arr
				} else {
					namesToPaths[key] = []string{path}
				}
				// if arr, ok := filePaths[pathToName(path)] = path
				// filePaths[pathToName(path)] = path
			}
		} else {
			break
		}
	}

	return added, filePaths, namesToPaths, nil
}

func generateFilePaths(db *Database, jobs chan string) error {
	options := &godirwalk.Options{
		Callback: func(path string, dirent *godirwalk.Dirent) error {
			if dirent.ModeType().IsRegular() {
				jobs <- path
			}
			return nil
		},

		// faster traversal by not sorting ^_^
		Unsorted: true,
	}

	err := godirwalk.Walk(".", options)
	close(jobs)
	return err
}

func dealWithFilePath(db *Database, path string) (bool, error) {
	if !db.ContainsEntry(path) {
		info, err := os.Stat(path)
		if err != nil {
			return false, err
		}

		// TODO: could improve this, do something like git's .gitignore, but
		// I think this is sufficient for now
		if go_utils.ContainsString(db.settings.Extensions, filepath.Ext(path)) {
			return true, db.CreateEntryForFile(path, info)
		}
	}

	// skip
	return false, nil
}

func makeJobs(db *Database) chan string {
	// TODO: make the number of jobs a function of the size of the database,
	// if it's already been created. If we get this right we can save, like, 50%,
	// or ~40ms / 100K files
	min_size := 10000
	max_size := 200000
	size := min_size // some sane default

	if db.isNew {
		// if the db is new,
		size = max_size
	} else {
		len, err := db.Len()
		if err != nil || len == 0 {
			// either there's an error getting the length or the database is new
			// either way
			size = max_size
		} else {
			// NOTE: unclear if this is okay. Might get too big
			size = len
		}
	}

	if size < min_size {
		size = min_size
	} else if size > max_size {
		size = max_size
	}

	return make(chan string, size)
}

func findDuplicateFiles(namesToPaths map[string][]string) error {
	// helpful, but kind of beside the point
	for _, paths := range namesToPaths {
		if len(paths) > 1 {
			return fmt.Errorf("ERROR: duplicate basenames in paths: %v", paths)
			// fmt.Printf("WARNING: duplicate basenames in paths: %v", paths)
		}
	}

	return nil
}

// NOTE: this will necessarily be slow
// func getAllFilePathsAsEntryIdMap(db *Database) (map[string]bool, error) {
// 	filePaths := make(map[string]bool)
// 	options := &godirwalk.Options{
// 		Callback: func(filepath string, dirent *godirwalk.Dirent) error {
// 			if dirent.ModeType().IsRegular() {
// 				filePaths[getEntryId(filepath)] = true
// 			}
// 			return nil
// 		},

// 		// faster traversal by not sorting ^_^
// 		Unsorted: true,
// 	}

// 	err := godirwalk.Walk(".", options)
// 	return filePaths, err
// }

// // func getAllFileNamesAndPathsMap(db *Database) (map[string][]string, error) {
// func (db *Database) getFileMaps() (map[string]bool, map[string][]string, error) {
// 	filePaths := make(map[string]bool)
// 	namesToPaths := make(map[string][]string)

// 	options := &godirwalk.Options{
// 		Callback: func(path string, dirent *godirwalk.Dirent) error {
// 			if dirent.ModeType().IsRegular() {
// 				if go_utils.ContainsString(db.settings.Extensions, filepath.Ext(path)) {
// 					filePaths[getEntryId(path)] = true

// 					key := filepath.Base(path)
// 					if arr, ok := namesToPaths[key]; ok {
// 						arr = append(arr, path)
// 						namesToPaths[key] = arr
// 					} else {
// 						namesToPaths[key] = []string{path}
// 					}
// 					// if arr, ok := filePaths[pathToName(path)] = path
// 					// filePaths[pathToName(path)] = path
// 				}
// 			}
// 			return nil
// 		},

// 		// faster traversal by not sorting ^_^
// 		Unsorted: true,
// 	}

// 	err := godirwalk.Walk(".", options)
// 	return filePaths, namesToPaths, err
// }

// func mapFilePaths(db *Database) (map[string]bool, error) {
// 	filePaths := make(map[string]bool)
// 	options := &godirwalk.Options{
// 		Callback: func(path string, dirent *godirwalk.Dirent) error {
// 			if dirent.ModeType().IsRegular() {
// 				filePaths[getEntryId(path)] = true
// 			}
// 			return nil
// 		},

// 		// faster traversal by not sorting ^_^
// 		Unsorted: true,
// 	}

// 	err := godirwalk.Walk(".", options)
// 	return filePaths, err
// }
