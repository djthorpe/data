package sets_test

import (
	"os"
	"testing"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/sets"
	"github.com/djthorpe/data/pkg/table"
)

const (
	DATASET_A = "../../etc/dataset/12-15-2020.csv"
)

func Test_Series_001(t *testing.T) {
	s := sets.NewSeries()
	if s == nil {
		t.Fatal("Expected non-nil series")
	}
	t.Log(s)
}

func Test_Series_002(t *testing.T) {
	// Create table and series
	series := sets.NewSeries()
	table := table.NewTable()
	if series == nil || table == nil {
		t.Fatal("Expected non-nil series and table")
	}

	// Read dataset
	fh, err := os.Open(DATASET_A)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()

	// Read longitude and latitude from CSV
	if err := table.Read(fh, table.OptHeader(), table.OptType(data.Nil|data.Float)); err != nil {
		t.Fatal(err)
	} else if err := series.Read(table, func(i int, row []interface{}) ([]interface{}, error) {
		// Row 0 is state, 3+4 is LatLong, 5 is Cases and 6 is Deaths
		if i == -1 {
			return []interface{}{"state", "longlat", "cases", "deaths"}, nil
		} else {
			return []interface{}{
				row[0].(string),
				data.Point{data.Float32(row[3]), data.Float32(row[4])},
				row[5].(float64),
				row[6].(float64),
			}, nil
		}
	}); err != nil {
		t.Error(err)
	}

	// Output longitude and latitude values
	t.Log(series)

}

/*
func Test_Series_002(t *testing.T) {


	// Generate series of data from table
	series := viz.NewSeries()


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
}
*/
