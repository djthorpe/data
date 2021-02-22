package main

/*
	Create an SVG file which includes colors within a swatch
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

	// Gray background
	c.Rect(c.Origin(), c.Size()).Style(
		c.Fill(color.Ivory, 1.0),
	)

	// Determine how many items across and down assuming
	// a square
	palette := color.Palette(data.ColorAll)
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
	for i, col := range palette {
		// Check how far color is from black to set foreground
		text := color.DarkSlateGray
		if color.Distance(col, color.Black) < 450 {
			text = color.White
		}
		// Set origin and size for the widget
		origin := data.Point{
			X: float32(i%int(across)) * size.W,
			Y: float32(i/int(across)) * size.H,
		}
		// Add color widget to the canvas
		AddColor(c, size, col, text).Transform(
			c.Translate(origin),
			c.Scale(data.Size{0.95, 0.95}),
		)
	}

	// Output as SVG
	c.Write(data.SVG, os.Stdout)
}

func AddColor(c data.Canvas, size data.Size, bg data.Color, fg data.Color) data.CanvasGroup {
	name := color.Name(bg)
	value := color.HashString(bg)
	return c.Group(
		c.Rect(data.ZeroPoint, size).Style(
			c.NoStroke(),
			c.Fill(bg, 1.0),
		),
		c.Rect(data.ZeroPoint, size).Style(
			c.NoFill(),
			c.Stroke(color.LightGray, 1.0),
			c.StrokeWidth(0.1),
		),
		c.Group(
			c.Text(
				geom.CentrePoint(data.ZeroPoint, size),
				false,
				c.TextSpan(name),
			),
			c.Text(
				geom.CentrePoint(data.ZeroPoint, size),
				false,
				c.TextSpan(value).Offset(data.Point{0, 3}),
			),
		).Style(
			c.TextAnchor(data.Middle),
			c.FontFamily("Arial, Helvetica, sans-serif"),
			c.FontSize(3, data.None),
			c.Fill(fg, 1.0),
		),
	)
}
