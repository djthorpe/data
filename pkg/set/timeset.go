package set

import (
	"fmt"
	"strings"
	"time"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type TimeSet struct {
	name      string
	v         []time.Time
	min, max  time.Time
	precision time.Duration
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

// NewTimeSet returns an empty array of time.Time values
func NewTimeSet(name string) data.TimeSet {
	this := new(TimeSet)
	this.name = name
	this.min, this.max = time.Time{}, time.Time{}
	return this
}

/////////////////////////////////////////////////////////////////////
// METHODS

func (set *TimeSet) Name() string {
	return set.name
}

func (set *TimeSet) SetName(name string) {
	set.name = name
}

func (set *TimeSet) Len() int {
	return len(set.v)
}

func (set *TimeSet) Min() time.Time {
	return set.min
}

func (set *TimeSet) Max() time.Time {
	return set.max
}

func (set *TimeSet) Precision() time.Duration {
	return set.precision
}

func (set *TimeSet) Append(values ...time.Time) {
	for _, v := range values {
		// TODO
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (set *TimeSet) String() string {
	str := "<timeset"
	str += fmt.Sprintf(" name=%q", set.Name())
	if min := set.Min(); min.IsZero() == false {
		str += fmt.Sprintf(" min=%v", min.Format(time.RFC3339))
	}
	if max := set.Max(); max.IsZero() == false {
		str += fmt.Sprintf(" max=%v", max.Format(time.RFC3339))
	}
	if precision := set.Precision(); precision != 0 {
		str += fmt.Sprintf(" precision=%v", this.toStringPrecision(precision))
	}
	if len(set.v) > 0 {
		str += " <"
		for _, v := range set.v {
			str += fmt.Sprintf("%s,", this.toStringTime(v))
		}
		str = strings.TrimSuffix(str, ",") + ">"
	}
	return str + ">"
}
