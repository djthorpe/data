package canvas_test

import (
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

	// Second circle 0.5 times first and rotate by 45 degrees
	c.Group(
		c.Circle(centre, f32.Min(size.W, size.H)).Transform(
			c.Scale(data.Size{0.5, 0.5}),
			c.Rotate(45),
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
