package canvas

import (
	"strings"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/f32"
)

/////////////////////////////////////////////////////////////////////
// GROUP ELEMENTS

func (this *Canvas) Group(children ...data.CanvasElement) data.CanvasGroup {
	g, err := this.NewElement("g")
	if err != nil {
		return nil
	}
	// Append children. If any children are nil, then return nil to bubble up
	// any errors
	for _, child := range children {
		if child == nil {
			return nil
		} else if elem, ok := child.(*Element); ok == false {
			return nil
		} else if err := g.AddChild(elem.Node); err != nil {
			return nil
		}
	}
	return g
}

func (this *Canvas) Defs(children ...data.CanvasElement) data.CanvasGroup {
	// Create a defs element
	g, err := this.NewElement("defs")
	if err != nil {
		return nil
	}

	// Append children
	for _, child := range children {
		if child == nil {
			return nil
		} else if elem, ok := child.(*Element); ok == false {
			return nil
		} else if err := g.AddChild(elem.Node); err != nil {
			return nil
		}
	}

	// Move defs to the first element of the canvas
	this.InsertChildBefore(g, this.FirstChild())

	// Return canvas group
	return g
}

func (this *Canvas) Marker(pt data.Point, sz data.Size, children ...data.CanvasElement) data.CanvasGroup {
	m, err := this.NewElement("marker")
	if err != nil {
		return nil
	}

	// Set attributes on element
	if pt != data.ZeroPoint {
		m.SetAttr("refX", f32.String(pt.X))
		m.SetAttr("refY", f32.String(pt.Y))
	}
	if sz != data.ZeroSize {
		m.SetAttr("markerWidth", f32.String(sz.W))
		m.SetAttr("markerHeight", f32.String(sz.H))
	}
	// Set auto
	if m.OrientationAngle(0) == nil {
		return nil
	}

	// Append children. If any children are nil, then return nil to bubble up
	// any errors
	for _, child := range children {
		if child == nil {
			return nil
		} else if elem, ok := child.(*Element); ok == false {
			return nil
		} else if err := m.AddChild(elem.Node); err != nil {
			return nil
		}
	}
	return m
}

func (this *Element) Desc(cdata string) data.CanvasGroup {
	cdata = strings.TrimSpace(cdata)

	// Remove existing desc tags
	if desc := this.Document.GetElementsByTagNameNS("desc", data.XmlNamespaceSVG); len(desc) != 0 {
		for _, child := range desc {
			this.Document.RemoveChild(child)
		}
	}

	// Create a new desc tag - put at top of element
	if cdata != "" {
		desc := this.Document.CreateElementNS("desc", data.XmlNamespaceSVG)
		if err := desc.AddChild(this.Document.CreateText(cdata)); err != nil {
			return nil
		} else if err := this.Document.InsertChildBefore(desc, this.Document.FirstChild()); err != nil {
			return nil
		}
		// If there is a title tag, then put the desc tag after the title tag
		if title := this.Document.GetElementsByTagNameNS("title", data.XmlNamespaceSVG); len(title) != 0 {
			this.Document.InsertChildBefore(desc, title[0].NextSibling())
		}
	}

	// Return success
	return this
}

func (this *Element) OrientationAngle(angle float32) data.CanvasGroup {
	if this.isElement("marker") == false {
		return nil
	}
	switch angle {
	case 0:
		if err := this.SetAttr("orient", "auto"); err != nil {
			return nil
		}
	case 180:
		if err := this.SetAttr("orient", "auto-start-reverse"); err != nil {
			return nil
		}
	default:
		if err := this.SetAttr("orient", f32.String(angle)); err != nil {
			return nil
		}
	}
	// Return group
	return this
}
