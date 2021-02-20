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

func (this *Element) Style(styles ...data.CanvasStyle) data.CanvasElement {
	// Get styles
	if styles := NewStyles(styles); styles == nil {
		return nil
	} else if value := styles.String(); value != "" {
		if err := this.SetAttr("style", value); err != nil {
			return nil
		}
	}

	// Return this element
	return this
}

func (this *Element) Transform(transforms ...data.CanvasTransform) data.CanvasElement {
	// Get transforms into a string array
	attr := make([]string, 0, len(transforms))
	for _, path := range transforms {
		if segment, ok := path.(TransformOperation); ok == false {
			return nil
		} else if segment != "" {
			attr = append(attr, string(segment))
		}
	}

	// Set attribute
	if len(attr) == 0 {
		return this
	} else if err := this.SetAttr("transform", strings.Join(attr, " ")); err != nil {
		return nil
	} else {
		return this
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Element) String() string {
	return fmt.Sprint(this.Node)
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Element) isElement(tags ...string) bool {
	name := this.Node.Name()
	if name.Space != data.XmlNamespaceSVG {
		return false
	}
	for _, tag := range tags {
		if name.Local == tag {
			return true
		}
	}
	return false
}
