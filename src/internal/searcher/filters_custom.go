package searcher

import (
	. "internal/data"

	"github.com/loremdipso/go_utils"
)

func FilterMaxCapacity(maxCapacity int) CustomFilterFunc {
	return func(entries Entries) Entries {
		return entries[:go_utils.Min(maxCapacity, len(entries))]
	}
}
