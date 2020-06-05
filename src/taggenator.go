package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	. "internal/data"

	"internal/database"
	"internal/processor"

	"github.com/fatih/color"
)

const PROJECT_NAME = "taggenator"

//var db *database.Database = nil

func main() {
	log.SetOutput(ioutil.Discard)

	db, err := database.New()
	if err != nil {
		log.Printf("%v", err)
		return
	}

	canExit := make(chan bool, 10)
	canExit <- true

	// NTOE: close handler will also close the database
	queryProcessor := processor.New()
	setupCloseHandler(db, canExit, queryProcessor)
	defer db.Close()

	go tagRemoverTimeout(db, canExit)
	go shrinkerTimeout(db, canExit)

	// createFakeEntries(db, 100)
	// its := 1
	// err = db.TraverseTreeFast()
	// db.Clear()
	// db.Shrink()
	// clearAndShrink(db)

	//db.Shrink()
	// go_utils.Timer(func() {
	// err = addNewFiles(db)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// err = removeMissingFiles(db)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	err = resolveMovedFiles(db)
	if err != nil {
		fmt.Println(err)
		return
	}

	// }, "both :)", its)

	err = queryProcessor.ProcessQuery(os.Args[1:], db, PROJECT_NAME)
	if err != nil {
		log.Println(err)
		return
	}

	//db.Print()
	cleanUp(canExit)
}

func setupCloseHandler(db *database.Database, canExit chan bool, queryProcessor *processor.QueryProcessor) {
	c := make(chan os.Signal)
	signal.Ignore(os.Interrupt, syscall.SIGTERM)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c

		// TODO: how to detect if DB is being written to, currently?
		cleanUp(canExit)
		queryProcessor.Close()
		db.Close()
		os.Exit(0)
	}()
}

func cleanUp(canExit chan bool) {
	log.Println("Cleaning up...")
	for <-canExit == false {
	}
}

// Wait a second, see if the program's still running. If so, then we can shrink
func shrinkerTimeout(db *database.Database, canExit chan bool) {
	time.Sleep(time.Second * 3)
	if <-canExit {
		canExit <- false
		log.Println("Starting shrinking...")
		db.Shrink()
		log.Println("Ending shrinking...")
		canExit <- true
	}
}

func tagRemoverTimeout(db *database.Database, canExit chan bool) {
	time.Sleep(time.Second * 2)
	_, err := db.RemoveMissingTags()
	if <-canExit {
		canExit <- false
		if err != nil {
			// TODO: log this?
			log.Println(err)
			// return len(addedArr), merged, removed, err
		}
		canExit <- true
	}
}

//#region helpers
func clearAndShrink(db *database.Database) {
	db.Clear()
	db.Shrink()
}

func resolveMovedFiles(db *database.Database) error {
	added, modified, removed, err := db.UpdateFiles()
	if err != nil {
		return err
	}

	fmt.Printf("# Added %s entries\n", color.HiGreenString("%d", added))
	fmt.Printf("# Removed %s entries\n", color.HiRedString("%d", removed))
	fmt.Printf("# Modified %s entries\n", color.HiYellowString("%d", modified))

	return nil
}

func createFakeEntries(db *database.Database, numFake int) {
	for i := 0; i < numFake; i++ {
		db.UpdateEntry(&Entry{Name: fmt.Sprintf("Name B_%d", i), Tags: []string{"shouldmatch", "tag 2"}})
	}
}

func searchTest(db *database.Database) {
	// var _ Entry
	// log.Println("searching...")
	// temp, err := searcher.TestSearch(db)
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println("matches...", len(temp))

	// for _, value := range temp {
	// 	fmt.Println(value)
	// }

	// tmp := db.SimpleSearch(func(entry *database.Entry) bool {
	// 	return true
	// 	if entry.Name == "Name A" {
	// 		return true
	// 	}
	// 	return false
	// })
}

//#endregion helpers
