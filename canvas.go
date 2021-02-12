package data

import (
	"fmt"
	"io"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Point struct {
	X, Y float32
}

type Size struct {
	W, H float32
}

type Color struct {
	R, G, B uint8
}

type Unit int

/////////////////////////////////////////////////////////////////////
// INTERFACES

type Canvas interface {
	CanvasGroup

	// Set canvas properties
	Title(string) Canvas
	Version(string) Canvas

	// Write SVG
	Write(w io.Writer) error

	// Drawing primitives
	Circle(Point, float32) CanvasElement
	Ellipse(Point, Size) CanvasElement
	Path([]Point) CanvasElement
}

type CanvasGroup interface {
	CanvasElement
	Desc(string) CanvasGroup
	Group(...CanvasElement) CanvasGroup
}

type CanvasElement interface {
	Id(string)
	Class(string)
}

type CanvasStyle interface{}

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

var (
	ZeroSize  = Size{0, 0}
	ZeroPoint = Point{0, 0}
)

var (
	A4PortraitSize  = Size{594, 841}
	A4LandscapeSize = Size{841, 594}
)

var (
	White = Color{0xFF, 0xFF, 0xFF}
	Black = Color{0x00, 0x00, 0x00}
)

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (s Size) UnitString(u Unit) (string, string) {
	return fmt.Sprint(s.W, u.String()), fmt.Sprint(s.H, u.String())
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
