package canvas_test

import (
	"strings"
	"testing"

	data "github.com/djthorpe/data"
	canvas "github.com/djthorpe/data/pkg/canvas2"
)

func Test_Canvas_001(t *testing.T) {
	c1 := canvas.NewCanvas(data.Size{16, 16}, data.PX).Version("1.1").Title("Hello, World")

	// Write SVG
	b := new(strings.Builder)
	if err := c1.Write(data.SVG, b); err != nil {
		t.Fatal(err)
	} else {
		t.Log(b.String())
	}
}
