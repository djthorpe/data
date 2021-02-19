package data

import (
	"fmt"
	"io"
	"math"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Point struct {
	X, Y float32
}

type Size struct {
	W, H float32
}

type (
	Unit     int
	Align    int
	LineCap  int
	LineJoin int
	FillRule int
	Writer   int
)

/////////////////////////////////////////////////////////////////////
// INTERFACES

type Canvas interface {
	// Get and set canvas viewbox
	Origin() Point
	Size() Size
	SetViewBox(Point, Size) error

	// Set canvas properties
	Title(string) Canvas
	Version(string) Canvas

	// Return canvas as an XML document
	DOM() Document

	// Write to data stream
	Write(Writer, io.Writer) error

	// Implements CanvasGroup
	CanvasGroup

	// Create a group and attach elements to group
	Group(...CanvasElement) CanvasGroup

	// Define a marker and attach elements to marker
	Marker(Point, Size, ...CanvasElement) CanvasGroup

	// Drawing primitives
	Circle(Point, float32) CanvasElement
	Ellipse(Point, Size) CanvasElement
	Line(Point, Point) CanvasElement
	Rect(Point, Size) CanvasElement
	Text(Point, ...CanvasText) CanvasElement
	Path(...CanvasPath) CanvasElement
	Polyline(...Point) CanvasElement
	Polygon(...Point) CanvasElement

	// Path primitives
	MoveTo(Point) CanvasPath
	LineTo(Point) CanvasPath
	QuadraticTo(pt, c Point) CanvasPath
	CubicTo(pt, c1, c2 Point) CanvasPath
	ClosePath() CanvasPath

	// Transform primitives
	Scale(Size) CanvasTransform
	Translate(Point) CanvasTransform
	Rotate(float32) CanvasTransform
	RotateAround(float32, Point) CanvasTransform
	SkewX(float32) CanvasTransform
	SkewY(float32) CanvasTransform

	// Style primitives
	NoFill() CanvasStyle
	NoStroke() CanvasStyle
	Fill(Color, float32) CanvasStyle
	FillRule(FillRule) CanvasStyle
	Stroke(Color, float32) CanvasStyle
	StrokeWidth(float32) CanvasStyle
	LineCap(LineCap) CanvasStyle
	LineJoin(LineJoin) CanvasStyle
	MiterLimit(float32) CanvasStyle
	UseMarker(Align, string) CanvasStyle

	/*
		// Text primitives
		Span(string) CanvasText
		FontSize(float32, Unit) CanvasStyle
		TextAnchor(TextAlign) CanvasStyle
	*/
}

type CanvasGroup interface {
	CanvasElement

	Desc(string) CanvasGroup

	// Marker orientation, when not set or zero, uses "auto"
	OrientationAngle(float32) CanvasGroup
}

type CanvasElement interface {
	Id(string) CanvasElement
	Class(string) CanvasElement
	Style(...CanvasStyle) CanvasElement
	Transform(...CanvasTransform) CanvasElement
}

type CanvasPath interface{}
type CanvasStyle interface{}
type CanvasTransform interface{}
type CanvasText interface{}

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	None Unit = iota
	PX
	CM
	MM
	IN
	PC
	PT
	EX
	EM
)

const (
	SVG    Writer = 0
	Minify Writer = (1 << iota) // Do not indent output
	// TODO: PDF, PNG, etc
)

const (
	Start Align = (1 << iota)
	Middle
	End
)

const (
	NonZero FillRule = iota
	EvenOdd
)

const (
	CapButt LineCap = iota
	CapRound
	CapSquare
)

const (
	JoinMiter LineJoin = iota
	JoinMiterClip
	JoinArcs
	JoinRound
	JoinBevel
)

var (
	ZeroSize  = Size{0, 0}
	ZeroPoint = Point{0, 0}
	NilPoint  = Point{float32(math.NaN()), float32(math.NaN())}
	NilSize   = Size{float32(math.NaN()), float32(math.NaN())}
)

var (
	// Paper sizes in mm
	A4PortraitSize      = Size{210, 297}
	A4LandscapeSize     = Size{297, 210}
	LetterPortraitSize  = Size{215.9, 279.4}
	LetterLandscapeSize = Size{279.4, 215.9}
	LegalPortraitSize   = Size{215.9, 355.6}
	LegalLandscapeSize  = Size{355.6, 215.9}
)

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (s Size) UnitString(u Unit) (string, string) {
	return fmt.Sprint(s.W, u.String()), fmt.Sprint(s.H, u.String())
}

func (a Align) String() string {
	switch a {
	case Middle:
		return "middle"
	case End:
		return "end"
	case Start:
		fallthrough
	default:
		return "start"
	}
}

func (u Unit) String() string {
	switch u {
	case PX:
		return "px"
	case CM:
		return "cm"
	case MM:
		return "mm"
	case IN:
		return "in"
	case PC:
		return "pc"
	case PT:
		return "pt"
	case EX:
		return "ex"
	case EM:
		return "em"
	default:
		return ""
	}
}

func (c LineCap) String() string {
	switch c {
	case CapRound:
		return "round"
	case CapSquare:
		return "square"
	case CapButt:
		fallthrough
	default:
		return "butt"
	}
}

func (j LineJoin) String() string {
	switch j {
	case JoinMiterClip:
		return "miter-clip"
	case JoinArcs:
		return "arcs"
	case JoinRound:
		return "round"
	case JoinBevel:
		return "bevel"
	case JoinMiter:
		fallthrough
	default:
		return "miter"
	}
}

func (r FillRule) String() string {
	switch r {
	case EvenOdd:
		return "evenodd"
	case NonZero:
		fallthrough
	default:
		return "nonzero"
	}
}
