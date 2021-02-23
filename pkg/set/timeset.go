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
	precision TimePrecision
	tz        *time.Location
}

type TimePrecision time.Duration

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	Day         = TimePrecision(24 * time.Hour)
	Hour        = TimePrecision(time.Hour)
	Minute      = TimePrecision(time.Minute)
	Second      = TimePrecision(time.Second)
	Millisecond = TimePrecision(time.Millisecond)
	Microsecond = TimePrecision(time.Microsecond)
	Nanosecond  = TimePrecision(time.Nanosecond)
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

// NewTimeSet returns an empty array of time.Time values
func NewTimeSet(name string) data.TimeSet {
	this := new(TimeSet)
	this.name = name
	this.min, this.max = time.Time{}, time.Time{}
	this.precision = TimePrecision(0)
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
	return time.Duration(set.precision)
}

func (set *TimeSet) Duration() time.Duration {
	if set.min.IsZero() || set.max.IsZero() {
		return 0
	} else {
		return set.max.Sub(set.min)
	}
}

func (set *TimeSet) Append(values ...time.Time) {
	for _, v := range values {
		set.min = timeMin(v, set.min)
		set.max = timeMax(v, set.max)
		set.precision = set.precision.Min(toPrecision(v))
		set.v = append(set.v, v)
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (set *TimeSet) String() string {
	str := "<timeset"
	str += fmt.Sprintf(" name=%q", set.Name())
	if min := set.Min(); min.IsZero() == false {
		str += fmt.Sprintf(" min=%q", toTimeString(set.precision, min))
	}
	if max := set.Max(); max.IsZero() == false {
		str += fmt.Sprintf(" max=%q", toTimeString(set.precision, max))
	}
	if precision := set.Precision(); precision != 0 {
		str += fmt.Sprintf(" precision=%q", TimePrecision(precision))
	}
	if len(set.v) > 0 {
		str += " <"
		for _, v := range set.v {
			str += fmt.Sprintf("%q,", toTimeString(set.precision, v))
		}
		str = strings.TrimSuffix(str, ",") + ">"
	}
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func timeMin(a, b time.Time) time.Time {
	if a.IsZero() {
		return b
	} else if b.IsZero() {
		return a
	} else if a.After(b) {
		return b
	} else {
		return a
	}
}

func timeMax(a, b time.Time) time.Time {
	if a.IsZero() {
		return b
	} else if b.IsZero() {
		return a
	} else if b.After(a) {
		return b
	} else {
		return a
	}
}

func (p TimePrecision) Min(q TimePrecision) TimePrecision {
	if p == 0 {
		return q
	} else if q == 0 {
		return p
	} else if p < q {
		return p
	} else {
		return q
	}
}

func toPrecision(t time.Time) TimePrecision {
	if t.Truncate(time.Duration(Day)) == t {
		return Day
	}
	if t.Truncate(time.Hour) == t {
		return Hour
	}
	if t.Truncate(time.Minute) == t {
		return Minute
	}
	if t.Truncate(time.Second) == t {
		return Second
	}
	if t.Truncate(time.Millisecond) == t {
		return Millisecond
	}
	if t.Truncate(time.Microsecond) == t {
		return Microsecond
	}
	// By default, return a minimum precision of nanosecond
	return Nanosecond
}

func (p TimePrecision) String() string {
	switch {
	case p == Nanosecond:
		return "ns"
	case p == Microsecond:
		return "us"
	case p == Millisecond:
		return "ms"
	case p == Second:
		return "second"
	case p == Minute:
		return "minute"
	case p == Hour:
		return "hour"
	case p == Day:
		return "day"
	default:
		return ""
	}
}

func toTimeString(p TimePrecision, t time.Time) string {
	if t.IsZero() {
		return ""
	}
	switch p {
	case Day:
		return t.Format("2006-01-02")
	case Hour, Minute, Second:
		return t.Format(time.RFC3339)
	case Millisecond, Microsecond, Nanosecond:
		return t.Format(time.RFC3339Nano)
	default:
		return ""
	}
}
