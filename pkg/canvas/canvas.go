package canvas

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	xmlNs = "http://www.w3.org/2000/svg"
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewCanvas(size data.Size, units data.Unit) data.Canvas {
	c := NewElement("svg", "", nil)
	c.XMLName.Space = xmlNs
	w, h := size.UnitString(units)
	if w != "" && w != "0" {
		c.Attr("width", w)
	}
	if h != "" && h != "0" {
		c.Attr("height", h)
	}

	// Set origin and size
	c.origin = data.Point{0, 0}
	c.size = size

	// Set the root element
	c.root = c

	// Return the root element
	return c
}

/////////////////////////////////////////////////////////////////////
// METHODS

func (e *Element) Title(value string) data.Canvas {
	e.addChild(NewElement("title", value, e.root))
	return e
}

func (e *Element) Version(value string) data.Canvas {
	e.Attr("version", value)
	return e
}

func (e *Element) Origin() data.Point {
	return e.root.origin
}

func (e *Element) Size() data.Size {
	return e.root.size
}

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

func (e *Element) Write(w io.Writer) error {
	w.Write([]byte(xml.Header))
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(e); err != nil {
		return err
	}
	if err := enc.Flush(); err != nil {
		return err
	}
	if _, err := w.Write([]byte("\n")); err != nil {
		return err
	}
	// Return success
	return nil
}
