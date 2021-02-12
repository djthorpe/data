package viz_test

import (
	"os"
	"testing"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/canvas"
	"github.com/djthorpe/data/pkg/table"
	"github.com/djthorpe/data/pkg/viz"
)

const (
	DATASET_A = "../../etc/dataset/12-15-2020.csv"
)

func Test_Points_001(t *testing.T) {
	points := viz.NewPoints("a")
	if points == nil {
		t.Fatal("Expected non-nil points")
	}
	t.Log(points)
}

func Test_Points_002(t *testing.T) {
	table := table.NewTable()
	fh, err := os.Open(DATASET_A)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()
	if err := table.Read(fh); err != nil {
		t.Fatal(err)
	}

	// "Confirmed Deaths" is column 5
	points := viz.NewPoints(table.Col(5).Name())
	canvas := canvas.NewCanvas(data.A4LandscapeSize, data.MM)
	points.Read(table, func(i int, row []interface{}) (data.Point, error) {
		if y, ok := row[5].(uint64); ok == false {
			return data.ZeroPoint, data.ErrSkipTransform
		} else {
			return data.Point{float32(i), float32(y)}, nil
		}
	})
	// Write the series to the canvas
	points.WritePath(canvas)
}
