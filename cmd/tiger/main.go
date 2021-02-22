package main

/*
	Create an SVG file from tiger data file, which is stored
	in etc/datasets/tiger_data.txt in the repository
*/

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/canvas"
	"github.com/djthorpe/data/pkg/geom"
)

///////////////////////////////////////////////////////////////////////////////

var (
	reOpcode = regexp.MustCompile("'(\\w)'")
	reValue  = regexp.MustCompile("([0-9\\.]*[0-9]+)f?")
)

var (
	flagRotate = flag.Float64("rotate", 0, "Rotation angle")
)

///////////////////////////////////////////////////////////////////////////////

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "Missing filename argument")
		os.Exit(-1)
	}

	// Create canvas
	c := canvas.NewCanvas(data.LetterPortraitSize, data.MM)
	c.Title("Tiger")

	// Set size in points 1mm = 0.3527pt (72 points per inch)
	c.SetViewBox(data.ZeroPoint, geom.DivideSize(c.Size(), 0.3527))

	// Create a group for rotation
	g := c.Group()
	g.Transform(
		c.RotateAround(float32(*flagRotate), geom.CentrePoint(c.Origin(), c.Size())),
	)

	// Read instructions from files
	for _, path := range flag.Args() {
		if err := ParseFile(c, g, path); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(-1)
		}
	}

	// Output as SVG
	c.Write(data.SVG, os.Stdout)
}

///////////////////////////////////////////////////////////////////////////////

func ParseFile(c data.Canvas, g data.CanvasGroup, path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// Read opcodes and values
	if opcodes := reOpcode.FindAllSubmatch(data, -1); opcodes == nil {
		return errors.New("Invalid data file, no opcodes")
	} else if values := reValue.FindAllSubmatch(data, -1); values == nil {
		return errors.New("Invalid data file, no values")
	} else if values, err := ParseValues(values); err != nil {
		return err
	} else if err := Render(c, g, opcodes, values); err != nil {
		return err
	}

	// Return success
	return nil
}

func ParseValues(bytes [][][]byte) ([]float32, error) {
	floats := make([]float32, len(bytes))
	for i, value := range bytes {
		if value, err := strconv.ParseFloat(string(value[1]), 32); err != nil {
			return nil, err
		} else {
			floats[i] = float32(value)
		}
	}

	// Return success
	return floats, nil
}

func Render(c data.Canvas, g data.CanvasGroup, opcodes [][][]byte, values []float32) error {
	i := 0
	v := 0
	for i < len(opcodes) && v < len(values) {
		op := new(Operation)

		// Fill opcode
		if err := op.Fill(c, string(opcodes[i][1])); err != nil {
			return err
		}
		i++

		// Stroke opcode
		if err := op.Stroke(c, string(opcodes[i][1])); err != nil {
			return err
		}
		i++

		// Linecap opcode
		if err := op.LineCap(c, string(opcodes[i][1])); err != nil {
			return err
		}
		i++

		// Linejoin opcode
		if err := op.LineJoin(c, string(opcodes[i][1])); err != nil {
			return err
		}
		i++

		// Miter Limit value
		if err := op.MiterLimit(c, values[v]); err != nil {
			return err
		}
		v++

		// Stroke Width value
		if err := op.StrokeWidth(c, values[v]); err != nil {
			return err
		}
		v++

		// Stroke color
		if err := op.StrokeColor(c, values[v], values[v+1], values[v+2]); err != nil {
			return err
		}
		v += 3

		// Fill color
		if err := op.FillColor(c, values[v], values[v+1], values[v+2]); err != nil {
			return err
		}
		v += 3

		// Path elements
		count := int(values[v])
		if err := op.CreatePath(c, count); err != nil {
			return err
		}
		v++

		segments := make([]data.CanvasPath, count)
		for j := 0; j < count; j++ {
			if segment, vinc, err := op.AddPoints(c, string(opcodes[i][1]), values, v); err != nil {
				return err
			} else {
				segments[j] = segment
				v += vinc
			}
			i++
		}

		// Add path to group
		g.Append(c.Path(segments...).Style(op.style...))
	}

	// Return success
	return nil
}
