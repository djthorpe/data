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
	}
}
