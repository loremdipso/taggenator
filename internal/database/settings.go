package database

import (
	"encoding/json"
	"internal/data"
	"io/ioutil"
	"os"

	"github.com/loremdipso/go_utils"
)

type Settings struct {
	Extensions   []string
	Synonyms     map[string]string
	Prefixes     []string
	Commands     map[string]string
	Tagger       map[string]string
	OpenerConfig data.OpenerConfig
}

func getSettings(filename string) (*Settings, error) {
	if !go_utils.FileExists(filename) {
		return getDefaultSettings(), nil
	}

	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteArr, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var settings Settings
	err = json.Unmarshal(byteArr, &settings)
	if err != nil {
		return nil, err
	}

	return &settings, nil
}

func getDefaultSettings() *Settings {
	return &Settings{Extensions: []string{"*"}, Synonyms: map[string]string{}}
}

func (db *Database) GetOpenerConfig() data.OpenerConfig {
	return db.settings.OpenerConfig
}

func (db *Database) GetPrefixes() []string {
	return db.settings.Prefixes
}
