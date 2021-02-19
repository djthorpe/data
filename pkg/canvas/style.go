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
	defs map[styleop]string
}

type styleop uint

type styledef struct {
	Op      styleop
	Color   data.Color
	Opacity float32
	Width   float32
	Unit    data.Unit
	Align   data.Align
	Cap     data.LineCap
	Join    data.LineJoin
	Rule    data.FillRule
	Uri     string
}

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	fillNone styleop = (1 << iota)
	strokeNone
	fillColor
	fillOpacity
	strokeColor
	strokeOpacity
	strokeWidth
	fontSize
	textAnchor
	lineCap
	lineJoin
	miterLimit
	fillRule
	markerStart
	markerMid
	markerEnd
	styleNone styleop = 0
	styleMin          = fillNone
	styleMax          = markerEnd
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewStyles(defs []data.CanvasStyle) *Style {
	this := new(Style)
	this.defs = make(map[styleop]string, len(defs))

	// Add styles
	for _, def := range defs {
		if def == nil {
			return nil
		} else if def, ok := def.(*styledef); ok == false {
			return nil
		} else if err := this.setDef(def); err != nil {
			return nil
		}
	}

	return this
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Style) setDef(def *styledef) error {
	for v := styleMin; v <= styleMax; v <<= 1 {
		if def.Op&v != v {
			continue
		}
		if str, err := v.StyleString(def); err != nil {
			return err
		} else if str != "" {
			this.defs[v] = str
		}
	}
	// Return success
	return nil
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (*Canvas) NoFill() data.CanvasStyle {
	return &styledef{Op: fillNone}
}

func (*Canvas) NoStroke() data.CanvasStyle {
	return &styledef{Op: strokeNone}
}

func (*Canvas) Fill(color data.Color, opacity float32) data.CanvasStyle {
	return &styledef{Op: fillColor | fillOpacity, Color: color, Opacity: opacity}
}

func (*Canvas) FillRule(rule data.FillRule) data.CanvasStyle {
	return &styledef{Op: fillRule, Rule: rule}
}

func (*Canvas) Stroke(color data.Color, opacity float32) data.CanvasStyle {
	return &styledef{Op: strokeColor | strokeOpacity, Color: color, Opacity: opacity}
}

func (*Canvas) StrokeWidth(width float32) data.CanvasStyle {
	if width == 0 {
		return &styledef{Op: strokeNone}
	} else {
		return &styledef{Op: strokeWidth, Width: width}
	}
}

func (*Canvas) LineCap(cap data.LineCap) data.CanvasStyle {
	return &styledef{Op: lineCap, Cap: cap}
}

func (*Canvas) LineJoin(join data.LineJoin) data.CanvasStyle {
	return &styledef{Op: lineJoin, Join: join}
}

func (*Canvas) MiterLimit(limit float32) data.CanvasStyle {
	return &styledef{Op: miterLimit, Width: limit}
}

func (*Canvas) UseMarker(pos data.Align, uri string) data.CanvasStyle {
	// TODO:
	// if uri is #?([a-zA-Z\-]+[a-zA-Z0-9\-]*) (ie, an id) then
	// wrap it in url(#id)
	op := markerStart | markerMid | markerEnd
	if pos != 0 && pos&data.Start == 0 {
		op ^= markerStart
	}
	if pos != 0 && pos&data.Middle == 0 {
		op ^= markerMid
	}
	if pos != 0 && pos&data.End == 0 {
		op ^= markerEnd
	}
	return &styledef{Op: op, Uri: uri}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (f styleop) String() string {
	switch f {
	case styleNone:
		return "none"
	case fillNone, fillColor:
		return "fill"
	case fillRule:
		return "fill-rule"
	case strokeNone, strokeColor:
		return "stroke"
	case fillOpacity:
		return "fill-opacity"
	case strokeOpacity:
		return "stroke-opacity"
	case strokeWidth:
		return "stroke-width"
	case fontSize:
		return "font-size"
	case textAnchor:
		return "text-anchor"
	case lineCap:
		return "stroke-linecap"
	case lineJoin:
		return "stroke-linejoin"
	case miterLimit:
		return "stroke-miterlimit"
	case markerStart:
		return "marker-start"
	case markerMid:
		return "marker-mid"
	case markerEnd:
		return "marker-end"
	default:
		return "[?? invalid styleop value]"
	}
}

func (f styleop) StyleString(args *styledef) (string, error) {
	switch f {
	case fillNone, strokeNone:
		return fmt.Sprint(f, ":none;"), nil
	case fillColor, strokeColor:
		return fmt.Sprint(f, ":", color.String(args.Color), ";"), nil
	case fillOpacity, strokeOpacity:
		return fmt.Sprint(f, ":", f32.String(args.Opacity), ";"), nil
	case strokeWidth, miterLimit:
		return fmt.Sprint(f, ":", f32.String(args.Width), ";"), nil
	case lineCap:
		return fmt.Sprint(f, ":", args.Cap, ";"), nil
	case lineJoin:
		return fmt.Sprint(f, ":", args.Join, ";"), nil
	case fillRule:
		return fmt.Sprint(f, ":", args.Rule, ";"), nil
	case markerStart, markerMid, markerEnd:
		return fmt.Sprint(f, ":", args.Uri, ";"), nil
	default:
		return "", data.ErrBadParameter.WithPrefix("SetStyle: ", f)
	}
}

func (s *Style) String() string {
	attrs := make([]string, 0, len(s.defs))
	for v := styleMin; v <= styleMax; v <<= 1 {
		if style, exists := s.defs[v]; exists {
			attrs = append(attrs, style)
		}
	}
	return strings.Join(attrs, " ")
}

func UnitString(name styleop, value float32, unit data.Unit) string {
	return fmt.Sprintf("%v: %s%s", name, f32.String(value), unit.String())
}
