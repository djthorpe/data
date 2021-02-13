package canvas_test

import (
	"strings"
	"testing"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/canvas"
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
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	c.Title("Two circles")
	c.Desc("Two circles, one of which is scaled")
	c.Group(
		c.Circle(data.Point{8, 8}, 8).Id("circle1"),
	).Id("g1")
	c.Group(
		c.Circle(data.Point{8, 8}, 8).Transform(
			c.Scale(data.Size{2.5, 2.5}),
			c.Rotate(45),
		).Id("circle2"),
	).Id("g2")

	// Output as SVG
	b := new(strings.Builder)
	if err := c.Write(b); err != nil {
		t.Error(err)
	} else {
		t.Log(b.String())
	}
}
