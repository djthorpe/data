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
	Unit    data.Unit
	Align   data.TextAlign
	Cap     data.LineCap
	Join    data.LineJoin
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
	styleFontSize
	styleTextAnchor
	styleLineCap
	styleLineJoin
	styleMiterLimit
	styleNone styleop = 0
	styleMin          = styleFillNone
	styleMax          = styleMiterLimit
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
	if width == 0 {
		return e.NoStroke()
	} else {
		return &styledef{Op: styleStrokeWidth, Width: width}
	}
}

func (e *Element) FontSize(size float32, unit data.Unit) data.CanvasStyle {
	return &styledef{Op: styleFontSize, Width: size, Unit: unit}
}

func (e *Element) TextAnchor(align data.TextAlign) data.CanvasStyle {
	return &styledef{Op: styleTextAnchor, Align: align}
}

func (e *Element) LineCap(cap data.LineCap) data.CanvasStyle {
	return &styledef{Op: styleLineCap, Cap: cap}
}

func (e *Element) LineJoin(join data.LineJoin) data.CanvasStyle {
	return &styledef{Op: styleLineJoin, Join: join}
}

func (e *Element) MiterLimit(limit float32) data.CanvasStyle {
	return &styledef{Op: styleMiterLimit, Width: limit}
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
	case styleFontSize:
		return "font-size"
	case styleTextAnchor:
		return "text-anchor"
	case styleLineCap:
		return "stroke-linecap"
	case styleLineJoin:
		return "stroke-linejoin"
	case styleMiterLimit:
		return "stroke-miterlimit"
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
			return FloatString(op, def.Width)
		}
	case styleFontSize:
		return UnitString(op, def.Width, def.Unit)
	case styleTextAnchor:
		return StyleString(op, def.Cap)
	case styleLineCap:
		if _, exists := s.defs[styleStrokeNone]; exists == false {
			return StyleString(op, def.Cap)
		}
	case styleLineJoin:
		if _, exists := s.defs[styleStrokeNone]; exists == false {
			return StyleString(op, def.Cap)
		}
	case styleMiterLimit:
		if _, exists := s.defs[styleStrokeNone]; exists == false {
			return FloatString(op, def.Width)
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

func UnitString(name styleop, value float32, unit data.Unit) string {
	return fmt.Sprintf("%v: %s%s", name, f32.String(value), unit.String())
}

func StyleString(name styleop, value interface{}) string {
	return fmt.Sprintf("%v: %v", name, value)
}
