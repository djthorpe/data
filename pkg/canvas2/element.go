package canvas

import (
	"fmt"
	"strings"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Element struct {
	data.Node
	*Canvas
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Canvas) NewElement(name string) (*Element, error) {
	elem := new(Element)
	if node := this.Document.CreateElementNS(name, data.XmlNamespaceSVG); node == nil {
		return nil, data.ErrInternalAppError.WithPrefix("NewElement")
	} else if err := this.Document.AddChild(node); err != nil {
		return nil, err
	} else {
		elem.Node = node
		elem.Canvas = this
		return elem, nil
	}
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Element) Id(value string) data.CanvasElement {
	value = strings.TrimSpace(value)
	if value == "" {
		if err := this.Node.RemoveAttr("id"); err != nil {
			return nil
		}
	} else {
		if err := this.Node.SetAttr("id", value); err != nil {
			return nil
		}
	}
	return this
}

func (this *Element) Class(value string) data.CanvasElement {
	value = strings.TrimSpace(value)
	if value == "" {
		if err := this.Node.RemoveAttr("class"); err != nil {
			return nil
		}
	} else {
		if err := this.Node.SetAttr("class", value); err != nil {
			return nil
		}
	}
	return this
}

func (this *Element) Style(...data.CanvasStyle) data.CanvasElement {
	return nil
}

func (this *Element) Transform(...data.CanvasTransform) data.CanvasElement {
	return nil
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Element) String() string {
	return fmt.Sprint(this.Node)
}
