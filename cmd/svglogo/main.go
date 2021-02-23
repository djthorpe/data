package main

/*
	Create the logo
*/

import (
	"os"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/canvas"
	"github.com/djthorpe/data/pkg/color"
	"github.com/djthorpe/data/pkg/geom"
)

func main() {
	c := canvas.NewCanvas(data.Size{200, 200}, data.PX)

	// Top of document
	c.Title("Data logo")

	// Black background
	c.Rect(c.Origin(), c.Size()).Style(
		c.Fill(color.Black, 1.0),
	)

	// Divider
	c.Path(
		c.MoveTo(geom.AddPoint(c.ViewBox())),
		c.LineTo(c.Origin()),
		c.ClosePath(),
	).Style(
		c.Stroke(color.DarkGray, 1.0),
	).Transform(
		c.RotateAround(90, geom.CentrePoint(c.ViewBox())),
	)

	// Text group
	c.Group(
		c.Text(
			geom.AddPoint(geom.CentrePoint(c.ViewBox()), data.Size{-8, 25}),
			false,
			c.TextSpan("d"),
			c.TextSpan("a").Offset(data.Point{-58, 21}),
		),
	).Style(
		c.FontFamily("Damion"),
		c.FontSize(128, data.PT),
		c.TextAnchor(data.Middle),
		c.Fill(color.White, 1.0),
	)

	// Output as SVG
	c.Write(data.SVG, os.Stdout)
}
