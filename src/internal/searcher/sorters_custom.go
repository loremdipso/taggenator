package searcher

import (
	. "internal/data"
	"math/rand"
	"time"
)

func SortRandom(entries Entries) Entries {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(entries), func(i, j int) { entries[i], entries[j] = entries[j], entries[i] })
	return entries
}

func SortReverse(entries Entries) Entries {
	reverse(entries)
	return entries
}

func reverse(numbers Entries) Entries {
	for i := 0; i < len(numbers)/2; i++ {
		j := len(numbers) - i - 1
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	return numbers
}
