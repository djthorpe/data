package canvas_test

import (
	"fmt"
	"strings"
	"testing"

	data "github.com/djthorpe/data"
	canvas "github.com/djthorpe/data/pkg/canvas"
	color "github.com/djthorpe/data/pkg/color"
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
	} else if str := fmt.Sprint(shape); str != "<ellipse cx=\"0\" cy=\"0\" rx=\"20\" ry=\"30\"></ellipse>" {
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
	if shape := c.Path(
		c.MoveTo(data.Point{0, 10}),
		c.LineTo(data.Point{20, 30}),
	); shape == nil {
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
	} else if str := fmt.Sprint(g); str != "<g><g id=\"test\"><circle cx=\"0\" cy=\"0\" r=\"10\"></circle></g></g>" {
		t.Error("Unexpected return, got: ", str)
	}
}

func Test_Canvas_011(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	if g := c.Group(); g == nil {
		t.Error("Unexpected nil from c.Group")
	} else if g.Transform() == nil {
		t.Error("Unexpected nil from g.Transform")
	} else if str := fmt.Sprint(g); str != "<g></g>" {
		t.Error("Unexpected return, got: ", str)
	}
}

func Test_Canvas_012(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	if g := c.Group(); g == nil {
		t.Error("Unexpected nil from c.Group")
	} else if g.Transform(
		c.Translate(data.ZeroPoint),
	) == nil {
		t.Error("Unexpected nil from g.Transform")
	} else if str := fmt.Sprint(g); str != "<g></g>" {
		t.Error("Unexpected return, got: ", str)
	}
}

func Test_Canvas_013(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	if g := c.Group(); g == nil {
		t.Error("Unexpected nil from c.Group")
	} else if g.Transform(
		c.Translate(data.Point{10, 10}),
	) == nil {
		t.Error("Unexpected nil from g.Transform")
	} else if str := fmt.Sprint(g); str != "<g transform=\"translate(10,10)\"></g>" {
		t.Error("Unexpected return, got: ", str)
	}
}

func Test_Canvas_014(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	if g := c.Group(); g == nil {
		t.Error("Unexpected nil from c.Group")
	} else if g.Transform(
		c.Translate(data.Point{10, 10}),
		c.Rotate(90),
	) == nil {
		t.Error("Unexpected nil from g.Transform")
	} else if str := fmt.Sprint(g); str != "<g transform=\"translate(10,10) rotate(90)\"></g>" {
		t.Error("Unexpected return, got: ", str)
	}
}

func Test_Canvas_015(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	if g := c.Group(); g == nil {
		t.Error("Unexpected nil from c.Group")
	} else if g.Transform(
		c.Translate(data.Point{10, 10}),
		c.RotateAround(90, data.Point{10, 10}),
	) == nil {
		t.Error("Unexpected nil from g.Transform")
	} else if str := fmt.Sprint(g); str != "<g transform=\"translate(10,10) rotate(90,10,10)\"></g>" {
		t.Error("Unexpected return, got: ", str)
	}
}

func Test_Canvas_016(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	if g := c.Group(); g == nil {
		t.Error("Unexpected nil from c.Group")
	} else if g.Transform(
		c.Scale(data.Size{1, 1}),
	) == nil {
		t.Error("Unexpected nil from g.Transform")
	} else if str := fmt.Sprint(g); str != "<g></g>" {
		t.Error("Unexpected return, got: ", str)
	}
}

func Test_Canvas_017(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	if g := c.Group(); g == nil {
		t.Error("Unexpected nil from c.Group")
	} else if g.Transform(
		c.Scale(data.Size{2, 2}),
		c.Scale(data.Size{1, 3}),
	) == nil {
		t.Error("Unexpected nil from g.Transform")
	} else if str := fmt.Sprint(g); str != "<g transform=\"scale(2) scale(1,3)\"></g>" {
		t.Error("Unexpected return, got: ", str)
	}
}

func Test_Canvas_018(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	if g := c.Group(); g == nil {
		t.Error("Unexpected nil from c.Group")
	} else if g.Transform(
		c.SkewX(2),
		c.SkewY(3),
	) == nil {
		t.Error("Unexpected nil from g.Transform")
	} else if str := fmt.Sprint(g); str != "<g transform=\"skewx(2) skewy(3)\"></g>" {
		t.Error("Unexpected return, got: ", str)
	}
}

func Test_Canvas_019(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	if g := c.Group(); g == nil {
		t.Error("Unexpected nil from c.Group")
	} else if g.Style(
		c.NoFill(),
		c.NoStroke(),
	) == nil {
		t.Error("Unexpected nil from g.Style")
	} else if str := fmt.Sprint(g); str != `<g style="fill:none; stroke:none;"></g>` {
		t.Error("Unexpected return, got: ", str)
	}
}

