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
	c := NewElement("svg", "")
	c.XMLName.Space = xmlNs
	w, h := size.UnitString(units)
	if w != "" && w != "0" {
		c.Attr("width", w)
	}
	if h != "" && h != "0" {
		c.Attr("height", h)
	}
	return c
}

/////////////////////////////////////////////////////////////////////
// METHODS

func (e *Element) Title(value string) data.Canvas {
	e.addChild(NewElement("title", value))
	return e
}

func (e *Element) Version(value string) data.Canvas {
	e.Attr("version", value)
	return e
}

func (e *Element) Circle(centre data.Point, radius float32) data.CanvasElement {
	c := NewElement("circle", "")
	c.Attr("cx", centre.X)
	c.Attr("cy", centre.Y)
	c.Attr("r", radius)
	e.addChild(c)
	return e
}

func (e *Element) Ellipse(centre data.Point, radius data.Size) data.CanvasElement {
	c := NewElement("ellipse", "")
	c.Attr("cx", centre.X)
	c.Attr("cy", centre.Y)
	c.Attr("rx", radius.W)
	c.Attr("ry", radius.H)
	e.addChild(c)
	return e
}

func (e *Element) Path(pts []data.Point) data.CanvasElement {
	// We require at least two elements for a path
	if len(pts) < 2 {
		return nil
	}
	c := NewElement("path", "")
	attr := ""
	for i, pt := range pts {
		if i == 0 {
			attr += fmt.Sprint("M ", pt.X, pt.Y)
		} else {
			attr += fmt.Sprint(" L ", pt.X, pt.Y)
		}
	}
	c.Attr("d", attr)
	e.addChild(c)
	fmt.Println(c)
	return e
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
