package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	. "internal/data"

	"github.com/loremdipso/go_utils"

	"github.com/fatih/color"
	"github.com/tidwall/buntdb"
)

func (db *Database) Close() {
	log.Println("Closing database...")
	// NOTE: this takes up some time, so maybe don't do it always?
	// Or, maybe do it synchronously?
	db.bdb.Close()
}

func (db *Database) Shrink() {
	db.bdb.Shrink()
}

func (db *Database) GetEntry(s_entryId string) (*Entry, error) {
	value, err := db.getEntryString(s_entryId)
	if err != nil {
		return nil, err
	}

	return deserializeEntry(value)
}

func (db *Database) getEntryString(s_entryId string) (string, error) {
	value, err := db.Get(getEntryId(s_entryId))
	if err != nil {
		return "", err
	}

	return value, nil
}

func (db *Database) RenameEntry(entry *Entry, newfilename string) error {
	oldLocation := entry.Location
	newName := pathToName(newfilename)
	newLocation := filepath.Join(filepath.Dir(oldLocation), newfilename+filepath.Ext(oldLocation))

	if go_utils.FileExists(newLocation) {
		return errors.New(fmt.Sprintf("ERROR: %s already exists", newLocation))
	}

	entry.Name = newName
	entry.Location = newLocation
	serializedEntry, err := serializeEntry(entry)
	if err != nil {
		return err
	}

	err = os.Rename(oldLocation, newLocation)
	if err != nil {
		return err
	}

	db.Remove(getEntryId(oldLocation))
	// TODO: is this right? Shouldn't it be the path?
	// return db.Add(getTagId(pathToName(newLocation)), string(serializedEntry))
	return db.Add(getEntryId(newLocation), string(serializedEntry))
}

func (db *Database) UpdateEntry(value *Entry) error {
	serializedEntry, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// NOTE: IMPORTANT! We're treating the path as a kind of UID. Might be slow.
	values := make([]string, 1+len(value.Tags))
	keys := make([]string, 1+len(value.Tags))
	keys[0] = getEntryId(value.Location)
	values[0] = string(serializedEntry)
	for i := 0; i < len(value.Tags); i++ {
		keys[i+1] = getTagId(value.Tags[i])
		values[i+1] = value.Tags[i]
	}
	return db.AddMultiple(keys, values)
}

func (db *Database) CreateEntryForFile(path string, info os.FileInfo) error {
	// TODO: make more generic
	// Special Case: any files that begin with 0trimmedx will be initialized with a copy of
	// the base file, if it exists

	// TODO: implement this
	entry := Entry{
		Name:     pathToName(path),
		Location: path,
		Size:     info.Size(),

		Date_Created: info.ModTime(),
		Date_Added:   time.Now(),
		Times_Opened: 0,
		// Tags:     []string{filepath},
	}

	return db.UpdateEntry(&entry)
}

