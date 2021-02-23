package viz_test

import (
	"testing"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/canvas"
	"github.com/djthorpe/data/pkg/set"
	"github.com/djthorpe/data/pkg/viz"
)

func Test_Viz_001(t *testing.T) {
	v := viz.NewViz(canvas.NewCanvas(data.A4LandscapeSize, data.MM))
	if v == nil {
		t.Fatal("Unexpected return from NewViz")
	} else {
		t.Log(v)
	}
}

func Test_Viz_002(t *testing.T) {
	v := viz.NewViz(canvas.NewCanvas(data.A4LandscapeSize, data.MM))
	if v == nil {
		t.Fatal("Unexpected nil return from NewViz")
	}
	set := set.NewRealSet("numbers")
	if set == nil {
		t.Fatal("Unexpected nil return from NewRealSet")
	} else {
		set.Append(1, 2, 4, 6, 8, 9)
	}
	if scale := v.RealScale(set, data.Horizontal); scale == nil {
		t.Fatal("Unexpected return from viz.GraphPaper")
	} else {
		t.Log(v)
	}
}
