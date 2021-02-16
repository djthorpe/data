package canvas

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/dom"
	"github.com/djthorpe/data/pkg/f32"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Canvas struct {
	data.Document
	origin data.Point
	size   data.Size
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewCanvas(size data.Size, units data.Unit) data.Canvas {
	this := new(Canvas)
	this.Document = dom.NewDocumentNS("svg", data.XmlNamespaceSVG)
	this.origin = data.ZeroPoint
	this.size = size

	// Set width and height attributes
	w, h := size.UnitString(units)
	if w != "" && w != "0" {
		this.Document.SetAttr("width", w)
	}
	if h != "" && h != "0" {
		this.Document.SetAttr("height", h)
	}

	// Set the viewbox
	if err := setViewBox(this.Document, this.origin, this.size); err != nil {
		return nil
	}

	// Return the canvas
	return this
}

/////////////////////////////////////////////////////////////////////
// METHODS

// Get viewBox origin
func (this *Canvas) Origin() data.Point {
	return this.origin
}

// Get viewBox size
func (this *Canvas) Size() data.Size {
	return this.size
}

// Set viewBox
func (this *Canvas) SetViewBox(pt data.Point, sz data.Size) error {
	if err := setViewBox(this.Document, pt, sz); err != nil {
		return err
	} else {
		this.origin = pt
		this.size = sz
	}

	// Return success
	return nil
}

func (this *Canvas) DOM() data.Document {
	return this.Document
}

func (this *Canvas) Title(cdata string) data.Canvas {
	cdata = strings.TrimSpace(cdata)

	// Remove existing title tags
	if title := this.Document.GetElementsByTagNameNS("title", data.XmlNamespaceSVG); len(title) != 0 {
		for _, child := range title {
			this.Document.RemoveChild(child)
		}
	}

	// Create a new title
	if cdata != "" {
		title := this.Document.CreateElementNS("title", data.XmlNamespaceSVG)
		if err := title.AddChild(this.Document.CreateText(cdata)); err != nil {
			return nil
		} else if err := this.Document.InsertChildBefore(title, this.Document.FirstChild()); err != nil {
			return nil
		}
	}

	// Return success
	return this
}

func (this *Canvas) Version(value string) data.Canvas {
	value = strings.TrimSpace(value)

	if value != "" {
		// Set version attribute
		this.Document.SetAttr("version", value)
	} else {
		// Remove attribute
		this.Document.RemoveAttr("version")
	}

	// Return success
	return this
}

/////////////////////////////////////////////////////////////////////
// WRITE CANVAS

func (this *Canvas) Write(fmt data.Writer, w io.Writer) error {
	switch {
	case fmt&data.SVG == data.SVG:
		return this.writeSVG(fmt, w)
	default:
		return data.ErrNotImplemented
	}
}

func (this *Canvas) writeSVG(fmt data.Writer, w io.Writer) error {
	opts := data.DOMWriteDirective
	if fmt&data.Minify != data.Minify {
		opts |= data.DOMWriteIndentSpace2
	}
	return this.Document.WriteEx(w, opts)
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func setViewBox(element data.Node, origin data.Point, size data.Size) error {
	// Check parameters
	if size.W == 0 || size.H == 0 {
		return data.ErrBadParameter.WithPrefix("setViewBox")
	}

	// Update viewbox
	viewBox := f32.Join([]float32{origin.X, origin.Y, f32.Abs(size.W), f32.Abs(size.H)}, " ")
	if err := element.SetAttr("viewBox", viewBox); err != nil {
		return err
	}

	// Return success
	return nil
}

func viewBoxFromAttr(element data.Node) (data.Point, data.Size, error) {
	viewbox, exists := element.Attr("viewBox")
	if exists == false || viewbox.Value == "" {
		return data.ZeroPoint, data.ZeroSize, nil
	}
	r := strings.NewReader(viewbox.Value)
	var pt data.Point
	var sz data.Size
	if n, err := fmt.Fscanf(r, "%g %g %g %g", &pt.X, &pt.Y, &sz.W, &sz.H); err != nil {
		return data.ZeroPoint, data.ZeroSize, err
	} else if n != 4 {
		return data.ZeroPoint, data.ZeroSize, data.ErrBadParameter.WithPrefix("Invalid viewBox: ", strconv.Quote(viewbox.Value))
	}
	return pt, sz, nil
}
