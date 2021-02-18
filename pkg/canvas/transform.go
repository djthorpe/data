package canvas

import (
	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/f32"
)

type TransformOperation string

var (
	NilTransformOperation = TransformOperation("")
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewTransformOperation(op string, args ...float32) TransformOperation {
	if len(args) > 0 {
		return TransformOperation(op + "(" + f32.Join(args, ",") + ")")
	} else {
		return TransformOperation(op + "()")
	}
}

/////////////////////////////////////////////////////////////////////
// METHODS

func (*Canvas) Translate(pt data.Point) data.CanvasTransform {
	if pt.X != 0 || pt.Y != 0 {
		return NewTransformOperation("translate", pt.X, pt.Y)
	} else {
		return NilTransformOperation
	}
}

func (*Canvas) Scale(size data.Size) data.CanvasTransform {
	if size.W == 1.0 && size.H == 1.0 {
		return TransformOperation("")
	} else if size.W == size.H {
		return NewTransformOperation("scale", size.W)
	} else {
		return NewTransformOperation("scale", size.W, size.H)
	}
}

func (*Canvas) Rotate(deg float32) data.CanvasTransform {
	if deg != 0 {
		return NewTransformOperation("rotate", deg)
	} else {
		return NilTransformOperation
	}
}

func (*Canvas) RotateAround(deg float32, pt data.Point) data.CanvasTransform {
	if deg == 0 {
		return NilTransformOperation
	} else if pt == data.ZeroPoint {
		return NewTransformOperation("rotate", deg)
	} else {
		return NewTransformOperation("rotate", deg, pt.X, pt.Y)
	}
}

func (*Canvas) SkewX(deg float32) data.CanvasTransform {
	if deg == 0 {
		return NilTransformOperation
	} else {
		return NewTransformOperation("skewx", deg)
	}
}

func (*Canvas) SkewY(deg float32) data.CanvasTransform {
	if deg == 0 {
		return NilTransformOperation
	} else {
		return NewTransformOperation("skewy", deg)
	}
}
