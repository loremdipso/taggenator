package searcher

import (
	. "internal/data"
)

func SortNewestAdded(entries Entries, i int, j int) bool {
	// if i - j < 0, then i < j
	return entries[i].Date_Added.Sub(entries[j].Date_Added).Nanoseconds() > 0
}

func SortOldestAdded(entries Entries, i int, j int) bool {
	return !SortNewestAdded(entries, i, j)
}

func SortMostRecentlyOpened(entries Entries, i int, j int) bool {
	return entries[i].Date_Last_Accessed.Sub(entries[j].Date_Last_Accessed).Nanoseconds() > 0
}

func SortLeastRecentlyOpened(entries Entries, i int, j int) bool {
	return !SortMostRecentlyOpened(entries, i, j)
}

func SortMostFrequentlyOpened(entries Entries, i int, j int) bool {
	return entries[i].Times_Opened > entries[j].Times_Opened
}

func SortLeastFrequentlyOpened(entries Entries, i int, j int) bool {
	return entries[i].Times_Opened > entries[j].Times_Opened
}

func SortLargest(entries Entries, i int, j int) bool {
	return entries[i].Size > entries[j].Size
}
func SortSmallest(entries Entries, i int, j int) bool {
	return !SortLargest(entries, i, j)
}

func SortLocation(entries Entries, i int, j int) bool {
	return entries[i].Location < entries[j].Location
}

func SortName(entries Entries, i int, j int) bool {
	return entries[i].Name < entries[j].Name
}
