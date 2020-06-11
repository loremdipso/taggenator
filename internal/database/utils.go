package database

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	. "internal/data"

	"github.com/tidwall/buntdb"
)

func (db *Database) Print() {
	db.bdb.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend("", func(key, value string) bool {
			fmt.Printf("%s: %s\n", key, value)
			return true
		})
		return err
	})
}

// func (db *Database) Print() {
// 	db.db.View(func(tx *buntdb.Tx) error {
// 		//err := tx.Ascend("", func(key, value string) bool {
// 		// err := tx.AscendKeys()
// 		err := tx.Ascend(index_name, func(key, value string) bool {
// 			fmt.Printf("%s: %s\n", key, value)
// 			return true
// 		})
// 		return err
// 	})
// }

func (db *Database) GetKeys() []string {
	var results []string
	db.bdb.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend("", func(key, _ string) bool {
			results = append(results, key)
			return true
		})
		return err
	})
	return results
}

func (db *Database) GetValues() []string {
	var results []string
	db.bdb.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend("", func(_, value string) bool {
			results = append(results, value)
			return true
		})
		return err
	})
	return results
}

func (db *Database) Clear() error {
	err := db.bdb.Update(func(tx *buntdb.Tx) error {
		err := tx.DeleteAll()
		return err
	})
	return err
}

func serializeEntry(entry *Entry) ([]byte, error) {
	return json.Marshal(entry)
}

func deserializeEntry(value string) (*Entry, error) {
	var entry Entry
	err := json.Unmarshal([]byte(value), &entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func PathToName(pathname string) string {
	return pathToName(pathname)
}

func pathToName(pathname string) string {
	basename := filepath.Base(pathname)
	without_extension := strings.TrimSuffix(basename, filepath.Ext(basename))
	return strings.ToLower(without_extension)
}

func getEntryId(s_entryKey string) string {
	return id_entry_prefix + s_entryKey
}

func getEntryKeyFromId(s_entryId string) string {
	return s_entryId[len(id_entry_prefix):]
}

func getTagId(s_tag string) string {
	return id_tag_prefix + s_tag
}
