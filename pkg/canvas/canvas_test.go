package canvas_test

import (
	"os"
	"testing"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/canvas"
)

func Test_Canvas_001(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX).Version("1.1")
	if err := c.Write(os.Stdout); err != nil {
		t.Error(err)
	}
}

func Test_Canvas_002(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	c.Title("Style inheritance and the use element")
	c.Desc("Two circles, one of which is a re-styled clone of the other.")
	c.Group(
		c.Circle(data.Point{8, 8}, 8),
	)
	c.Group(
		c.Circle(data.Point{8, 8}, 8),
	)
	if err := c.Write(os.Stdout); err != nil {
		t.Error(err)
	}
}
