package canvas_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/canvas"
	"github.com/djthorpe/data/pkg/color"
	"github.com/djthorpe/data/pkg/f32"
	"github.com/djthorpe/data/pkg/geom"
)

func Test_Canvas_001(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX).Version("1.1")

	// Output as SVG
	b := new(strings.Builder)
	if err := c.Write(b); err != nil {
		t.Error(err)
	} else {
		t.Log(b.String())
	}
}

func Test_Canvas_002(t *testing.T) {
	c := canvas.NewCanvas(data.A4LandscapeSize, data.MM)
	c.Title("Two circles")
	c.Desc("Two circles, one of which is scaled")

	// Gray background
	c.Rect(c.Origin(), c.Size()).Style(
		c.Fill(color.LightGray, 1.0),
	)

	// Centre point
	centre := geom.CentrePoint(c.Origin(), c.Size())
	size := geom.MultiplySize(c.Size(), 0.5)

	// First circle
	c.Group(
		c.Circle(centre, f32.Min(size.W, size.H)).Id("circle1"),
	).Style(
		c.NoFill(),
		c.Stroke(color.Black, 1.0),
	).Id("g1")

	// Second circle 0.5 times first
	c.Group(
		c.Circle(data.ZeroPoint, f32.Min(size.W, size.H)).Transform(
			c.Scale(data.Size{0.5, 0.5}),
			c.Translate(centre),
		).Id("circle2"),
	).Style(
		c.Fill(color.Red, 1.0),
		c.NoStroke(),
	).Id("g2")

	// Output as SVG
	b := new(strings.Builder)
	if err := c.Write(b); err != nil {
		t.Error(err)
	} else {
		t.Log(b.String())
	}
}

func Test_Canvas_003(t *testing.T) {
	c := canvas.NewCanvas(data.A4LandscapeSize, data.MM)

	// Top of document
	c.Title("Text")
	c.Desc("Text on centre of page")

	// Gray background
	c.Rect(c.Origin(), c.Size()).Style(
		c.Fill(color.LightGray, 1.0),
	)

	// Centre point
	centre := geom.CentrePoint(c.Origin(), c.Size())

	// Text
	c.Text(centre,
		c.Span("Hello, world!"),
	).Style(
		c.FontSize(12, data.PT),
		c.TextAnchor(data.Middle),
	)

	// Output as SVG
	b := new(strings.Builder)
	if err := c.Write(b); err != nil {
		t.Error(err)
	} else {
		t.Log(b.String())
	}
}

func Test_Canvas_004(t *testing.T) {
	c := canvas.NewCanvas(data.A4LandscapeSize, data.MM)

	// Top of document
	c.Title("Color Palette")
	c.Desc("Display all colors in the palette")

	// Gray background
	c.Rect(c.Origin(), c.Size()).Style(
		c.Fill(color.LightGray, 1.0),
	)

	// Determine how many items across and down assuming
	// a square
	palette := color.Palette()
	count := float32(len(palette))
	across := f32.Floor(f32.Sqrt(count))
	size := geom.DivideSize(c.Size(), across)
	for i, cr := range palette {
		origin := data.Point{
			X: float32(i%int(across)) * size.W,
			Y: float32(i/int(across)) * size.H,
		}
		c.Rect(origin, size).Style(
			c.NoStroke(),
			c.Fill(cr, 1.0),
		)
		c.Rect(origin, size).Style(
			c.NoFill(),
			c.Stroke(color.LightGray, 1.0),
			c.StrokeWidth(0.5),
		)
		c.Text(
			geom.CentrePoint(origin, size),
			c.Span(color.String(cr)),
			c.Span(fmt.Sprint("D:", color.Distance(cr, color.Black))),
		).Style(
			c.TextAnchor(data.Middle),
			c.FontSize(9, data.PT),
		)
	}

	// Output as SVG
	b := new(strings.Builder)
	if err := c.Write(b); err != nil {
		t.Error(err)
	} else {
		t.Log(b.String())
	}
}
