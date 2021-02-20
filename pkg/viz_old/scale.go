package viz

import (
	"fmt"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type scale struct {
	name     string
	min, max float32
}

/////////////////////////////////////////////////////////////////////
// CONSTANTS

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

// NewScale returns an X or Y axis scale given minimum
// and maximum values that need to be included on the scale
// currently a linear scale
func NewScale(name string, min, max float32) data.Scale {
	this := new(scale)
	this.name = name

	// Check parameters
	if isNaN(min) || isNaN(max) || min > max {
		return nil
	} else {
		this.min = min
		this.max = max
	}

	return this
}

/////////////////////////////////////////////////////////////////////
// METHODS

func (p *scale) Name() string {
	return p.name
}

func (p *scale) Min() float32 {
	return p.min
}

func (p *scale) Max() float32 {
	return p.max
}

func (p *scale) WritePath(data.Canvas) data.CanvasGroup {
	return c.Group(
		c.Line(data.Point{p.min, 0}, data.Point{p.max, 0}),
	)
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (p *scale) String() string {
	str := "<scale"
	str += fmt.Sprintf(" min=%f", p.min)
	str += fmt.Sprintf(" max=%f", p.max)
	return str + ">"
}
