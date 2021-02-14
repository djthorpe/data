package main

import (
	"fmt"
	"os"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/color"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Operation struct {
	style []data.CanvasStyle
	path  data.CanvasPath
}

///////////////////////////////////////////////////////////////////////////////
// METHODS

func (op *Operation) Fill(c data.Canvas, code string) error {
	switch code {
	case "N":
		op.style = append(op.style, c.NoFill())
	case "F":
		op.style = append(op.style, c.Fill(color.White, 1.0))
	case "E":
		op.style = append(op.style, c.Fill(color.White, 1.0))
		fmt.Fprintln(os.Stderr, "TODO:EVENODD Fill Rule")
	default:
		return fmt.Errorf("Invalid Fill opcode %q", code)
	}
	return nil
}

func (op *Operation) Stroke(c data.Canvas, code string) error {
	switch code {
	case "N":
		op.style = append(op.style, c.NoStroke())
	case "S":
		op.style = append(op.style, c.Stroke(color.White, 1.0))
	default:
		return fmt.Errorf("Invalid Stroke opcode %q", code)
	}
	return nil
}

func (op *Operation) LineCap(c data.Canvas, code string) error {
	switch code {
	case "B":
		op.style = append(op.style, c.LineCap(data.CapButt))
	case "R":
		op.style = append(op.style, c.LineCap(data.CapRound))
	case "S":
		op.style = append(op.style, c.LineCap(data.CapSquare))
	default:
		return fmt.Errorf("Invalid LineCap opcode %q", code)
	}
	return nil
}

func (op *Operation) LineJoin(c data.Canvas, code string) error {
	switch code {
	case "M":
		op.style = append(op.style, c.LineJoin(data.JoinMiter))
	case "R":
		op.style = append(op.style, c.LineJoin(data.JoinRound))
	case "B":
		op.style = append(op.style, c.LineJoin(data.JoinBevel))
	default:
		return fmt.Errorf("Invalid LineJoin opcode %q", code)
	}
	return nil
}

func (op *Operation) MiterLimit(c data.Canvas, value float32) error {
	op.style = append(op.style, c.MiterLimit(value))
	return nil
}

func (op *Operation) StrokeWidth(c data.Canvas, value float32) error {
	op.style = append(op.style, c.StrokeWidth(value))
	return nil
}

func (op *Operation) StrokeColor(c data.Canvas, r, g, b float32) error {
	color := data.Color{uint8(r * 255.0), uint8(g * 255.0), uint8(b * 255.0)}
	op.style = append(op.style, c.Stroke(color, 1.0))
	return nil
}

func (op *Operation) FillColor(c data.Canvas, r, g, b float32) error {
	color := data.Color{uint8(r * 255.0), uint8(g * 255.0), uint8(b * 255.0)}
	op.style = append(op.style, c.Fill(color, 1.0))
	return nil
}

func (op *Operation) CreatePath(c data.Canvas, cap int) error {
	op.path = c.Path([]data.Point{})
	return nil
}

func (op *Operation) AddPoints(c data.Canvas, code string, values []float32, i int) (int, error) {
	switch code {
	case "M":
		op.path.MoveTo(data.Point{values[i], values[i+1]})
		return 2, nil
	case "L":
		op.path.LineTo(data.Point{values[i], values[i+1]})
		return 2, nil
	case "C":
		op.path.CubicTo(
			data.Point{values[i+4], values[i+5]},
			data.Point{values[i], values[i+1]},
			data.Point{values[i+2], values[i+3]},
		)
		return 6, nil
	case "E":
		op.path.ClosePath()
		return 0, nil
	default:
		return 0, fmt.Errorf("Invalid AddPoints opcode %q", code)
	}
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (op *Operation) String() string {
	str := "<op"
	for _, style := range op.style {
		str += " " + fmt.Sprint(style)
	}
	if op.path != nil {
		str += " " + fmt.Sprint(op.path)
	}
	return str + ">"
}
