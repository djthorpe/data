package canvas

import (
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
	return nil
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
	if err := elem.SetAttr("d", strings.Join(d, " ")); err != nil {
		return nil
	} else {
		return elem
	}
}
