package viz

import (
	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

// RealScale returns an element which draws either an X-Axis scale
// or Y-Axis scale for real values
func (this *Viz) RealScale(set data.RealSet, orientation data.Orientation) data.VizScale {
	if element := this.NewGroup(data.ClassScale); element == nil {
		return nil
	} else {
		element.Orientation = orientation

		// Return group
		return element
	}
}
