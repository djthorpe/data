package canvas

import (
	"fmt"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/f32"
)

func (e *Element) MoveTo(pt data.Point) data.CanvasPath {
	e.Attr("d", fmt.Sprint("M ", f32.String(pt.X), " ", f32.String(pt.Y)))
	return e
}

func (e *Element) LineTo(pt data.Point) data.CanvasPath {
	e.Attr("d", fmt.Sprint("L ", f32.String(pt.X), " ", f32.String(pt.Y)))
	return e
}

func (e *Element) QuadraticTo(pt, c data.Point) data.CanvasPath {
	e.Attr("d", fmt.Sprint("Q ",
		f32.String(c.X), " ", f32.String(c.Y), ", ",
		f32.String(pt.X), " ", f32.String(pt.Y),
	))
	return e
}

func (e *Element) CubicTo(pt, c1, c2 data.Point) data.CanvasPath {
	e.Attr("d", fmt.Sprint("C ",
		f32.String(c1.X), " ", f32.String(c1.Y), ", ",
		f32.String(c2.X), " ", f32.String(c2.Y), ", ",
		f32.String(pt.X), " ", f32.String(pt.Y),
	))
	return e
}

func (e *Element) ClosePath() data.CanvasPath {
	e.Attr("d", "Z")
	return e
}
