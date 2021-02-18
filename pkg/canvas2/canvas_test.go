package canvas_test

import (
	"fmt"
	"strings"
	"testing"

	data "github.com/djthorpe/data"
	canvas "github.com/djthorpe/data/pkg/canvas2"
)

func CheckError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Canvas_001(t *testing.T) {
	c1 := canvas.NewCanvas(data.Size{16, 16}, data.PX).Version("1.1").Title("Hello, World")

	// Write SVG
	b := new(strings.Builder)
	if err := c1.Write(data.SVG|data.Minify, b); err != nil {
		t.Fatal(err)
	} else if b.String() != "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"16px\" height=\"16px\" viewBox=\"0 0 16 16\" version=\"1.1\"><title>Hello, World</title></svg>" {
		t.Error("Unexpected return, got: ", b.String())
	}

	// Change version and title
	c1.Version("1.2").Title("Goodbye cruel world")

	b = new(strings.Builder)
	if err := c1.Write(data.SVG|data.Minify, b); err != nil {
		t.Fatal(err)
	} else if b.String() != "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"16px\" height=\"16px\" viewBox=\"0 0 16 16\" version=\"1.2\"><title>Goodbye cruel world</title></svg>" {
		t.Error("Unexpected return, got: ", b.String())
	}

	// Remove version and title
	c1.Version("").Title("")

	b = new(strings.Builder)
	if err := c1.Write(data.SVG|data.Minify, b); err != nil {
		t.Fatal(err)
	} else if b.String() != "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"16px\" height=\"16px\" viewBox=\"0 0 16 16\"></svg>" {
		t.Error("Unexpected return, got: ", b.String())
	}

	// Change viewbox
	CheckError(t, c1.SetViewBox(data.ZeroPoint, data.Size{32, 32}))

}

func Test_Canvas_002(t *testing.T) {
	c1 := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	if attr, exists := c1.DOM().Attr("viewBox"); exists {
		if attr.Value != "0 0 16 16" {
			t.Error("Unexpected viewBox:", attr.Value)
		}
	}

	// Change viewbox
	CheckError(t, c1.SetViewBox(data.ZeroPoint, data.Size{32, 32}))
	if attr, exists := c1.DOM().Attr("viewBox"); exists {
		if attr.Value != "0 0 32 32" {
			t.Error("Unexpected viewBox:", attr.Value)
		}
	}

	// Change viewbox
	CheckError(t, c1.SetViewBox(data.Point{-10, -10}, data.Size{32, 32}))
	if attr, exists := c1.DOM().Attr("viewBox"); exists {
		if attr.Value != "-10 -10 32 32" {
			t.Error("Unexpected viewBox:", attr.Value)
		}
	}

}

func Test_Canvas_003(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	if circle := c.Circle(data.ZeroPoint, 10); circle == nil {
		t.Error("Unexpected nil from c.Circle")
	} else if str := fmt.Sprint(circle); str != "<circle cx=\"0\" cy=\"0\" r=\"10\"></circle>" {
		t.Error("Unexpected return, got: ", str)
	}
}

func Test_Canvas_004(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	if shape := c.Ellipse(data.ZeroPoint, data.Size{20, 30}); shape == nil {
		t.Error("Unexpected nil from c.Ellipse")
	} else if str := fmt.Sprint(shape); str != "<ellipse cx=\"0\" cy=\"0\" rx=\"20\" rx=\"30\"></ellipse>" {
		t.Error("Unexpected return, got: ", str)
	}
}

func Test_Canvas_005(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	if shape := c.Rect(data.Point{0, 10}, data.Size{20, 30}); shape == nil {
		t.Error("Unexpected nil from c.Rect")
	} else if str := fmt.Sprint(shape); str != "<rect x=\"0\" y=\"10\" width=\"20\" height=\"30\"></rect>" {
		t.Error("Unexpected return, got: ", str)
	}
}

func Test_Canvas_006(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	if shape := c.Path(data.Point{0, 10}, data.Point{20, 30}); shape == nil {
		t.Error("Unexpected nil from c.Path")
	} else if str := fmt.Sprint(shape); str != "<path d=\"M 0 10 L 20 30\"></path>" {
		t.Error("Unexpected return, got: ", str)
	}
}

func Test_Canvas_007(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	if shape := c.Path(
		c.MoveTo(data.Point{0, 10}),
		c.ClosePath(),
	); shape == nil {
		t.Error("Unexpected nil from c.Path")
	} else if str := fmt.Sprint(shape); str != "<path d=\"M 0 10 Z\"></path>" {
		t.Error("Unexpected return, got: ", str)
	}
}

func Test_Canvas_008(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	if shape := c.Path(
		c.MoveTo(data.Point{0, 10}),
		c.ClosePath(),
	); shape == nil {
		t.Error("Unexpected nil from c.Path")
	} else if str := fmt.Sprint(shape); str != "<path d=\"M 0 10 Z\"></path>" {
		t.Error("Unexpected return, got: ", str)
	}
}

func Test_Canvas_009(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	if shape := c.Path(
		c.MoveTo(data.Point{0, 10}),
		c.LineTo(data.Point{20, 30}),
		c.ClosePath(),
	); shape == nil {
		t.Error("Unexpected nil from c.Path")
	} else if str := fmt.Sprint(shape); str != "<path d=\"M 0 10 L 20 30 Z\"></path>" {
		t.Error("Unexpected return, got: ", str)
	}
}

func Test_Canvas_010(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	if g := c.Group(c.Group(
		c.Circle(data.Point{0, 0}, 10),
	).Id("test")); g == nil {
		t.Error("Unexpected nil from c.Group")
	} else if str := fmt.Sprint(g); str != "<g><g id=\"test\"><circle x=\"0\" y=\"0\" r=\"10\"></circle></g></g>" {
		t.Error("Unexpected return, got: ", str)
	}
}
