package canvas

import (
	"fmt"
	"strings"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/color"
	"github.com/djthorpe/data/pkg/f32"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Style struct {
	defs map[styleop]*styledef
}

type styleop uint

type styledef struct {
	Op      styleop
	Color   data.Color
	Opacity float32
	Width   float32
}

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	styleFillNone styleop = (1 << iota)
	styleStrokeNone
	styleFillColor
	styleFillOpacity
	styleStrokeColor
	styleStrokeOpacity
	styleStrokeWidth
	styleNone styleop = 0
	styleMin          = styleFillColor
	styleMax          = styleStrokeWidth
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewStyle(styles []data.CanvasStyle) *Style {
	this := new(Style)
	this.defs = make(map[styleop]*styledef)

	// Append styles
	this.Append(styles)

	// Return success
	return this
}

/////////////////////////////////////////////////////////////////////
// STYLE METHODS

func (this *Style) Append(styles []data.CanvasStyle) {
	for _, style := range styles {
		if style_, ok := style.(*styledef); ok {
			this.setStyle(style_)
		}
	}
}

func (this *Style) setStyle(style *styledef) {
	for f := styleMin; f <= styleMax; f <<= 1 {
		if style.Op&f == f {
			this.defs[f] = style
		}
	}
}

/////////////////////////////////////////////////////////////////////
// ELEMENT METHODS

func (e *Element) NoFill() data.CanvasStyle {
	return &styledef{Op: styleFillNone}
}

func (e *Element) NoStroke() data.CanvasStyle {
	return &styledef{Op: styleStrokeNone}
}

func (e *Element) Fill(color data.Color, opacity float32) data.CanvasStyle {
	return &styledef{Op: styleFillColor | styleFillOpacity, Color: color, Opacity: opacity}
}

func (e *Element) Stroke(color data.Color, opacity float32) data.CanvasStyle {
	return &styledef{Op: styleStrokeColor | styleStrokeOpacity, Color: color, Opacity: opacity}
}

func (e *Element) StrokeWidth(width float32) data.CanvasStyle {
	return &styledef{Op: styleStrokeWidth, Width: width}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (f styleop) FlagString() string {
	switch f {
	case styleNone:
		return "none"
	case styleFillNone, styleFillColor:
		return "fill"
	case styleStrokeNone, styleStrokeColor:
		return "stroke"
	case styleFillOpacity:
		return "fill-opacity"
	case styleStrokeOpacity:
		return "stroke-opacity"
	case styleStrokeWidth:
		return "stroke-width"
	default:
		return "[?? invalid styleop value]"
	}
}

func (f styleop) String() string {
	str := ""
	if f == 0 {
		return f.FlagString()
	}
	for v := styleMin; v <= styleMax; v <<= 1 {
		if f&v == v {
			str += v.FlagString() + "|"
		}
	}
	return strings.TrimSuffix(str, "|")
}

func (s *Style) String() string {
	str := ""
	for f := styleMin; f <= styleMax; f <<= 1 {
		if def, exists := s.defs[f]; exists {
			if value := s.DefString(f, def); value != "" {
				str += value + "; "
			}
		}
	}
	return strings.TrimSuffix(str, " ")
}

func (s *Style) DefString(op styleop, def *styledef) string {
	switch op {
	case styleFillNone:
		return "fill: none"
	case styleStrokeNone:
		return "stroke: none"
	case styleFillColor:
		if _, exists := s.defs[styleFillNone]; exists == false {
			return ColorString(op, def.Color)
		}
	case styleFillOpacity:
		if _, exists := s.defs[styleFillNone]; exists == false {
			return FloatString(op, def.Opacity)
		}
	case styleStrokeColor:
		if _, exists := s.defs[styleStrokeNone]; exists == false {
			return ColorString(op, def.Color)
		}
	case styleStrokeOpacity:
		if _, exists := s.defs[styleStrokeNone]; exists == false {
			return FloatString(op, def.Opacity)
		}
	case styleStrokeWidth:
		if _, exists := s.defs[styleStrokeNone]; exists == false {
			return FloatString(op, def.Opacity)
		}
	}
	// By default return empty string
	return ""
}

func ColorString(name styleop, value data.Color) string {
	return fmt.Sprintf("%v: %v", name, color.String(value))
}

func FloatString(name styleop, value float32) string {
	return fmt.Sprintf("%v: %s", name, f32.String(value))
}
