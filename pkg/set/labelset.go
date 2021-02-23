package set

import (
	"fmt"
	"strings"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type LabelSet struct {
	name string
	v    []string
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

// NewLabelSet returns an empty array of labels
func NewLabelSet(name string) data.LabelSet {
	this := new(LabelSet)
	this.name = name
	return this
}

/////////////////////////////////////////////////////////////////////
// METHODS

func (set *LabelSet) Name() string {
	return set.name
}

func (set *LabelSet) SetName(name string) {
	set.name = name
}

func (set *LabelSet) Len() int {
	return len(set.v)
}

func (set *LabelSet) Append(labels ...string) {
	set.v = append(set.v, labels...)
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (set *LabelSet) String() string {
	str := "<labelset"
	str += fmt.Sprintf(" name=%q", set.Name())
	if len(set.v) > 0 {
		str += " <"
		for _, v := range set.v {
			str += fmt.Sprintf("%q,", v)
		}
		str = strings.TrimSuffix(str, ",") + ">"
	}
	return str + ">"
}
