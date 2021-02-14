// Geometry functions
package geom

import (
	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/f32"
)

/////////////////////////////////////////////////////////////////////
// METHODS

// Return the centre point in a rectangle
func CentrePoint(pt data.Point, sz data.Size) data.Point {
	return data.Point{
		X: pt.X + f32.Abs(sz.W)/2.0,
		Y: pt.Y + f32.Abs(sz.H)/2.0,
	}
}

// Multiply size by a constant
func MultiplySize(sz data.Size, k float32) data.Size {
	if k == 0 {
		return data.ZeroSize
	}
	return data.Size{
		W: f32.Abs(sz.W) * k,
		H: f32.Abs(sz.H) * k,
	}
}

// Divide size by a constant
func DivideSize(sz data.Size, k float32) data.Size {
	if k == 0 {
		return data.NilSize
	}
	return data.Size{
		W: f32.Abs(sz.W) / k,
		H: f32.Abs(sz.H) / k,
	}
}