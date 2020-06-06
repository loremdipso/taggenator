package searcher

import (
	"fmt"
	"strings"

	. "internal/data"

	"github.com/loremdipso/go_utils"
)

func FilterNone(*Entry) bool {
	return false
}

func FilterAny(*Entry) bool {
	return true
}

func FilterSeen(entry *Entry) bool {
	fmt.Println(entry.Times_Opened)
	return entry.Times_Opened > 0
}

func FilterUnseen(entry *Entry) bool {
	fmt.Println(entry.Times_Opened)
	return entry.Times_Opened == 0
}

func FilterTouched(entry *Entry) bool {
	return entry.HaveManuallyTouched
}
func FilterUnTouched(entry *Entry) bool {
	return !entry.HaveManuallyTouched
}

func FilterSearchInclusive(tokens []string) SimpleFilterFunc {
	return func(entry *Entry) bool {
		for _, token := range tokens {
			if filterSearch(entry, token) {
				return true
			}
		}
		return false
	}
}

func FilterSearchExclusive(tokens []string) SimpleFilterFunc {
	return func(entry *Entry) bool {
		for _, token := range tokens {
			if !filterSearch(entry, token) {
				return false
			}
		}
		return true
	}
}

func filterSearch(entry *Entry, token string) bool {
	if isNegative(token) {
		return containsTagHelper(entry.Tags, token, false) && containsSubstringHelper(entry.Name, token)
	} else {
		return containsTagHelper(entry.Tags, token, false) || containsSubstringHelper(entry.Name, token)
	}
}

func FilterNameInclusive(tokens []string) SimpleFilterFunc {
	return func(entry *Entry) bool {
		for _, token := range tokens {
			if containsSubstringHelper(entry.Name, token) {
				return true
			}
		}
		return false
	}
}

func FilterNameExclusive(tokens []string) SimpleFilterFunc {
	return func(entry *Entry) bool {
		for _, token := range tokens {
			if !containsSubstringHelper(entry.Name, token) {
				return false
			}
		}
		return true
	}
}

func FilterPathInclusive(tokens []string) SimpleFilterFunc {
	return func(entry *Entry) bool {
		for _, token := range tokens {
			if containsSubstringHelper(entry.Location, token) {
				return true
			}
		}
		return false
	}
}

func FilterPathExclusive(tokens []string) SimpleFilterFunc {
	return func(entry *Entry) bool {
		for _, token := range tokens {
			if !containsSubstringHelper(entry.Location, token) {
				return false
			}
		}
		return true
	}
}

func FilterTagsInclusive(tokens []string) SimpleFilterFunc {
	return func(entry *Entry) bool {
		for _, token := range tokens {
			if containsTagHelper(entry.Tags, token, true) {
				return true
			}
		}
		return false
	}
}

func FilterTagsExclusive(tokens []string) SimpleFilterFunc {
	return func(entry *Entry) bool {
		for _, token := range tokens {
			if !containsTagHelper(entry.Tags, token, true) {
				return false
			}
		}
		return true
	}
}

// TODO: a little slow, but cleaner and less work for now
func containsTagHelper(tags []string, token string, strict bool) bool {
	exclude := false
	if isNegative(token) {
		token = token[1:]
		exclude = true
	}

	doesContain := go_utils.ContainsStringFast(tags, token)
	if !doesContain && !strict {
		// Special case: tag wasn't matched exactly and we're not being overly strict,
		// go ahead and check each tag for a substring
		for _, tag := range tags {
			if strings.Contains(tag, token) {
				doesContain = true
				break
			}
		}
	}

	if exclude {
		return !doesContain
	} else {
		return doesContain
	}
}

func containsSubstringHelper(str string, token string) bool {
	exclude := false
	if isNegative(token) {
		token = token[1:]
		exclude = true
	}

	doesContain := strings.Contains(str, token)
	if exclude {
		return !doesContain
	} else {
		return doesContain
	}
}

func isNegative(token string) bool {
	return token[0] == '-'
}

func FilterA(entry *Entry) bool {
	return strings.Index(entry.Name, "A") > -1
}

func FilterB(entry *Entry) bool {
	return strings.Index(entry.Name, "B") > -1
}
