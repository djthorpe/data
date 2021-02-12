package viz_test

import (
	"os"
	"testing"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/table"
	"github.com/djthorpe/data/pkg/viz"
)

const (
	DATASET_A = "../../etc/dataset/12-15-2020.csv"
)

func Test_Series_001(t *testing.T) {
	s := viz.NewSeries()
	if s == nil {
		t.Fatal("Expected non-nil series")
	}
	t.Log(s)
}

func Test_Series_002(t *testing.T) {
	table := table.NewTable()
	fh, err := os.Open(DATASET_A)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()

	// Read float, string and nil values from CSV
	if err := table.Read(fh, table.OptHeader(), table.OptType(data.Nil|data.Float)); err != nil {
		t.Fatal(err)
	}

	// Generate series of data from table
	series := viz.NewSeries()
	if err := series.Read(table, func(i int, row []interface{}) ([]interface{}, error) {
		// Row 0 is state, 3+4 is LatLong, 5 is Cases and 6 is Deaths

		switch i {
		case -1:
			return []interface{}{
				table.Col(0).Name(),
				table.Col(3).Name() + "," + table.Col(4).Name(),
				table.Col(5).Name(),
				table.Col(6).Name(),
			}, nil
		default:
			return []interface{}{
				row[0].(string),
				data.Point{data.Float32(row[3]), data.Float32(row[4])},
				data.Float32(row[5]),
				data.Float32(row[6]),
			}, nil
		}
	}); err != nil {
		t.Error(err)
	}

	for _, set := range series.Sets() {
		t.Log(set)
		if v, ok := set.(data.Values); ok {
			t.Log(v.Scale())
		}
	}

	/*
		// Create an A4 canvas
		canvas := canvas.NewCanvas(data.A4LandscapeSize, data.MM)

		// Write the series to the canvas, scaling values so that (0,0) translates to (0,0) and
		// (W,H) translates to maximum point in series
		points.WritePath(canvas).Scale(data.Point{0, 0}, points.Max())

		// Output SVG
		b := new(strings.Builder)
		if err := canvas.Write(b); err != nil {
			t.Error(err)
		} else {
			t.Log("\n" + b.String())
		}
	*/
}
