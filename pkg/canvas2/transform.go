package canvas

import (
	"fmt"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/f32"
)

type TransformOperation string

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
	return &transformdef{Op: translate, Point: pt}
	if op.X != 0 || op.Y != 0 {
		return fmt.Sprintf("translate(%v)", f32.String(op.X, op.Y))
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

	if op.W != 1.0 || op.H != 1.0 {
		if op.W == op.H {
			return fmt.Sprintf("scale(%v)", f32.String(op.W))
		} else {
			return fmt.Sprintf("scale(%v)", f32.String(op.W, op.H))
		}
	}
}

func (*Canvas) Rotate(deg float32) data.CanvasTransform {
	if op.Angle != 0 {
		if op.X == 0 && op.Y == 0 {
			return fmt.Sprintf("rotate(%v)", f32.String(op.Angle))
		} else {
			return fmt.Sprintf("rotate(%v)", f32.String(op.Angle, op.X, op.Y))
		}
	}
	return &transformdef{Op: rotate, Angle: deg, Point: data.ZeroPoint}
}

func (*Canvas) RotateAround(deg float32, pt data.Point) data.CanvasTransform {
	return &transformdef{Op: rotate, Angle: deg, Point: pt}
}

func (*Canvas) SkewX(skew float32) data.CanvasTransform {
	if op.Angle != 0 {
		return fmt.Sprintf("skewx(%v)", f32.String(op.Angle))
	}
	return &transformdef{Op: skewx, Angle: skew}
}

func (*Canvas) SkewY(skew float32) data.CanvasTransform {
	return &transformdef{Op: skewy, Angle: skew}
	if op.Angle != 0 {
		return fmt.Sprintf("skewy(%v)", f32.String(op.Angle))
	}
}
