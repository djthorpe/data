package set_test

import (
	"os"
	"testing"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/f32"
	"github.com/djthorpe/data/pkg/set"
	"github.com/djthorpe/data/pkg/table"
)

const (
	DATASET_A = "../../etc/dataset/12-15-2020.csv"
)

func Test_Series_001(t *testing.T) {
	s := set.NewSeries()
	if s == nil {
		t.Fatal("Expected non-nil series")
	}
	t.Log(s)
}

func Test_Series_002(t *testing.T) {
	// Create table and series
	series := set.NewSeries()
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
				data.Point{f32.Cast(row[3]), f32.Cast(row[4])},
				row[5].(float64),
				row[6].(float64),
			}, nil
		}
	}); err != nil {
		t.Error(err)
	}

	// Output state, longlat, cases and deaths
	t.Log(series)
}
