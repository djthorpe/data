package canvas

import (
	"net/url"
	"strings"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/f32"
)

/////////////////////////////////////////////////////////////////////
// DRAWING PRIMITIVES

func (this *Canvas) Circle(centre data.Point, radius float32) data.CanvasElement {
	elem, err := this.NewElement("circle")
	if err != nil {
		return nil
	}

	elem.SetAttr("cx", f32.String(centre.X))
	elem.SetAttr("cy", f32.String(centre.Y))
	elem.SetAttr("r", f32.String(f32.Abs(radius)))
	return elem
}

func (this *Canvas) Ellipse(centre data.Point, radius data.Size) data.CanvasElement {
	elem, err := this.NewElement("ellipse")
	if err != nil {
		return nil
	}

	elem.SetAttr("cx", f32.String(centre.X))
	elem.SetAttr("cy", f32.String(centre.Y))
	elem.SetAttr("rx", f32.String(radius.W))
	elem.SetAttr("ry", f32.String(radius.H))
	return elem
}

func (this *Canvas) Line(p1, p2 data.Point) data.CanvasElement {
	elem, err := this.NewElement("line")
	if err != nil {
		return nil
	}

	elem.SetAttr("x1", f32.String(p1.X))
	elem.SetAttr("y1", f32.String(p1.Y))
	elem.SetAttr("x2", f32.String(p2.X))
	elem.SetAttr("y2", f32.String(p2.Y))
	return elem
}

func (this *Canvas) Polyline(pts ...data.Point) data.CanvasElement {
	elem, err := this.NewElement("polyline")
	if err != nil {
		return nil
	}

	points := make([]string, 0, len(pts))
	for _, pt := range pts {
		points = append(points, f32.String(pt.X, pt.Y))
	}

	if attr := strings.Join(points, " "); attr != "" {
		if err := elem.SetAttr("points", attr); err != nil {
			return nil
		}
	}

	return elem
}

func (this *Canvas) Polygon(pts ...data.Point) data.CanvasElement {
	elem, err := this.NewElement("polygon")
	if err != nil {
		return nil
	}

	points := make([]string, 0, len(pts))
	for _, pt := range pts {
		points = append(points, f32.String(pt.X, pt.Y))
	}

	if attr := strings.Join(points, " "); attr != "" {
		if err := elem.SetAttr("points", attr); err != nil {
			return nil
		}
	}

	return elem
}

func (this *Canvas) Rect(pt data.Point, sz data.Size) data.CanvasElement {
	elem, err := this.NewElement("rect")
	if err != nil {
		return nil
	}

	elem.SetAttr("x", f32.String(pt.X))
	elem.SetAttr("y", f32.String(pt.Y))
	elem.SetAttr("width", f32.String(sz.W))
	elem.SetAttr("height", f32.String(sz.H))
	return elem
}

func (this *Canvas) Text(data.Point, ...data.CanvasText) data.CanvasElement {
	// TODO
	return nil
}

func (this *Canvas) Image(pt data.Point, sz data.Size, u string) data.CanvasElement {
	if url, err := url.Parse(u); err != nil {
		return nil
	} else if elem, err := this.NewElement("image"); err != nil {
		return nil
	} else {
		elem.SetAttr("x", f32.String(pt.X))
		elem.SetAttr("y", f32.String(pt.Y))
		elem.SetAttr("width", f32.String(sz.W))
		elem.SetAttr("height", f32.String(sz.H))
		elem.SetAttr("src", url.String())
		return elem
	}
}

func (this *Canvas) Path(paths ...data.CanvasPath) data.CanvasElement {
	// Get path elements into a string array
	d := make([]string, 0, len(paths))
	for _, path := range paths {
		if segment, ok := path.(PathSegment); ok == false {
			return nil
		} else if segment != "" {
			d = append(d, string(segment))
		}
	}

	// Create path element
	elem, err := this.NewElement("path")
	if err != nil {
		return nil
	}

	// Set attribute
	if attr := strings.Join(d, " "); attr != "" {
		if err := elem.SetAttr("d", attr); err != nil {
			return nil
		}
	}

	// Return relement
	return elem
}
