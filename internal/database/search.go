package database

import (
	"errors"

	"internal/data"
	. "internal/data"

	"github.com/tidwall/buntdb"
)

// TODO: think hard about this
// how to compose efficiently?
func (db *Database) SimpleDatabaseSearch(index string, filterFunc func(*Entry) bool) (Entries, error) {
	if filterFunc == nil {
		return nil, errors.New("Missing function")
	}

	var results Entries
	db.bdb.View(func(tx *buntdb.Tx) error {
		//err := tx.Ascend(index_name, func(key, value string) bool {
		//err := tx.Ascend("", func(key, value string) bool {
		//err := tx.Ascend(id_entry_prefix+"*", func(key, value string) bool {
		err := tx.AscendKeys(id_entry_prefix+"*", func(key, value string) bool {
			// fmt.Println(key, value)
			entry, err := deserializeEntry(value)
			if err != nil {
				// TODO: how do I handle this?
				// panic(err)
			} else {
				if filterFunc(entry) {
					results = append(results, entry)
				}
			}
			return true
		})
		return err
	})

	return results, nil
}

func (db *Database) AsyncDatabaseSearch(index string, filterFunc data.SimpleFilterFunc, entryChannel chan *Entry, getIsChannelClosed func() bool) int {
	total := 0
	db.bdb.View(func(tx *buntdb.Tx) error {
		err := tx.AscendKeys(index, func(key, value string) bool {
			if getIsChannelClosed() {
				return false
			}

			entry, err := deserializeEntry(value)
			if err != nil {
				// TODO: how do I handle this?
				// panic(err)
			} else {
				if filterFunc(entry) {
					entryChannel <- entry
					total++
				}
			}
			return true
		})
		return err
	})

	return total
}
