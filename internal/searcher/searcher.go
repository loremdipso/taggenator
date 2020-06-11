package searcher

import (
	"errors"
	"sort"

	. "internal/data"

	"internal/database"
)

func TestSearch(db *database.Database) (Entries, error) {
	searchAndFilters := []*SearchAndFilter{
		// {SimpleFilter: FilterAny},
		{SimpleFilter: FilterNone},
		// {SimpleFilter: InclusiveSearch([]string{"shouldmatch"})},
		//{SimpleFilter: FilterB},
	}

	return BatchSearchAndFilter(db, database.KeyEntries, searchAndFilters)
}

type Searcher struct {
	searchAndFilters []*SearchAndFilter
	db               *database.Database
	initialIndex     string
}

func New(db *database.Database, sizes ...int) *Searcher {
	var initialCapacity int = 5
	if len(sizes) > 0 {
		initialCapacity = sizes[0]
	}
	return &Searcher{make([]*SearchAndFilter, 0, initialCapacity), db, database.KeyEntries}
}

func (s *Searcher) Execute() (Entries, error) {

	return BatchSearchAndFilter(s.db, s.initialIndex, s.searchAndFilters)
}

func (s *Searcher) SetInitialIndex(initialIndex string) {
	s.initialIndex = initialIndex
}

func (s *Searcher) AppendSimpleFilter(fn SimpleFilterFunc) {
	s.searchAndFilters = append(s.searchAndFilters, &SearchAndFilter{SimpleFilter: fn})
}

func (s *Searcher) AppendCustomFilter(fn CustomFilterFunc) {
	s.searchAndFilters = append(s.searchAndFilters, &SearchAndFilter{CustomFilter: fn})
}

func (s *Searcher) AppendSimpleSorter(fn SimpleSortFunc) {
	s.searchAndFilters = append(s.searchAndFilters, &SearchAndFilter{SimpleSort: fn})
}

func (s *Searcher) AppendCustomSorter(fn CustomSortFunc) {
	s.searchAndFilters = append(s.searchAndFilters, &SearchAndFilter{CustomSort: fn})
}

func BatchSearchAndFilter(db *database.Database, initialIndex string, searchAndFilters []*SearchAndFilter) (Entries, error) {
	var entries Entries
	var err error
	if len(searchAndFilters) == 0 {
		return nil, errors.New("No searches/filters, hauss")
	}

	for i, searchAndFilter := range searchAndFilters {
		if i == 0 {
			entries, err = db.SimpleDatabaseSearch(initialIndex, searchAndFilter.SimpleFilter)
			if err != nil {
				return nil, err
			}
		} else if searchAndFilter.SimpleFilter != nil {
			entries = filterEntries(entries, searchAndFilter.SimpleFilter)
		}

		if searchAndFilter.CustomFilter != nil {
			entries = searchAndFilter.CustomFilter(entries)
		}

		if searchAndFilter.SimpleSort != nil {
			entries = sortEntries(entries, searchAndFilter.SimpleSort)
		}
		if searchAndFilter.CustomSort != nil {
			entries = searchAndFilter.CustomSort(entries)
		}

	}
	return entries, nil
}

// func (s *Searcher) ExecuteAsync(entryChannel chan *Entry, finishedCallback func(int), getIsChannelClosed func() bool) {
// 	finishedCallback(BatchSearchAndFilterAsync(s.db, s.initialIndex, s.CombinedFilter, entryChannel, getIsChannelClosed))
// }

// func (s *Searcher) CombinedFilter(entry *Entry) bool {
// 	// fmt.Println("ugg")
// 	// time.Sleep(time.Second)
// 	// time.Sleep(time.Millisecond)
// 	for _, searchAndFilter := range s.searchAndFilters {
// 		if searchAndFilter.SimpleFilter != nil {
// 			if !searchAndFilter.SimpleFilter(entry) {
// 				return false
// 			}
// 			// TODO: if we go back to this, implement the other kinds of sorts/filters
// 		}

// 		// TODO: does this still make sense? No, right?
// 		// if searchAndFilter.CustomFilter != nil {
// 		// 	entries = searchAndFilter.CustomFilter(entries)
// 		// }
// 		// if searchAndFilter.CustomSort != nil {
// 		// 	entries = searchAndFilter.CustomSort(entries)
// 		// }
// 	}

// 	return true
// }

// func BatchSearchAndFilterAsync(db *database.Database, initialIndex string, combinedFilter data.SimpleFilterFunc, entryChannel chan *Entry, getIsChannelClosed func() bool) int {
// 	return db.AsyncDatabaseSearch(initialIndex, combinedFilter, entryChannel, getIsChannelClosed)
// }

func filterEntries(values Entries, f SimpleFilterFunc) Entries {
	matches := make(Entries, 0, len(values))
	for _, v := range values {
		if f(v) {
			matches = append(matches, v)
		}
	}
	return matches
}

func sortEntries(entries Entries, fn SimpleSortFunc) Entries {
	sort.SliceStable(entries, func(i, j int) bool {
		return fn(entries, i, j)
	})
	return entries
}

func SimpleMemorySearch(entries Entries, filterFunc func(*Entry) bool) Entries {
	var results []*Entry
	for _, entry := range entries {
		if filterFunc(entry) {
			results = append(results, entry)
		}
	}
	return results
}
