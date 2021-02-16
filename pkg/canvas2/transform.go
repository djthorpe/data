package canvas

/*
/////////////////////////////////////////////////////////////////////
// METHODS

func (elem *Element) Scale(size data.Size) data.CanvasTransform {
	return nil
}

func (elem *Element) Translate(pt data.Point) data.CanvasTransform {
	return &transformdef{Op: translate, Point: pt}
}

func (elem *Element) Rotate(deg float32) data.CanvasTransform {
	return &transformdef{Op: rotate, Angle: deg, Point: data.ZeroPoint}
}

func (elem *Element) RotateAround(deg float32, pt data.Point) data.CanvasTransform {
	return &transformdef{Op: rotate, Angle: deg, Point: pt}
}

func (elem *Element) SkewX(skew float32) data.CanvasTransform {
	return &transformdef{Op: skewx, Angle: skew}
}

func (elem *Element) SkewY(skew float32) data.CanvasTransform {
	return &transformdef{Op: skewy, Angle: skew}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (op *transformdef) String() string {
	switch op.Op {
	case translate:
		if op.X != 0 || op.Y != 0 {
			return fmt.Sprintf("translate(%v)", f32.String(op.X, op.Y))
		}
	case scale:
		if op.W != 1.0 || op.H != 1.0 {
			if op.W == op.H {
				return fmt.Sprintf("scale(%v)", f32.String(op.W))
			} else {
				return fmt.Sprintf("scale(%v)", f32.String(op.W, op.H))
			}
		}
	case rotate:
		if op.Angle != 0 {
			if op.X == 0 && op.Y == 0 {
				return fmt.Sprintf("rotate(%v)", f32.String(op.Angle))
			} else {
				return fmt.Sprintf("rotate(%v)", f32.String(op.Angle, op.X, op.Y))
			}
		}
	case skewx:
		if op.Angle != 0 {
			return fmt.Sprintf("skewx(%v)", f32.String(op.Angle))
		}
	case skewy:
		if op.Angle != 0 {
			return fmt.Sprintf("skewy(%v)", f32.String(op.Angle))
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
*/
