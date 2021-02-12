package viz

import (
	"fmt"
	"math"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type scale struct {
	name string
}

/////////////////////////////////////////////////////////////////////
// CONSTANTS

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

// NewScale returns an X or Y axis scale given minimum
// and maximum values that need to be included on the scale
func NewScale(name string, min, max float32) data.Scale {
	this := new(scale)
	this.name = name

	// Calculate scale min and max
	minl := int(math.Floor(math.Log10(float64(min))))
	maxl := int(math.Ceil(math.Log10(float64(max))))
	fmt.Printf("min=%f log=%v pow=%f\n", min, minl, math.Pow10(minl-1))
	fmt.Printf("max=%f log=%v pow=%f\n", max, maxl, math.Pow10(maxl-1))

	return this
}

/////////////////////////////////////////////////////////////////////
// METHODS

func (p *scale) Name() string {
	return p.name
}
