package canvas

import (
	"io"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/dom"
	"github.com/djthorpe/data/pkg/f32"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Canvas struct {
	data.Document
	origin  data.Point
	size    data.Size
	title   data.Node
	version data.Node
}

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DOMOptions = data.DOMWriteDirective | data.DOMWriteIndentSpace2
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewCanvas(size data.Size, units data.Unit) data.Canvas {
	this := new(Canvas)
	this.Document = dom.NewDocumentNS("svg", data.XmlNamespaceSVG, DOMOptions)
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
	this.setViewBox(this.Document)

	// Return the canvas
	return this
}

/////////////////////////////////////////////////////////////////////
// METHODS

// Get and set canvas viewbox
func (this *Canvas) Origin() data.Point {
	return this.origin
}

func (this *Canvas) Size() data.Size {
	return this.size
}

func (this *Canvas) SetViewBox(pt data.Point, sz data.Size) error {
	this.origin = pt
	this.size = sz
	if err := this.setViewBox(this.Document); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *Canvas) Title(cdata string) data.Canvas {
	// Remove any existing title tag
	if this.title != nil {
		this.Document.RemoveChild(this.title)
	}

	// Create a new title
	this.title = this.Document.CreateElement("title")
	if err := this.title.AddChild(this.Document.CreateText(cdata)); err != nil {
		return nil
	}

	// Add title to document
	// TODO
	if err := this.Document.AddChild(this.title); err != nil {
		return nil
	}

	// Return success
	return this
}

func (this *Canvas) Version(value string) data.Canvas {
	// Set version attribute
	this.Document.SetAttr("version", value)

	// Return success
	return this
}

/////////////////////////////////////////////////////////////////////
// WRITE CANVAS

func (this *Canvas) Write(fmt data.Writer, w io.Writer) error {
	switch fmt {
	case data.SVG:
		return this.writeSVG(w)
	default:
		return data.ErrNotImplemented
	}
}

func (this *Canvas) writeSVG(w io.Writer) error {
	return this.Document.Write(w)
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Canvas) setViewBox(element data.Node) error {
	// Check parameters
	if this.size.W == 0 || this.size.H == 0 {
		return data.ErrBadParameter.WithPrefix("setViewBox")
	}

	// Update viewbox
	viewBox := f32.Join([]float32{this.origin.X, this.origin.Y, f32.Abs(this.size.W), f32.Abs(this.size.H)}, " ")
	if err := element.SetAttr("viewBox", viewBox); err != nil {
		return err
	}

	// Return success
	return nil
}

/*
	// Drawing primitives
	Circle(Point, float32) CanvasElement
	Ellipse(Point, Size) CanvasElement
	Line(Point, Point) CanvasElement
	Rect(Point, Size) CanvasElement
	Text(Point, ...CanvasText) CanvasElement
	Path([]Point) CanvasPath

	// Transform primitives
	Scale(Size) CanvasTransform
	Translate(Point) CanvasTransform
	Rotate(float32) CanvasTransform
	RotateAround(float32, Point) CanvasTransform
	SkewX(float32) CanvasTransform
	SkewY(float32) CanvasTransform

	// Style primitives
	Fill(Color, float32) CanvasStyle
	NoFill() CanvasStyle
	FillRule(FillRule) CanvasStyle
	Stroke(Color, float32) CanvasStyle
	StrokeWidth(float32) CanvasStyle
	NoStroke() CanvasStyle
	FontSize(float32, Unit) CanvasStyle
	TextAnchor(TextAlign) CanvasStyle
	LineCap(LineCap) CanvasStyle
	LineJoin(LineJoin) CanvasStyle
	MiterLimit(float32) CanvasStyle

	// Text primitives
	Span(string) CanvasText
*/
