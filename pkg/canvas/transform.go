package canvas

import (
	"fmt"
	"strings"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Transform struct {
	op []data.CanvasTransform
}

type transformop uint

type transformdef struct {
	Op transformop
	data.Point
	data.Size
	Angle float32
}

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	translate transformop = iota
	scale
	rotate
	skewx
	skewy
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewTransform(op []data.CanvasTransform) *Transform {
	this := &Transform{op}
	return this
}

/////////////////////////////////////////////////////////////////////
// METHODS

func (e *Element) Scale(size data.Size) data.CanvasTransform {
	return &transformdef{Op: scale, Size: size}
}

func (e *Element) Translate(pt data.Point) data.CanvasTransform {
	return &transformdef{Op: translate, Point: pt}
}

func (e *Element) Rotate(deg float32) data.CanvasTransform {
	return &transformdef{Op: rotate, Angle: deg, Point: data.ZeroPoint}
}

func (e *Element) RotateAround(deg float32, pt data.Point) data.CanvasTransform {
	return &transformdef{Op: rotate, Angle: deg, Point: pt}
}

func (e *Element) SkewX(skew float32) data.CanvasTransform {
	return &transformdef{Op: skewx, Angle: skew}
}

func (e *Element) SkewY(skew float32) data.CanvasTransform {
	return &transformdef{Op: skewy, Angle: skew}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (op *transformdef) String() string {
	switch op.Op {
	case translate:
		if op.X != 0 || op.Y != 0 {
			return fmt.Sprintf("translate(%v)", string2(op.X, op.Y))
		}
	case scale:
		if op.W != 1.0 || op.H != 1.0 {
			if op.W == op.H {
				return fmt.Sprintf("scale(%v)", string1(op.W))
			} else {
				return fmt.Sprintf("scale(%v)", string2(op.W, op.H))
			}
		}
	case rotate:
		if op.Angle != 0 {
			if op.X == 0 && op.Y == 0 {
				return fmt.Sprintf("rotate(%v)", string1(op.Angle))
			} else {
				return fmt.Sprintf("rotate(%v)", string3(op.Angle, op.X, op.Y))
			}
		}
	case skewx:
		if op.Angle != 0 {
			return fmt.Sprintf("skewx(%v)", string1(op.Angle))
		}
	case skewy:
		if op.Angle != 0 {
			return fmt.Sprintf("skewy(%v)", string1(op.Angle))
		}
	}

	// By default return empty string
	return ""
}

func (t *Transform) String() string {
	str := ""
	for _, op := range t.op {
		if ops := fmt.Sprint(op); ops != "" {
			str += ops + " "
		}
	}
	return strings.TrimSpace(str)
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func string1(x float32) string {
	if float32(int64(x)) == x {
		return fmt.Sprintf("%.0f", x)
	} else {
		return fmt.Sprintf("%f", x)
	}
}

func string2(x, y float32) string {
	return string1(x) + "," + string1(y)
}

func string3(x, y, z float32) string {
	return string1(x) + "," + string1(y) + "," + string1(z)
}
