package processor

import (
	"errors"
	"fmt"
	"internal/database"
	"internal/searcher"
	"strconv"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/loremdipso/go_utils"
	"github.com/robpike/filter"
)

func apply_tags(self *QueryProcessor, args []string, db *database.Database) error {
	// TODO: speed up. This takes ~5ms right now for just 1000K entries
	// fmt.Printf("DURATION: %v\n", time.Since(start))
	numThreads := 1

	tags := make([]string, 0)
	for i, arg := range args {
		switch arg {
		case "-tag", "--tag":
			tags = append(tags, args[i+1])

			// TODO: a little messy, but we can handle it
			args[i] = ""
			args[i+1] = ""
			break
		case "-threads", "--numthreads":
			var err error
			numThreads, err = strconv.Atoi(args[i+1])
			if err != nil {
				return err
			}

			// TODO: a little messy, but we can handle it
			args[i] = ""
			args[i+1] = ""
			break
		}
	}

	if len(tags) == 0 {
		return errors.New("no tags to add")
	} else {
		tags = filter.Choose(tags, func(el string) bool {
			return !strings.HasPrefix(el, "#")
		}).([]string)
		fmt.Printf("Tags to add: %s\n", go_utils.StringArrayToString(tags))
	}

	args = go_utils.RemoveEmpty(args)

	search := searcher.New(db)
	err := search.Parse(args)
	if err != nil {
		return err
	}

	entries, err := search.Execute()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	jobs := make(chan int, 1) // make 10 workers for this
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func(workerIndex int) {
			defer wg.Done()

			for {
				entryIndex, ok := <-jobs
				if !ok || entryIndex == -1 {
					return
				}

				entry := entries[entryIndex]
				fmt.Printf("\n\nWorker %d: %d / %d\n", workerIndex+1, entryIndex+1, len(entries))
				fmt.Printf("Auto-adding tags for: %s\n", color.HiGreenString(entry.Name))

				for _, tag := range tags {
					addTagStringToEntry(db, entry, tag, nil, nil)
				}
				db.UpdateEntry(entry)
			}
		}(i)
	}

	for i, _ := range entries {
		jobs <- i
	}

	jobs <- -1 // signal to close. Used so that we don't lose a final job when closing that channel
	close(jobs)
	wg.Wait()
	return nil
}
