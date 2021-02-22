package sets

import (
	"fmt"
	"strings"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/f32"
	"github.com/djthorpe/data/pkg/geom"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type PointSet struct {
	name     string
	v        []data.Point
	min, max data.Point
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

// NewPointSet returns an empty array of points
func NewPointSet(name string) data.PointSet {
	this := new(PointSet)
	this.name = name
	this.min, this.max = data.NilPoint, data.NilPoint
	return this
}

/////////////////////////////////////////////////////////////////////
// METHODS

func (p *PointSet) Name() string {
	return p.name
}

func (p *PointSet) SetName(name string) {
	p.name = name
}

func (p *PointSet) Len() int {
	return len(p.v)
}

func (p *PointSet) Min() data.Point {
	return p.min
}

func (p *PointSet) Max() data.Point {
	return p.max
}

func (p *PointSet) Append(pts ...data.Point) {
	for _, pt := range pts {
		p.min.X, p.min.Y = f32.Min(p.min.X, pt.X), f32.Min(p.min.Y, pt.Y)
		p.max.X, p.max.Y = f32.Max(p.min.X, pt.X), f32.Max(p.min.Y, pt.Y)
		p.v = append(p.v, pt)
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (p *PointSet) String() string {
	str := "<pointset"
	str += fmt.Sprintf(" name=%q", p.Name())
	if min := p.Min(); geom.IsNilPoint(min) == false {
		str += fmt.Sprintf(" min=%v", min)
	}
	if max := p.Max(); geom.IsNilPoint(max) == false {
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
