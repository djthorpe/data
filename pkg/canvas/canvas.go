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
	c.Attr("viewBox", fmt.Sprintf("0 0 %v %v", size.W, size.H))

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

/////////////////////////////////////////////////////////////////////
// WRITE CANVAS

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
