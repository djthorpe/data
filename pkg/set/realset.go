package set

import (
	"fmt"
	"math"
	"strings"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/f32"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type RealSet struct {
	name     string
	v        []float64
	min, max float64
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

// NewRealSet returns an empty array of float64 values
func NewRealSet(name string) data.RealSet {
	this := new(RealSet)
	this.name = name
	this.min, this.max = math.NaN(), math.NaN()
	return this
}

/////////////////////////////////////////////////////////////////////
// METHODS

func (set *RealSet) Name() string {
	return set.name
}

func (set *RealSet) SetName(name string) {
	set.name = name
}

func (set *RealSet) Len() int {
	return len(set.v)
}

func (set *RealSet) Min() float64 {
	return set.min
}

func (set *RealSet) Max() float64 {
	return set.max
}

func (set *RealSet) Append(values ...float64) {
	for _, v := range values {
		set.min = f64Min(set.min, v)
		set.max = f64Max(set.max, v)
		set.v = append(set.v, v)
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (set *RealSet) String() string {
	str := "<realset"
	str += fmt.Sprintf(" name=%q", set.Name())
	if min := set.Min(); math.IsNaN(min) == false {
		str += fmt.Sprintf(" min=%v", min)
	}
	if max := set.Max(); math.IsNaN(max) == false {
		str += fmt.Sprintf(" max=%v", max)
	}
	if len(set.v) > 0 {
		str += " <"
		for _, v := range set.v {
			str += fmt.Sprintf("%s,", f32.String(float32(v)))
		}
		str = strings.TrimSuffix(str, ",") + ">"
	}
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func f64Min(a, b float64) float64 {
	if math.IsNaN(a) {
		return b
	} else if math.IsNaN(a) {
		return a
	} else {
		return math.Min(a, b)
	}
}

func f64Max(a, b float64) float64 {
	if math.IsNaN(a) {
		return b
	} else if math.IsNaN(a) {
		return a
	} else {
		return math.Max(a, b)
	}
}
