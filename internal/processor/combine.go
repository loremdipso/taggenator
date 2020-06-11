package processor

import (
	"encoding/json"
	"fmt"
	"internal/data"
	"internal/database"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func absorb_old_database(self *QueryProcessor, args []string, db *database.Database) error {
	// TODO: remove or clean. This is just dirty
	dbname := args[0]

	// Open our jsonFile
	jsonFile, err := os.Open(dbname)
	// if we os.Open returns an error then handle it
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	var historic data.HistoricDB
	err = json.Unmarshal(byteValue, &historic)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(historic)

	// TODO: do we need to worry about casing?
	for _, entry := range historic {
		newEntry := data.Entry{
			Name:         database.PathToName(entry.Location),
			Tags:         entry.Tags,
			Location:     entry.Location,
			Times_Opened: entry.Times_Opened,
			Length:       getLength(entry.Length),

			// Date_Added:    entry.Creation_Date,
			Date_Added:         stringToDate(entry.Add_Date),
			Date_Last_Accessed: intToDate(entry.Access_Date),
			Date_Created:       stringToDate(entry.Creation_Date),
		}

		// fix tags by removing surrounding whitespace
		for i, tag := range newEntry.Tags {
			newEntry.Tags[i] = strings.Trim(tag, "	 \n")
		}

		fmt.Println()
		fmt.Printf("%+v\n", entry)
		fmt.Printf("%+v\n", newEntry)

		db.UpdateEntry(&newEntry)
	}

	return nil
}

func getLength(length interface{}) int64 {
	if rv, ok := length.(int64); ok {
		return rv
	}

	fmt.Println("copping out")
	return -1
}

func stringToDate(val string) time.Time {
	// example val = "2016-09-14 00:43:24 -0500"
	layout := "2006-01-02 15:04:05 -0700"
	dte, err := time.Parse(layout, val)
	if err != nil {
		fmt.Println(err)
	}
	return dte
}

func intToDate(val int64) time.Time {
	return time.Unix(val, 0)
}
