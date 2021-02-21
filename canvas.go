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
	Unit        int
	Align       int
	Adjust      int
	LineCap     int
	LineJoin    int
	FillRule    int
	FontVariant uint32
	Writer      int
)

/////////////////////////////////////////////////////////////////////
// INTERFACES

type Canvas interface {
	// Get and set canvas viewbox
	ViewBox() (Point, Size)
	SetViewBox(Point, Size) error
	Origin() Point
	Size() Size

	// Set canvas properties
	Title(string) Canvas
	Version(string) Canvas

	// Return canvas as an XML document
	DOM() Document

	// Remove a group or element
	Remove(...CanvasElement) error

	// Write to data stream
	Write(Writer, io.Writer) error

	// Create a defs group and attach elements to the group
	Defs(...CanvasElement) CanvasGroup

	// Create a group and attach elements to group
	Group(...CanvasElement) CanvasGroup

	// Define a marker and attach elements to marker
	Marker(Point, Size, ...CanvasElement) CanvasGroup

	// Drawing primitives
	Circle(Point, float32) CanvasElement
	Ellipse(Point, Size) CanvasElement
	Line(Point, Point) CanvasElement
	Rect(Point, Size) CanvasElement
	Path(...CanvasPath) CanvasElement
	Polyline(...Point) CanvasElement
	Polygon(...Point) CanvasElement
	Text(Point, bool, ...CanvasText) CanvasElement
	Image(Point, Size, string) CanvasElement

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

	// Fill styles
	NoFill() CanvasStyle
	Fill(Color, float32) CanvasStyle
	FillRule(FillRule) CanvasStyle

	// Stroke styles
	NoStroke() CanvasStyle
	Stroke(Color, float32) CanvasStyle
	StrokeWidth(float32) CanvasStyle
	LineCap(LineCap) CanvasStyle
	LineJoin(LineJoin) CanvasStyle
	MiterLimit(float32) CanvasStyle

	// Text styles
	FontSize(float32, Unit) CanvasStyle
	FontFamily(string) CanvasStyle
	FontVariant(FontVariant) CanvasStyle
	TextAnchor(Align) CanvasStyle

	// Other styles
	UseMarker(Align, string) CanvasStyle

	// Text primitives
	TextSpan(string) CanvasText
	TextPath(string) CanvasText // TODO
}

type CanvasGroup interface {
	CanvasElement

	// Description element for the group
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

type CanvasText interface {
	Offset(Point) CanvasText
	Length(float32, Adjust) CanvasText
}

type CanvasPath interface{}
type CanvasStyle interface{}
type CanvasTransform interface{}

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
	Spacing Adjust = iota
	SpacingAndGlyphs
)

const (
	Thin       FontVariant = (1 << iota) // 100
	ExtraLight                           // 200
	Light                                // 300
	Regular                              // 400 (same as normal)
	Medium                               // 500
	SemiBold                             // 600
	Bold                                 // 700
	ExtraBold                            // 800
	Black                                // 900
	Bolder
	Lighter
	Italic
	Oblique
	Normal = Regular
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

func (a Adjust) String() string {
	switch a {
	case SpacingAndGlyphs:
		return "spacingAndGlyphs"
	case Spacing:
		fallthrough
	default:
		return "spacing"
	}
}

func (v FontVariant) String() string {
	switch v {
	case Thin:
		return "100"
	case ExtraLight:
		return "200"
	case Light:
		return "300"
	case Medium:
		return "500"
	case SemiBold:
		return "600"
	case Bold:
		return "bold" // aka 700
	case ExtraBold:
		return "800"
	case Black:
		return "900"
	case Italic:
		return "italic"
	case Oblique:
		return "oblique"
	case Lighter:
		return "lighter"
	case Bolder:
		return "bolder"
	case Regular:
		fallthrough
	default:
		return "normal" // aka 400
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
