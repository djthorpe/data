package viz

import (
	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/geom"
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

// GraphPaper returns an element which contains a grid, with same
// number of squares in X and Y directions. Obviously this won't work
// for graphs, which will need to take X and Y scales as input instead
func (this *Viz) GraphPaper(major, minor uint) data.VizGraphPaper {
	// major box size
	majorSize := geom.DivideSize(this.Size(), float32(major+1))
	//minorSize := geom.DivideSize(this.Size(), float32((major+1)*(minor+1)))
	p1 := this.Origin()
	p2 := geom.AddPoint(this.Origin(), this.Size())

	// Elements
	elements := make([]data.CanvasElement, 0, (major+1)*(minor+1)+1)

	// Create major-xy lines
	if major > 0 {
		// Y-Axis
		for i := uint(0); i < major; i++ {
			elements = append(elements, this.Line(
				data.Point{p1.X, p1.Y + float32(i+1)*majorSize.H},
				data.Point{p2.X, p1.Y + float32(i+1)*majorSize.H},
			).Class(data.ClassGraphPaperYMajor))
		}
		// X-Axis
		for i := uint(0); i < major; i++ {
			elements = append(elements, this.Line(
				data.Point{p1.X + float32(i+1)*majorSize.W, p1.Y},
				data.Point{p1.X + float32(i+1)*majorSize.W, p2.Y},
			).Class(data.ClassGraphPaperXMajor))
		}
	}

	// Add border
	elements = append(elements, this.Rect(p1, this.Size()).Class(data.ClassGraphPaperBorder))

	// Return group
	return this.NewGroup(data.ClassGraphPaper, elements...)
}
