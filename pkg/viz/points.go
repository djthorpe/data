package viz

import (
	"fmt"
	"strings"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type points struct {
	name     string
	v        []data.Point
	min, max data.Point
}

/////////////////////////////////////////////////////////////////////
// CONSTANTS

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

// NewPoints returns an empty array of points
func NewPoints(name string) data.Points {
	this := new(points)
	this.name = name
	this.min, this.max = data.NilPoint, data.NilPoint
	return this
}

/////////////////////////////////////////////////////////////////////
// METHODS

func (p *points) Name() string {
	return p.name
}

func (p *points) SetName(name string) {
	p.name = name
}

func (p *points) Len() int {
	return len(p.v)
}

func (p *points) Min() data.Point {
	return p.min
}

func (p *points) Max() data.Point {
	return p.max
}

func (p *points) Append(point data.Point) {
	// Check min and max
	if isNaN(point.X) == false {
		if isNaN(p.min.X) || p.min.X > point.X {
			p.min.X = point.X
		}
		if isNaN(p.max.X) || p.max.X < point.X {
			p.max.X = point.X
		}
	}
	if isNaN(point.Y) == false {
		if isNaN(p.min.Y) || p.min.Y > point.Y {
			p.min.Y = point.Y
		}
		if isNaN(p.max.Y) || p.max.Y < point.Y {
			p.max.Y = point.Y
		}
	}

	// Append value
	p.v = append(p.v, point)
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (p *points) String() string {
	str := "<points"
	str += fmt.Sprintf(" name=%q", p.Name())
	if min := p.Min(); min.IsNil() == false {
		str += fmt.Sprintf(" min=%v", min)
	}
	if max := p.Max(); max.IsNil() == false {
		str += fmt.Sprintf(" max=%v", max)
	}
	if len(p.v) > 0 {
		str += " <"
		for _, v := range p.v {
			str += fmt.Sprintf("%v,", v)
		}
		str = strings.TrimSuffix(str, ",") + ">"
	}
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func isNaN(f float32) bool {
	return f != f
}

/*

func (p *points) WritePath(c data.Canvas) data.CanvasGroup {

}

*/
