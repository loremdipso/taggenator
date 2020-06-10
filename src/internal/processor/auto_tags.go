package processor

import (
	"fmt"
	"internal/data"
	"internal/database"
	"path/filepath"
	"strings"

	"github.com/loremdipso/go_utils"

	"github.com/fatih/color"
	combinations "github.com/mxschmitt/golang-combinations"
)

func addAutoTagsToEntry(db *database.Database, entry *data.Entry, getSortedTags func() []string, onTagAdd func(string) bool) []string {
	// TODO: deal with multi-token subsequences and things like synonyms
	// Also, make sure we don't auto-add unsafe sequences, like when a
	// sequence is a substring of a tag
	sorted_tags := getSortedTags()
	tokens := getTokens(strings.ToLower(entry.Location))
	tags_to_auto_add := make([]string, 0)

	prefixes := db.GetPrefixes()
	tags_to_auto_add = appendHelper(entry, tags_to_auto_add, sorted_tags, tokens, "")
	for _, prefix := range prefixes {
		tags_to_auto_add = appendHelper(entry, tags_to_auto_add, sorted_tags, tokens, prefix)
	}

	if len(tags_to_auto_add) > 0 {
		fmt.Printf("%s %s\n", color.HiYellowString("Auto Adding..."), go_utils.StringArrayToString(tags_to_auto_add))
		addTagsToEntry(db, entry, tags_to_auto_add, onTagAdd)
	}

	return tags_to_auto_add
}

func appendHelper(entry *data.Entry, tags_to_auto_add []string, sorted_tags []string, tokens []string, prefix string) []string {
	for _, token := range tokens {
		if len(prefix) > 0 {
			token = prefix + token
		}

		// make sure we don't already have this one
		if !go_utils.ContainsStringFast(entry.Tags, token) {

			// make sure this tag exists
			if go_utils.ContainsStringFast(sorted_tags, token) {
				tags_to_auto_add = append(tags_to_auto_add, token)
			}
		}
	}

	return tags_to_auto_add
}

func getTokens(path string) []string {
	// TODO: slow and silly. Fix this
	path = strings.TrimSuffix(path, filepath.Ext(path))
	badChars := []string{"_", "-", "/", ".", "\\", "[", "]"}
	for _, str := range badChars {
		path = strings.ReplaceAll(path, str, " ")
	}

	sillyChars := []string{"'", "\""}
	for _, str := range sillyChars {
		path = strings.ReplaceAll(path, str, "")
	}

	// TODO: permutations
	arr := strings.Split(path, " ")
	arr = go_utils.RemoveEmpty(arr)
	return append(arr, append(ngrams(arr, 2), ngrams(arr, 3)...)...)
}

func ngrams(arr []string, n int) []string {
	/* how?
	for each token
	*/
	raw := combinations.Combinations(arr, n)
	result := make([]string, len(raw))
	for i, set := range raw {
		result[i] = strings.Join(set, " ")
	}

	return result
}
