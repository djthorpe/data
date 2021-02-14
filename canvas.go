package data

import (
	"fmt"
	"io"
	"math"

	f32 "github.com/djthorpe/data/pkg/f32"
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
	Unit      int
	TextAlign int
	LineCap   int
	LineJoin  int
)

/////////////////////////////////////////////////////////////////////
// INTERFACES

type Canvas interface {
	CanvasGroup

	// Get canvas properties
	Origin() Point
	Size() Size

	// Set canvas properties
	Title(string) Canvas
	Version(string) Canvas

	// Write SVG
	Write(w io.Writer) error

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
}

type CanvasElement interface {
	Id(string) CanvasElement
	Class(string) CanvasElement
	Style(...CanvasStyle) CanvasElement
	Transform(...CanvasTransform) CanvasElement
}

type CanvasGroup interface {
	CanvasElement

	Desc(string) CanvasGroup
	Group(...CanvasElement) CanvasGroup
}

type CanvasPath interface {
	CanvasElement

	MoveTo(Point) CanvasPath
	LineTo(Point) CanvasPath
	QuadraticTo(pt, c Point) CanvasPath
	CubicTo(pt, c1, c2 Point) CanvasPath
	ClosePath() CanvasPath
}

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
	Start TextAlign = iota
	Middle
	End
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
	A4PortraitSize  = Size{594, 841}
	A4LandscapeSize = Size{841, 594}
)

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (s Size) UnitString(u Unit) (string, string) {
	return fmt.Sprint(s.W, u.String()), fmt.Sprint(s.H, u.String())
}

func (ta TextAlign) String() string {
	switch ta {
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

/////////////////////////////////////////////////////////////////////
// FUNCTIONS

func (p Point) IsNil() bool {
	return f32.IsNaN(p.X) || f32.IsNaN(p.Y)
}
