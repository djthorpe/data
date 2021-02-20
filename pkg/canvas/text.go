package canvas

import (
	"fmt"
	"strings"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/f32"
)

/////////////////////////////////////////////////////////////////////
// TEXT ELEMENTS

func (this *Canvas) Text(pt data.Point, rel bool, children ...data.CanvasText) data.CanvasElement {
	elem, err := this.NewElement("text")
	if err != nil {
		return nil
	}

	// Set attributes
	if rel {
		elem.SetAttr("dx", f32.String(pt.X))
		elem.SetAttr("dy", f32.String(pt.Y))
	} else {
		elem.SetAttr("x", f32.String(pt.X))
		elem.SetAttr("y", f32.String(pt.Y))
	}

	// Add child nodes, which must be elements
	for _, child := range children {
		// Check child
		if child == nil {
			return nil
		} else if child_, ok := child.(*Element); ok == false {
			return nil
		} else if err := elem.AddChild(child_.Node); err != nil {
			return nil
		}
	}

	return elem
}

func (this *Canvas) TextSpan(value string) data.CanvasText {
	elem, err := this.NewElement("tspan")
	if err != nil {
		return nil
	}

	// Add cdata
	value = strings.TrimSpace(value)
	if value != "" {
		if err := elem.AddChild(this.Document.CreateText(value)); err != nil {
			return nil
		}
	}

	// Return success
	return elem
}

func (this *Element) Origin(pt data.Point, rel bool) data.CanvasText {
	// Only possible on tspan elements
	if this.isElement("tspan") == false {
		return nil
	}

	// Set attributes
	if rel {
		this.SetAttr("dx", f32.String(pt.X))
		this.SetAttr("dy", f32.String(pt.Y))
	} else {
		this.SetAttr("x", f32.String(pt.X))
		this.SetAttr("y", f32.String(pt.Y))
	}

	// Success
	return this
}

func (this *Element) Offset(pt data.Point) data.CanvasText {
	return this.Origin(pt, true)
}

func (this *Element) Length(length float32, adjust data.Adjust) data.CanvasText {
	// Only possible on text and tspan elements
	if this.isElement("tspan", "text") == false {
		return nil
	}

	// Set attributes
	if err := this.SetAttr("textLength", f32.String(length)); err != nil {
		return nil
	}
	if textAdjust := fmt.Sprint(adjust); textAdjust != "spacing" {
		if err := this.SetAttr("textAdjust", textAdjust); err != nil {
			return nil
		}
	}

	// Success
	return this
}
