package search_iterator

import (
	"errors"
	"fmt"
	"internal/data"
	"internal/searcher"
)

type SearchIterator struct {
	// search       *searcher.Searcher
	entryChannel      chan *data.Entry
	cumulativeEntries data.Entries
	totalSize         int
	isFinished        bool
}

func New(search *searcher.Searcher) *SearchIterator {
	entryChannel := make(chan *data.Entry, 10000)                                // TODO: what size should this be?
	sit := &SearchIterator{entryChannel, make(data.Entries, 0, 1000), -1, false} // TODO: what size should THIS be?
	go search.ExecuteAsync(entryChannel, func(totalSize int) {
		sit.totalSize = totalSize
		sit.isFinished = true
		close(sit.entryChannel)
	}, func() bool {
		return sit.isFinished
	})
	return sit
}

func (sit *SearchIterator) Close() {
	// TODO: end search early if ctrl+c
	if !sit.isFinished {
		sit.isFinished = true
		fmt.Println("ENDING")
	}
	// close(sit.entryChannel)
}

// blocking
// func (sit *SearchIterator) GetTotalMatches() int {
// }

// also blocking
func (sit *SearchIterator) HasAtLeast(i int) bool {
	isFinished, totalSize := sit.TryGetTotalMatches()

	// need to off-by-one since arrays are 0-indexed
	i++
	if isFinished {
		return totalSize >= i
	} else {
		for len(sit.cumulativeEntries) < i && sit.Eat() {
		}
		return len(sit.cumulativeEntries) >= i
	}
}

// non-blocking. Returns: hasFinishedSearching, totalMatches
func (sit *SearchIterator) TryGetTotalMatches() (bool, int) {
	if sit.isFinished {
		return true, sit.totalSize
	}

	return false, -1
}

func (sit *SearchIterator) GetEntry(i int) (*data.Entry, error) {
	if sit.HasAtLeast(i) {
		return sit.cumulativeEntries[i], nil
	} else {
		return nil, errors.New("ERROR: out of range")
	}
}

// returns true if we successfully ate
func (sit *SearchIterator) Eat() bool {
	entry, ok := <-sit.entryChannel
	if !ok {
		return false
	}

	sit.cumulativeEntries = append(sit.cumulativeEntries, entry)
	return true
}
