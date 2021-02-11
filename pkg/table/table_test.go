package table_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	data "github.com/djthorpe/data"
	table "github.com/djthorpe/data/pkg/table"
)

const (
	TABLE_A = "0,1\n1,2\n1,\n,2\n"
	TABLE_B = "1,0\n1,2\n1,\n,2\n"
	TABLE_C = "2020-10-01\n1/1/2008\n30 Mar 2006\n6 Oct 06\nAug 14 1908"
)

const (
	DATASET_A = "../../etc/dataset/12-15-2020.csv"
	DATASET_B = "../../etc/dataset/time_series_covid19_confirmed_global.csv"
)

func Test_Table_001(t *testing.T) {
	if c := table.NewTable(data.ZeroSize); c == nil {
		t.Error("Expected Non-nil return from NewTable")
	}
	if c := table.NewTable(data.Size{-1, 0}); c != nil {
		t.Error("Expected Nil return from NewTable")
	}
	if c := table.NewTable(data.Size{0, -1}); c != nil {
		t.Error("Expected Nil return from NewTable")
	}
	if c := table.NewTable(data.Size{1.1, 1}); c != nil {
		t.Error("Expected Nil return from NewTable")
	}
	if c := table.NewTable(data.Size{1.999999, 1}); c != nil {
		t.Error("Expected Nil return from NewTable")
	}
	if c := table.NewTable(data.Size{0, 1.99999}); c != nil {
		t.Error("Expected Nil return from NewTable")
	}
	if c := table.NewTable(data.Size{0, 2}); c == nil {
		t.Error("Expected Non-nil return from NewTable")
	}
	if c := table.NewTable(data.Size{2, 2}); c == nil {
		t.Error("Expected Non-nil return from NewTable")
	}
}

func Test_Table_002(t *testing.T) {
	c := table.NewTable(data.ZeroSize)
	r := strings.NewReader(TABLE_A)
	if err := c.Read(r); err != nil {
		t.Error(err)
	} else {
		t.Log(c)
	}
}

func Test_Table_003(t *testing.T) {
	// Read in two tables, re-ordering rows in second read
	c := table.NewTable(data.ZeroSize)
	if err := c.Read(strings.NewReader(TABLE_A), c.OptHeader()); err != nil {
		t.Error(err)
	} else if err := c.Read(strings.NewReader(TABLE_B), c.OptHeader()); err != nil {
		t.Error(err)
	} else {
		t.Log(c)
	}
}

func Test_Table_004(t *testing.T) {
	// Read table with existing columns and rows
	c := table.NewTable(data.Size{2, 2})
	if err := c.Read(strings.NewReader(TABLE_A), c.OptHeader(), c.OptType(data.DefaultTypes|data.Nil)); err != nil {
		t.Error(err)
	} else if err := c.Write(os.Stdout, c.OptHeader()); err != nil {
		t.Error(err)
	} else {
		t.Log(c)
	}
}

func Test_Table_005(t *testing.T) {
	transformFloat := func(t data.Table) data.TransformFunc {
		return func(i, j int, v interface{}) (interface{}, error) {
			col := t.Col(j)
			if v, ok := v.(float64); ok {
				return fmt.Sprintf("%v %.2f", col, v), nil
			} else {
				return nil, data.ErrSkipTransform
			}
		}
	}
	// Read table with existing columns and rows
	c := table.NewTable(data.ZeroSize)
	fh, err := os.Open(DATASET_A)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()
	if err := c.Read(fh, c.OptHeader(), c.OptType(data.DefaultTypes|data.Nil)); err != nil {
		t.Error(err)
	} else if err := c.Write(os.Stdout, c.OptHeader(), c.OptAscii(0, data.BorderLines), c.OptTransform(transformFloat(c))); err != nil {
		t.Error(err)
	}
}

func Test_Table_006(t *testing.T) {
	// Read table with existing columns and rows
	c := table.NewTable(data.ZeroSize)
	fh, err := os.Open(DATASET_B)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()
	if err := c.Read(fh, c.OptHeader(), c.OptType(data.DefaultTypes|data.Nil)); err != nil {
		t.Error(err)
	} else if err := c.Write(os.Stdout, c.OptAscii(0, data.BorderLines)); err != nil {
		t.Error(err)
	}
}

func Test_Table_007(t *testing.T) {
	rowIterator := func(c data.Table) data.IteratorFunc {
		return func(i int, r []interface{}) error {
			t.Log(i, "=>", r)
			return nil
		}
	}
	// Read table with existing columns and rows
	c := table.NewTable(data.ZeroSize)
	if err := c.Read(strings.NewReader(TABLE_A), c.OptType(data.Duration|data.Nil), c.OptRowIterator(rowIterator(c))); err != nil {
		t.Error(err)
	} else if err := c.Write(os.Stdout, c.OptAscii(0, "")); err != nil {
		t.Error(err)
	}
}

func Test_Table_008(t *testing.T) {
	rowIterator := func(c data.Table) data.IteratorFunc {
		return func(i int, r []interface{}) error {
			t.Log(i, "=>", r)
			return nil
		}
	}

	compareFunc := func(c data.Table) data.CompareFunc {
		return func(a, b []interface{}) bool {
			a_ := a[0].(time.Time)
			b_ := b[0].(time.Time)
			return b_.After(a_)
		}
	}

	// Read table with existing columns and rows
	c := table.NewTable(data.ZeroSize)
	if err := c.Read(strings.NewReader(TABLE_C), c.OptType(data.DefaultTypes), c.OptRowIterator(rowIterator(c))); err != nil {
		t.Error(err)
	}
	c.Sort(compareFunc(c))
	// Write sorted data
	if err := c.Write(os.Stdout, c.OptHeader(), c.OptAscii(0, "")); err != nil {
		t.Error(err)
	}
}

func Test_Table_009(t *testing.T) {
	// Read table with existing columns and rows
	c := table.NewTable(data.ZeroSize)
	c.Append(0, 0, 0)
	c.Append(1)
	c.Append(2, 2, 2, 2, 2)
	// Write out
	if err := c.Write(os.Stdout, c.OptAscii(0, "")); err != nil {
		t.Error(err)
	}
}
