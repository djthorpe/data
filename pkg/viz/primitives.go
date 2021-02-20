package viz

import (
	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Viz) GraphPaper(major, minor uint) data.VizGraphPaper {
	// Create major-xy lines
	if major > 0 {
		// TODO
	}
	// Create minor-xy lines
	if minor > 0 {
		// TODO
	}
	elem := this.NewGroup("graphpaper",
		this.Rect(this.Origin(), this.Size()).Class("border"),
	)
	return elem
}
