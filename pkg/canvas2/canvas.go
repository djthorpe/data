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

func (this *Canvas) Title(cdata string) data.Canvas {
	// Get title tag
	title := this.Document.GetElementsByTagNameNS("title", data.XmlNamespaceSVG)
	if len(title) == 0 {
		title = []data.Node{this.Document.CreateElementNS("title", data.XmlNamespaceSVG)}
	}
	// Set title cdata

	// Create a new title
	this.title
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