func Test_Canvas_020(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	if g := c.Group(); g == nil {
		t.Error("Unexpected nil from c.Group")
	} else if g.Style(
		c.Fill(color.Red, 1.0),
	) == nil {
		t.Error("Unexpected nil from g.Style")
	} else if str := fmt.Sprint(g); str != `<g style="fill:red; fill-opacity:1;"></g>` {
		t.Error("Unexpected return, got: ", str)
	}
}

func Test_Canvas_021(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	if g := c.Group(); g == nil {
		t.Error("Unexpected nil from c.Group")
	} else if g.Style(
		c.Stroke(color.Red, 1.0),
	) == nil {
		t.Error("Unexpected nil from g.Style")
	} else if str := fmt.Sprint(g); str != `<g style="stroke:red; stroke-opacity:1;"></g>` {
		t.Error("Unexpected return, got: ", str)
	}
}

func Test_Canvas_022(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	if g := c.Group(); g == nil {
		t.Error("Unexpected nil from c.Group")
	} else if g.Style(
		c.Stroke(color.Red, 1.0),
		c.NoFill(),
	) == nil {
		t.Error("Unexpected nil from g.Style")
	} else if str := fmt.Sprint(g); str != `<g style="fill:none; stroke:red; stroke-opacity:1;"></g>` {
		t.Error("Unexpected return, got: ", str)
	}
}

func Test_Canvas_023(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	if g := c.Group(); g == nil {
		t.Error("Unexpected nil from c.Group")
	} else if g.Style(
		c.Stroke(color.Red, 1.0),
		c.StrokeWidth(2.0),
	) == nil {
		t.Error("Unexpected nil from g.Style")
	} else if str := fmt.Sprint(g); str != `<g style="stroke:red; stroke-opacity:1; stroke-width:2;"></g>` {
		t.Error("Unexpected return, got: ", str)
	}
}

func Test_Canvas_024(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	if g := c.Group(); g == nil {
		t.Error("Unexpected nil from c.Group")
	} else if g.Style(
		c.LineCap(data.CapRound),
		c.LineJoin(data.JoinMiterClip),
		c.MiterLimit(2.0),
	) == nil {
		t.Error("Unexpected nil from g.Style")
	} else if str := fmt.Sprint(g); str != `<g style="stroke-linecap:round; stroke-linejoin:miter-clip; stroke-miterlimit:2;"></g>` {
		t.Error("Unexpected return, got: ", str)
	}
}

func Test_Canvas_025(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	m := c.Marker(data.Point{5, 5}, data.Size{6, 6}, c.Path())
	if m == nil {
		t.Error("Unexpected nil from c.Group")
	} else if str := fmt.Sprint(m); str != `<marker refX="5" refY="5" markerWidth="6" markerHeight="6" orient="auto"><path></path></marker>` {
		t.Error("Unexpected return, got: ", str)
	}
	if m := m.OrientationAngle(0); m == nil {
		t.Error("Unexpected nil from c.Group")
	} else if str := fmt.Sprint(m); str != `<marker refX="5" refY="5" markerWidth="6" markerHeight="6" orient="auto"><path></path></marker>` {
		t.Error("Unexpected return, got: ", str)
	}
	if m := m.OrientationAngle(45); m == nil {
		t.Error("Unexpected nil from c.Group")
	} else if str := fmt.Sprint(m); str != `<marker refX="5" refY="5" markerWidth="6" markerHeight="6" orient="45"><path></path></marker>` {
		t.Error("Unexpected return, got: ", str)
	}
	if m := m.OrientationAngle(0); m == nil {
		t.Error("Unexpected nil from c.Group")
	} else if str := fmt.Sprint(m); str != `<marker refX="5" refY="5" markerWidth="6" markerHeight="6" orient="auto"><path></path></marker>` {
		t.Error("Unexpected return, got: ", str)
	}
	if m := m.OrientationAngle(180); m == nil {
		t.Error("Unexpected nil from c.Group")
	} else if str := fmt.Sprint(m); str != `<marker refX="5" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse"><path></path></marker>` {
		t.Error("Unexpected return, got: ", str)
	}
}

func Test_Canvas_026(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	if l := c.Polyline(); l == nil {
		t.Error("Unexpected nil from c.Polyline")
	} else if str := fmt.Sprint(l); str != `<polyline></polyline>` {
		t.Error("Unexpected return, got: ", str)
	}

	if l := c.Polyline(data.ZeroPoint); l == nil {
		t.Error("Unexpected nil from c.Polyline")
	} else if str := fmt.Sprint(l); str != `<polyline points="0,0"></polyline>` {
		t.Error("Unexpected return, got: ", str)
	}

	if l := c.Polyline(data.ZeroPoint, data.Point{5, 5}); l == nil {
		t.Error("Unexpected nil from c.Polyline")
	} else if str := fmt.Sprint(l); str != `<polyline points="0,0 5,5"></polyline>` {
		t.Error("Unexpected return, got: ", str)
	}
}

