package viz

import (
	"fmt"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type points struct {
	series string
	p      []data.Point
}

/////////////////////////////////////////////////////////////////////
// CONSTANTS

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

// NewPoints returns an empty array of points
func NewPoints(series string) data.Points {
	this := new(points)
	this.series = series
	return this
}

/////////////////////////////////////////////////////////////////////
// METHODS

func (p *points) Series() string {
	return p.series
}

func (p *points) Read(t data.Table, fn data.PointIteratorFunc) error {
	if t == nil || fn == nil {
		return data.ErrBadParameter
	}
	for i := 0; i < t.Len(); i++ {
		if pt, err := fn(i, t.Row(i)); err == data.ErrSkipTransform {
			continue
		} else if err != nil {
			return err
		} else {
			p.p = append(p.p, pt)
		}
	}

	// Return success
	return nil
}

func (p *points) WritePath(c data.Canvas) data.CanvasGroup {
	return c.Group(
		c.Path(p.p),
	)
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (p *points) String() string {
	str := "<points"
	str += fmt.Sprintf(" series=%q", p.Series())
	for _, pt := range p.p {
		str += fmt.Sprint(" <", pt.X, ",", pt.Y, ">")
	}
	return str + ">"
}