func (db *Database) FixAllEntries() error {
	return db.bdb.Update(func(tx *buntdb.Tx) error {
		// Need to do this since we can't read/write at the same time. Sign
		keys := make(map[string]string, 100) // TODO: what size should this be? Does it matter?
		tx.AscendKeys(KeyEntries, func(entryId, value string) bool {
			keys[entryId] = value
			return true
		})

		for key, value := range keys {
			originalEntry, err := deserializeEntry(value)
			if err != nil {
				return err
			}

			// TODO: do we even want this behavior?
			originalEntry.Size = go_utils.GetFileSize(getEntryKeyFromId(key))
			sort.Strings(originalEntry.Tags)

			newValue, err := serializeEntry(originalEntry)
			if err != nil {
				return err
			}

			_, _, err = tx.Set(key, string(newValue), nil)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (db *Database) DealWithTrimmedFiles(addedPaths []string) (int, error) {
	// TODO: make more generic
	// Special Case: any files that begin with 0trimmedx will be initialized with a copy of
	// the base file, if it exists
	var updated int
	trimmedPrefix := "0trimmedx "

	// namesToEntryIDs := make(map[string]string, len(addedPaths))
	err := db.bdb.Update(func(tx *buntdb.Tx) error {
		namesToPaths := make(map[string]string, 100) // TODO: what size should this be? Does it matter?

		tx.AscendKeys(KeyEntries, func(entryId, value string) bool {
			path := getEntryKeyFromId(entryId)
			namesToPaths[filepath.Base(path)] = path
			return true
		})

		for _, addedPath := range addedPaths {
			if strings.Contains(pathToName(addedPath), trimmedPrefix) {
				// TODO: this will fail if there're multiple matches for 0trimmedx
				queryName := filepath.Base(addedPath)[len(trimmedPrefix):]

				// also need to remove the end portion, which is a space followed by non-spaces followed by the extension
				pieces := strings.Split(queryName, " ")

				// TODO: make this support multiple extensions
				queryName = strings.Join(pieces[:len(pieces)-1], " ") + filepath.Ext(addedPath)

				queryPaths, ok := namesToPaths[queryName]
				if ok {
					queryPath := queryPaths
					queryEntryId := getEntryId(queryPath)

					value, err := tx.Get(queryEntryId)
					if err != nil {
						return err
					}

					originalEntry, err := deserializeEntry(value)
					if err != nil {
						return err
					}

					// TODO: do we even want this behavior?
					originalEntry.Name = pathToName(addedPath)
					originalEntry.Location = addedPath
					originalEntry.Size = go_utils.GetFileSize(addedPath)
					// originalEntry.Date_Added = time.Now()
					originalEntry.Times_Opened = 0

					newValue, err := serializeEntry(originalEntry)
					if err != nil {
						return err
					}

					fmt.Printf("Adding trimmed file: %s\n", color.HiMagentaString(addedPath))
					_, _, err = tx.Set(getEntryId(addedPath), string(newValue), nil)
					if err != nil {
						return err
					}

					updated++
				}
			}
		}
		return nil
	})
	return updated, err

}

// func (db *Database) CreateEntry(value *Entry) (string, error) {
// 	serializedEntry, err := json.Marshal(value)
// 	if err != nil {
// 		return "", err
// 	}

// 	// TODO: should we check for duplicates? Or is the location being the key enough?
// 	return db.Add(value.Location, string(serializedEntry))
// 	// return db.appendEntry(string(serializedEntry))
// }

// func (db *Database) appendEntry(value string) (string, error) {
// 	var s_entryId string
// 	err := db.bdb.Update(func(tx *buntdb.Tx) error {
// 		val, err := tx.Get(EntryID_Key)
// 		var entryIdValue int = 0
// 		if err == nil {
// 			separatorPos := strings.Index(val, id_separator)
// 			if separatorPos > -1 {
// 				entryIdValue, err = strconv.Atoi(val[(separatorPos + 1):])
// 				if err != nil {
// 					return err
// 				} else {
// 					entryIdValue++
// 				}
// 			}
// 		}

// 		s_entryId = id_entry_prefix + strconv.Itoa(entryIdValue)
// 		_, _, err = tx.Set(s_entryId, value, nil)
// 		if err != nil {
// 			return err
// 		}

// 		_, _, err = tx.Set(EntryID_Key, s_entryId, nil)
// 		return err
// 	})
// 	if err != nil {
// 		log.Println("Error adding entry: ", err)
// 	}

// 	return s_entryId, err
// }

func (db *Database) AddMultiple(keys, values []string) error {
	return db.SetMultiple(keys, values)
}

func (db *Database) Add(key, value string) error {
	return db.Set(key, value)
}

func (db *Database) Len() (int, error) {
	var count int
	err := db.bdb.View(func(tx *buntdb.Tx) error {
		var err error
		count, err = tx.Len()
		return err
	})
	return count, err
}

func (db *Database) SetMultiple(keys, values []string) error {
	if len(keys) != len(values) {
		return errors.New("Must have same length keys and values")
	}

	return db.bdb.Update(func(tx *buntdb.Tx) error {
		for i, key := range keys {
			_, _, err := tx.Set(key, values[i], nil)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (db *Database) Set(key, value string) error {
	err := db.bdb.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, value, nil)
		return err
	})
	return err
}

func (db *Database) ContainsEntry(s_entryId string) bool {
	value, err := db.getEntryString(s_entryId)
	return err == nil && len(value) > 0
}

func (db *Database) GetAllTags() ([]string, error) {
	// TODO: create some index, maybe?
	tags := make([]string, 0, 0) // TODO: make initial size the same as the number of tags
	err := db.bdb.View(func(tx *buntdb.Tx) error {
		var err error
		tx.AscendKeys(KeyTags, func(key, _ string) bool {
			// TODO: right now we could just store the actual key in the value and skip
			// this string op. But I'm thinking we might store other things in the tag,
			// like synonyms or w/e
			tag := key[len(id_tag_prefix):]
			tags = append(tags, tag)
			return true
		})
		return err
	})

	return tags, err
}

func (db *Database) Get(key string) (string, error) {
	var value string
	err := db.bdb.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(key)
		if err != nil {
			return err
		}
		value = val
		return nil
	})
	return value, err
}

// Removes entries not in the entry list. Returns number removed, number merged (renamed), err
func (db *Database) RemoveMissingTags() (int, error) {
	var removed int

	allTags, err := db.GetAllTags() // NOTE: should be sorted
	if err != nil {
		return -1, err
	}

	tagWhitelist := make(map[string]bool, 1000) // TODO: initial size?
	err = db.bdb.View(func(tx *buntdb.Tx) error {
		tx.AscendKeys(KeyEntries, func(key, value string) bool {
			var entry *Entry
			entry, err = deserializeEntry(value)
			if err != nil {
				return false
			}

			for _, tag := range entry.Tags {
				tagWhitelist[tag] = true
			}
			return true
		})

		return err
	})

	if err != nil {
		return -1, err
	}

	err = db.bdb.Update(func(tx *buntdb.Tx) error {
		for _, tag := range allTags {
			// if missing from whitelist, remove
			if _, ok := tagWhitelist[tag]; !ok {
				removed++

				if removed < 10 {
					log.Printf("# Removing unused tag %s...\n", color.HiRedString("%s", tag))
				} else if removed == 10 {
					log.Printf("# (removing a lot of tags, not listing them here)...\n")
				}

				_, err := tx.Delete(getTagId(tag))
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	return removed, err
}

// Removes entries not in the entry list. Returns number removed, number merged (renamed), err
func (db *Database) RemoveMissingEntries(entryIdMapWhitelist map[string]bool, addedPaths []string) (int, int, error) {
	var removed, merged int

	namesToEntryIDs := make(map[string]string, len(addedPaths))
	for _, path := range addedPaths {
		namesToEntryIDs[filepath.Base(path)] = getEntryId(path)
	}

	err := db.bdb.Update(func(tx *buntdb.Tx) error {
		// to_delete := make([]string, 0, 1000) // TODO: what size should this be?
		to_delete := make(map[string]string, 100) // TODO: what size should this be? Does it matter?

		tx.AscendKeys(KeyEntries, func(key, value string) bool {
			if _, ok := entryIdMapWhitelist[key]; !ok {
				to_delete[key] = value
				// to_delete = append(to_delete, key)
			}
			return true
		})

		for entryId, value := range to_delete {
			name := filepath.Base(getEntryKeyFromId(entryId))
			if newEntryId, ok := namesToEntryIDs[name]; ok {
				fmt.Printf("# Renaming %s to %s\n", color.HiGreenString("%s", getEntryKeyFromId(entryId)), color.HiGreenString("%s", getEntryKeyFromId(newEntryId)))

				entry, err := deserializeEntry(value)
				if err != nil {
					return err
				}

				entry.Location = getEntryKeyFromId(newEntryId)

				updatedValue, err := serializeEntry(entry)
				if err != nil {
					return err
				}

				if _, _, err := tx.Set(newEntryId, string(updatedValue), nil); err != nil {
					return err
				}
				merged++
			} else {
				if removed < 10 {
					fmt.Printf("# Removing %s...\n", color.HiRedString("%s", getEntryKeyFromId(entryId)))
				} else if removed == 10 {
					fmt.Printf("# (removing a lot, not listing them here)...\n")
				}

				removed++
			}

			_, err := tx.Delete(entryId)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return removed, merged, err
}

func (db *Database) Remove(key string) error {
	err := db.bdb.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(key)
		return err
	})
	return err
}

func (db *Database) RemoveArr(keys []string) error {
	err := db.bdb.Update(func(tx *buntdb.Tx) error {
		for _, key := range keys {
			_, err := tx.Delete(key)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
