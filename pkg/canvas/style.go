package canvas

import (
	"fmt"
	"strings"

	"github.com/djthorpe/data"
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
	Rule    data.FillRule
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
	styleFillRule
	styleNone styleop = 0
	styleMin          = styleFillNone
	styleMax          = styleFillRule
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

func (e *Element) FontSize(size float32, unit data.Unit) data.CanvasStyle {
	return &styledef{Op: styleFontSize, Width: size, Unit: unit}
}

func (e *Element) TextAnchor(align data.TextAlign) data.CanvasStyle {
	return &styledef{Op: styleTextAnchor, Align: align}
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
	case styleFillRule:
		if _, exists := s.defs[styleFillNone]; exists == false {
			return StyleString(op, def.Rule)
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

func StyleString(name styleop, value interface{}) string {
	return fmt.Sprintf("%v: %v", name, value)
}
