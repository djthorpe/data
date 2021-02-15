package viz

import (
	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Viz struct {
	c data.Canvas
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewViz(c data.Canvas) data.Viz {
	viz := new(Viz)
	if c == nil {
		return nil
	} else {
		viz.c = c
	}
	return viz
}
