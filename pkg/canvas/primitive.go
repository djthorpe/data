package canvas

import (
	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/f32"
)

/////////////////////////////////////////////////////////////////////
// METHODS

func (e *Element) Circle(centre data.Point, radius float32) data.CanvasElement {
	c := NewElement("circle", "", e.root)
	c.Attr("cx", f32.String(centre.X))
	c.Attr("cy", f32.String(centre.Y))
	c.Attr("r", f32.String(radius))
	e.addChild(c)
	return c
}

func (e *Element) Ellipse(centre data.Point, radius data.Size) data.CanvasElement {
	c := NewElement("ellipse", "", e.root)
	c.Attr("cx", f32.String(centre.X))
	c.Attr("cy", f32.String(centre.Y))
	c.Attr("rx", f32.String(radius.W))
	c.Attr("ry", f32.String(radius.H))
	e.addChild(c)
	return c
}

func (e *Element) Line(p1, p2 data.Point) data.CanvasElement {
	c := NewElement("line", "", e.root)
	c.Attr("x1", f32.String(p1.X))
	c.Attr("y1", f32.String(p1.Y))
	c.Attr("x2", f32.String(p2.X))
	c.Attr("y2", f32.String(p2.Y))
	e.addChild(c)
	return c
}

func (e *Element) Rect(pt data.Point, sz data.Size) data.CanvasElement {
	c := NewElement("rect", "", e.root)
	c.Attr("x", f32.String(pt.X))
	c.Attr("y", f32.String(pt.Y))
	c.Attr("width", f32.String(sz.W))
	c.Attr("height", f32.String(sz.H))
	e.addChild(c)
	return c
}

func (e *Element) Text(pt data.Point, children ...data.CanvasText) data.CanvasElement {
	t := NewElement("text", "", e.root)
	t.Attr("x", f32.String(pt.X))
	t.Attr("y", f32.String(pt.Y))
	e.addChild(t)
	for _, node := range children {
		e.removeChild(node.(*Element))
		t.addChild(node.(*Element))
	}
	return t
}

func (e *Element) Span(cdata string) data.CanvasText {
	c := NewElement("tspan", cdata, e.root)
	e.addChild(c)
	return c
}

func (e *Element) Path(pts []data.Point) data.CanvasPath {
	c := NewElement("path", "", e.root)
	for i, pt := range pts {
		if i == 0 {
			c.MoveTo(pt)
		} else {
			c.LineTo(pt)
		}
	}
	e.addChild(c)
	return c
}
