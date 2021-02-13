package canvas

import (
	"fmt"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// METHODS

func (e *Element) Circle(centre data.Point, radius float32) data.CanvasElement {
	c := NewElement("circle", "", e.root)
	c.Attr("cx", centre.X)
	c.Attr("cy", centre.Y)
	c.Attr("r", radius)
	e.addChild(c)
	return c
}

func (e *Element) Ellipse(centre data.Point, radius data.Size) data.CanvasElement {
	c := NewElement("ellipse", "", e.root)
	c.Attr("cx", centre.X)
	c.Attr("cy", centre.Y)
	c.Attr("rx", radius.W)
	c.Attr("ry", radius.H)
	e.addChild(c)
	return c
}

func (e *Element) Line(p1, p2 data.Point) data.CanvasElement {
	c := NewElement("line", "", e.root)
	c.Attr("x1", p1.X)
	c.Attr("y1", p1.Y)
	c.Attr("x2", p2.X)
	c.Attr("y2", p2.Y)
	e.addChild(c)
	return c
}

func (e *Element) Rect(pt data.Point, sz data.Size) data.CanvasElement {
	c := NewElement("rect", "", e.root)
	c.Attr("x", pt.X)
	c.Attr("y", pt.Y)
	c.Attr("width", sz.W)
	c.Attr("height", sz.H)
	e.addChild(c)
	return c
}

func (e *Element) Text(pt data.Point, children ...data.CanvasText) data.CanvasElement {
	t := NewElement("text", "", e.root)
	t.Attr("x", pt.X)
	t.Attr("y", pt.Y)
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

func (e *Element) Path(pts []data.Point) data.CanvasElement {
	// We require at least two elements for a path
	if len(pts) < 2 {
		return nil
	}
	c := NewElement("path", "", e.root)
	attr := ""
	for i, pt := range pts {
		if i == 0 {
			attr += fmt.Sprintf("M %f %f", pt.X, pt.Y)
		} else {
			attr += fmt.Sprintf(" L %f %f", pt.X, pt.Y)
		}
	}
	c.Attr("d", attr)
	e.addChild(c)
	return c
}
