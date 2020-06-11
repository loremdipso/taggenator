package processor

import (
	"strings"

	"github.com/loremdipso/go_utils"

	"github.com/fatih/color"
)

// TODO: make separate package?
type Completer struct {
	sorted_tags []string
}

// NOTE: will modify sorted_tags array
func NewCompleter(sorted_tags []string) *Completer {
	return &Completer{sorted_tags: sorted_tags}
}

func (comp *Completer) Complete(line string) (c []string) {
	line = strings.ToLower(line)
	for _, tag := range comp.sorted_tags {
		if strings.Contains(tag, line) {
			c = append(c, tag)
		}
	}
	return
}

//returns if thing was appended
func (comp *Completer) Append(new_tag string) (success bool) {
	success, comp.sorted_tags = go_utils.InsertIntoSortedListIfNotThereAlready(comp.sorted_tags, new_tag)
	if success {
		color.HiYellow("Brand new: \"%s\"", color.HiBlueString(new_tag))
	}
	return success
}

func (comp *Completer) GetSortedTags() []string {
	return comp.sorted_tags
}
