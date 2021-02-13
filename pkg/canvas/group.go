package canvas

import "github.com/djthorpe/data"

/////////////////////////////////////////////////////////////////////
// METHODS

func (e *Element) Group(children ...data.CanvasElement) data.CanvasGroup {
	g := NewElement("g", "", e.root)
	e.addChild(g)
	for _, node := range children {
		e.removeChild(node.(*Element))
		g.addChild(node.(*Element))
	}
	return g
}

func (e *Element) Transform(op ...data.CanvasTransform) data.CanvasElement {
	if e.transform == nil {
		e.transform = NewTransform(op)
	} else {
		e.transform.op = append(e.transform.op, op...)
	}
	return e
}

func (e *Element) Style(styles ...data.CanvasStyle) data.CanvasElement {
	if e.style == nil {
		e.style = NewStyle(styles)
	} else {
		e.style.Append(styles)
	}
	return e
}
