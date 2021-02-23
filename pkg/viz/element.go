package viz

import (
	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Element struct {
	data.CanvasGroup
	data.Orientation
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Viz) NewGroup(class string, elems ...data.CanvasElement) *Element {
	elem := new(Element)

	// Create group and add elements
	group := this.Canvas.Group(elems...)
	if group == nil {
		return nil
	} else {
		elem.CanvasGroup = group
	}

	// Set classes
	if class != "" {
		if group.Class(class) == nil {
			this.Remove(group)
			return nil
		}
	}

	// Return success
	return elem
}