func Test_Canvas_027(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	if l := c.Polygon(); l == nil {
		t.Error("Unexpected nil from c.Polygon")
	} else if str := fmt.Sprint(l); str != `<polygon></polygon>` {
		t.Error("Unexpected return, got: ", str)
	}

	if l := c.Polygon(data.ZeroPoint); l == nil {
		t.Error("Unexpected nil from c.Polygon")
	} else if str := fmt.Sprint(l); str != `<polygon points="0,0"></polygon>` {
		t.Error("Unexpected return, got: ", str)
	}

	if l := c.Polygon(data.ZeroPoint, data.Point{5, 5}); l == nil {
		t.Error("Unexpected nil from c.Polygon")
	} else if str := fmt.Sprint(l); str != `<polygon points="0,0 5,5"></polygon>` {
		t.Error("Unexpected return, got: ", str)
	}
}

func Test_Canvas_028(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	if l := c.Polygon().Style(
		c.UseMarker(data.Start, "start"),
		c.UseMarker(data.Middle, "mid"),
		c.UseMarker(data.End, "end"),
	); l == nil {
		t.Error("Unexpected nil from c.Polygon")
	} else if str := fmt.Sprint(l); str != `<polygon style="marker-start:start; marker-mid:mid; marker-end:end;"></polygon>` {
		t.Error("Unexpected return, got: ", str)
	}

	if l := c.Polygon().Style(
		c.UseMarker(data.Start|data.End, "both"),
		c.UseMarker(data.Middle, "mid"),
	); l == nil {
		t.Error("Unexpected nil from c.Polygon")
	} else if str := fmt.Sprint(l); str != `<polygon style="marker-start:both; marker-mid:mid; marker-end:both;"></polygon>` {
		t.Error("Unexpected return, got: ", str)
	}

	if l := c.Polygon().Style(
		c.UseMarker(0, "all"),
	); l == nil {
		t.Error("Unexpected nil from c.Polygon")
	} else if str := fmt.Sprint(l); str != `<polygon style="marker-start:all; marker-mid:all; marker-end:all;"></polygon>` {
		t.Error("Unexpected return, got: ", str)
	}

	if l := c.Polygon().Style(
		c.UseMarker(data.Middle, "one"),
	); l == nil {
		t.Error("Unexpected nil from c.Polygon")
	} else if str := fmt.Sprint(l); str != `<polygon style="marker-mid:one;"></polygon>` {
		t.Error("Unexpected return, got: ", str)
	}
}

func Test_Canvas_029(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	url := "https://mdn.mozillademos.org/files/6457/mdn_logo_only_color.png"
	if l := c.Image(data.ZeroPoint, data.Size{16, 16}, url); l == nil {
		t.Error("Unexpected nil from c.Image")
	} else if str := fmt.Sprint(l); str != `<image x="0" y="0" width="16" height="16" href="https://mdn.mozillademos.org/files/6457/mdn_logo_only_color.png"></image>` {
		t.Error("Unexpected return, got: ", str)
	}
}

func Test_Canvas_030(t *testing.T) {
	c := canvas.NewCanvas(data.Size{16, 16}, data.PX)
	if txt := c.Text(data.ZeroPoint, false); txt == nil {
		t.Error("Unexpected nil from c.Text")
	} else if str := fmt.Sprint(txt); str != `<text x="0" y="0"></text>` {
		t.Error("Unexpected return, got: ", str)
	}

	if txt := c.Text(data.ZeroPoint, true); txt == nil {
		t.Error("Unexpected nil from c.Text")
	} else if str := fmt.Sprint(txt); str != `<text dx="0" dy="0"></text>` {
		t.Error("Unexpected return, got: ", str)
	}

	if txt := c.Text(data.ZeroPoint, false, c.TextSpan("hello")); txt == nil {
		t.Error("Unexpected nil from c.Text")
	} else if str := fmt.Sprint(txt); str != `<text x="0" y="0"><tspan>hello</tspan></text>` {
		t.Error("Unexpected return, got: ", str)
	}

	if txt := c.Text(data.ZeroPoint, false,
		c.TextSpan("hello"),
		c.TextSpan("world").Offset(data.Point{0, 5}),
	); txt == nil {
		t.Error("Unexpected nil from c.Text")
	} else if str := fmt.Sprint(txt); str != `<text x="0" y="0"><tspan>hello</tspan><tspan dx="0" dy="5">world</tspan></text>` {
		t.Error("Unexpected return, got: ", str)
	}

	if txt := c.Text(data.ZeroPoint, false,
		c.TextSpan("hello").Length(100, 0),
		c.TextSpan("world").Length(50, data.SpacingAndGlyphs),
	); txt == nil {
		t.Error("Unexpected nil from c.Text")
	} else if str := fmt.Sprint(txt); str != `<text x="0" y="0"><tspan textLength="100">hello</tspan><tspan textLength="50" textAdjust="spacingAndGlyphs">world</tspan></text>` {
		t.Error("Unexpected return, got: ", str)
	} else {
		t.Log(txt)
	}
}
