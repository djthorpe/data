package viz_test

import (
	"testing"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/canvas"
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
		t.Fatal("Unexpected return from NewViz")
	}
	if g := v.GraphPaper(10, 2); g == nil {
		t.Fatal("Unexpected return from viz.GraphPaper")
	} else {
		t.Log(v)
	}
}
