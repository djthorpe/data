package canvas

import (
	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/f32"
)

type PathSegment string

func NewPathSegment(op string, args ...float32) PathSegment {
	if len(args) > 0 {
		return PathSegment(op + " " + f32.Join(args, " "))
	} else {
		return PathSegment(op)
	}
}

func (*Canvas) MoveTo(pt data.Point) data.CanvasPath {
	return NewPathSegment("M", pt.X, pt.Y)
}

func (*Canvas) LineTo(pt data.Point) data.CanvasPath {
	return NewPathSegment("L", pt.X, pt.Y)
}

func (*Canvas) QuadraticTo(pt, c data.Point) data.CanvasPath {
	return NewPathSegment("Q", c.X, c.Y, pt.X, pt.Y)
}

func (*Canvas) CubicTo(pt, c1, c2 data.Point) data.CanvasPath {
	return NewPathSegment("C", c1.X, c1.Y, c2.X, c2.Y, pt.X, pt.Y)
}

func (*Canvas) ClosePath() data.CanvasPath {
	return NewPathSegment("Z")
}
