package main

/*
	Create an SVG file which includes all the named colors
	on an A4 Landscape size. Demonstrates creating SVG files
	and calculating numerical difference between two colors
	in order to ensure text visibility
*/

import (
	"os"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/canvas"
	"github.com/djthorpe/data/pkg/color"
	"github.com/djthorpe/data/pkg/f32"
	"github.com/djthorpe/data/pkg/geom"
)

func main() {
	c := canvas.NewCanvas(data.A4LandscapeSize, data.MM)

	// Top of document
	c.Title("Color Palette")
	c.Desc("Display all colors in the palette")

	// Gray background
	c.Rect(c.Origin(), c.Size()).Style(
		c.Fill(color.Ivory, 1.0),
	)

	// Determine how many items across and down assuming
	// a square
	palette := color.Palette()
	count := float32(len(palette))
	across := f32.Floor(f32.Sqrt(count))
	down := f32.Ceil(count / across)
	var size data.Size
	for {
		size = geom.DivideSize(c.Size(), across)
		expected := data.Size{size.W * across, size.H * down}
		if expected.W > c.Size().W || expected.H > c.Size().H {
			across++
		} else {
			break
		}
	}
	for i, cr := range palette {
		var fill data.CanvasStyle
		if color.Distance(cr, color.Black) < 400 {
			fill = c.Fill(color.White, 1.0)
		} else {
			fill = c.Fill(color.Black, 1.0)
		}
		origin := data.Point{
			X: float32(i%int(across)) * size.W,
			Y: float32(i/int(across)) * size.H,
		}
		c.Group(
			c.Rect(data.ZeroPoint, size).Style(
				c.NoStroke(),
				c.Fill(cr, 1.0),
			),
			c.Rect(data.ZeroPoint, size).Style(
				c.NoFill(),
				c.Stroke(color.LightGray, 1.0),
				c.StrokeWidth(0.5),
			),
			c.Text(
				geom.CentrePoint(data.ZeroPoint, size),
				c.Span(color.Name(cr)),
			).Style(
				c.TextAnchor(data.Middle),
				c.FontSize(9, data.PT),
				fill,
			),
		).Transform(
			c.Translate(origin),
			c.Scale(data.Size{0.95, 0.95}),
		)
	}

	// Output as SVG
	c.Write(os.Stdout)
}
