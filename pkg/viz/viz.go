package viz

import (
	"io"
	"strings"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Viz struct {
	data.Canvas
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewViz(c data.Canvas) data.Viz {
	viz := new(Viz)
	if c == nil {
		return nil
	} else {
		viz.Canvas = c
	}
	return viz
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Viz) String() string {
	b := new(strings.Builder)
	if err := this.Write(0, b); err != nil {
		panic(err)
	} else {
		return b.String()
	}
}

/////////////////////////////////////////////////////////////////////
// WRITE VIZ

func (this *Viz) Write(opts data.Writer, w io.Writer) error {
	return this.Canvas.Write(opts, w)
}
