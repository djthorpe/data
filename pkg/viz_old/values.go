package viz

import (
	"fmt"
	"math"
	"strings"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type values struct {
	name     string
	v        []float32
	min, max float32
}

/////////////////////////////////////////////////////////////////////
// CONSTANTS

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

// NewValues returns an empty array of values
func NewValues(name string) data.Values {
	this := new(values)
	this.name = name
	this.min, this.max = float32(math.NaN()), float32(math.NaN())
	return this
}

/////////////////////////////////////////////////////////////////////
// METHODS

func (p *values) Name() string {
	return p.name
}

func (p *values) SetName(name string) {
	p.name = name
}

func (p *values) Len() int {
	return len(p.v)
}

func (p *values) Min() float32 {
	return p.min
}

func (p *values) Max() float32 {
	return p.max
}

func (p *values) Append(value float32) {
	// Check min and max
	if isNaN(value) == false {
		if isNaN(p.min) || p.min > value {
			p.min = value
		}
		if isNaN(p.max) || p.max < value {
			p.max = value
		}
	}

	// Append value
	p.v = append(p.v, value)
}

func (p *values) Scale() data.Scale {
	if isNaN(p.min) || isNaN(p.max) {
		return nil
	}
	return NewScale(p.name, p.min, p.max)
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (p *values) String() string {
	str := "<values"
	str += fmt.Sprintf(" name=%q", p.Name())
	if min := p.Min(); math.IsNaN(float64(min)) == false {
		str += fmt.Sprintf(" min=%f", min)
	}
	if max := p.Max(); math.IsNaN(float64(max)) == false {
		str += fmt.Sprintf(" max=%f", max)
	}
	if len(p.v) > 0 {
		str += " <"
		for _, v := range p.v {
			str += fmt.Sprintf("%f,", v)
		}
		str = strings.TrimSuffix(str, ",") + ">"
	}
	return str + ">"
}
