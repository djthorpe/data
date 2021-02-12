package table_test

import (
	"fmt"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	data "github.com/djthorpe/data"
	table "github.com/djthorpe/data/pkg/table"
)

const (
	TABLE_A = "0,1,2,3\n1,2,3,\n2,3,,\n,,3,4\n"
	TABLE_B = "1,0\n1,2\n1,\n,2\n"
	TABLE_C = "2020-10-01,1h\n1/1/2008,2h\n30 Mar 2006,3h30m\n6 Oct 06,40ns\nAug 14 1908,5.5"
	TABLE_D = "8.8.8.8\n127.0.0.1\n192.168.0.1\n\"\"\n"
)

const (
	DATASET_A = "../../etc/dataset/12-15-2020.csv"
	DATASET_B = "../../etc/dataset/time_series_covid19_confirmed_global.csv"
)

func Test_Table_001(t *testing.T) {
	// Create an empty table
	c := table.NewTable()
	if c == nil {
		t.Fatal("Expected Non-nil return from NewTable")
	} else if c.Len() != 0 {
		t.Error("Expected empty table")
	}

	// Append an empty row, then one with four values
	c.Append()
	c.Append(1, 2, 3, 4)

	// Check table
	if c.Len() != 2 {
		t.Error("Expected table length to be two")
	}
	for i := 0; i < 4; i++ {
		if col := c.Col(i); col == nil {
			t.Error("Unexpected non-nil column")
		} else if col.Type()&data.Int == 0 {
			t.Error("Unexpected int type for column, got <", col.Type(), ">")
		} else if col.Type()&data.Nil == 0 {
			t.Error("Unexpected nil type for column, got <", col.Type(), ">")
		} else {
			t.Log("col", i, "=", col)
		}
	}
	for i := 0; i < c.Len(); i++ {
		row := c.Row(i)
		if len(row) != 4 {
			t.Error("Unexpected row width")
		}
	}

	// Output table
	b := new(strings.Builder)
	if err := c.Write(b, c.OptHeader(), c.OptAscii(0, "")); err != nil {
		t.Error(err)
	} else {
		t.Log(b.String())
	}
}

func Test_Table_002(t *testing.T) {
	// The following creates a table of four columns and adds four rows from TABLE_A
	c := table.NewTable("", "a", "b", "a")
	if c == nil {
		t.Fatal("Expected Non-nil return from NewTable")
	}
	if err := c.Read(strings.NewReader(TABLE_A)); err != nil {
		t.Fatal(err)
	}
	// There should be four columns and four rows, fourth column name
	// should be "a" but the key is not
	if c.Len() != 4 {
		t.Error("Expected table length to be 4")
	}
	for i := 0; i < 4; i++ {
		col := c.Col(i)
		if col == nil {
			t.Error("Unexpected non-nil column")
		} else if col.Type()&data.Uint == 0 {
			t.Error("Unexpected uint type for column, got <", col.Type(), ">")
		} else if col.Type()&data.Nil == 0 {
			t.Error("Unexpected nil type for column, got <", col.Type(), ">")
		}
		if i == 3 {
			if col.Name() != "a" {
				t.Error("Unexpected name for column 4, got <", col.Name(), ">")
			}
		}
	}

	// Output table
	b := new(strings.Builder)
	if err := c.Write(b, c.OptHeader(), c.OptAscii(0, "")); err != nil {
		t.Error(err)
	} else {
		t.Log(b.String())
	}
}

func Test_Table_003(t *testing.T) {
	// In this test we do same ingestion but don't allow nil values
	// in which case each column has types uint|string
	c := table.NewTable()
	if c == nil {
		t.Fatal("Expected Non-nil return from NewTable")
	}
	if err := c.Read(strings.NewReader(TABLE_A), c.OptType(data.DefaultTypes)); err != nil {
		t.Fatal(err)
	}
	// There should be four columns and four rows
	if c.Len() != 4 {
		t.Error("Expected table length to be 4")
	}
	for i := 0; i < 4; i++ {
		col := c.Col(i)
		if col == nil {
			t.Error("Unexpected non-nil column")
		} else if col.Type()&data.Uint == 0 {
			t.Error("Unexpected uint type for column, got <", col.Type(), ">")
		} else if col.Type()&data.Nil != 0 {
			t.Error("Unexpected nil type for column, got <", col.Type(), ">")
		} else if col.Type()&data.String == 0 {
			t.Error("Unexpected non-string type for column, got <", col.Type(), ">")
		}
	}

	// Output table
	b := new(strings.Builder)
	if err := c.Write(b, c.OptHeader(), c.OptAscii(0, "")); err != nil {
		t.Error(err)
	} else {
		t.Log(b.String())
	}
}

func Test_Table_004(t *testing.T) {
	// Read in two tables, re-ordering rows in second read
	c := table.NewTable()
	if err := c.Read(strings.NewReader(TABLE_A), c.OptHeader()); err != nil {
		t.Fatal(err)
	}
	// On first input, header rows are called "0","1","2","3"
	for i := 0; i < 4; i++ {
		col := c.Col(i)
		if col.Name() != fmt.Sprint(i) {
			t.Error("Unexpected colname, got <", col.Type(), ">")
		} else if col.Type().Is(data.Uint|data.Nil) == false {
			t.Error("Unexpected type uint|nil for column, got <", col.Type(), ">")
		} else if col.Type().Is(data.String) {
			t.Error("Unexpected string type for column, got <", col.Type(), ">")
		}
	}

	// Merge in second table, but don't allow nils
	if err := c.Read(strings.NewReader(TABLE_B), c.OptHeader(), c.OptType(data.DefaultTypes)); err != nil {
		t.Fatal(err)
	}

	// Headings are still 0,1,2,3
	for i := 0; i < 4; i++ {
		col := c.Col(i)
		if col.Name() != fmt.Sprint(i) {
			t.Error("Unexpected colname, got <", col.Type(), ">")
		} else if col.Type().Is(data.Uint|data.Nil|data.String) == false {
			t.Error("Unexpected type uint|nil|string for column, got <", col.Type(), ">")
		}
	}

	// Output table
	b := new(strings.Builder)
	if err := c.Write(b, c.OptHeader(), c.OptAscii(0, "")); err != nil {
		t.Error(err)
	} else {
		t.Log(b.String())
	}
}

func Test_Table_005(t *testing.T) {
	// Read table with existing columns and rows
	// allowing date and duration types and interpreting duration values in seconds
	// by default
	c := table.NewTable("A", "B")
	if err := c.Read(strings.NewReader(TABLE_C), c.OptType(data.Date|data.Duration)); err != nil {
		t.Error(err)
	}

	// Two headings
	for i := 0; i < 2; i++ {
		col := c.Col(i)
		switch i {
		case 0:
			if col.Name() != "A" {
				t.Error("Unexpected column name", col.Name())
			}
			if col.Type() != data.Date {
				t.Error("Unexpected column type", col.Type())
			}
		case 1:
			if col.Name() != "B" {
				t.Error("Unexpected column name", col.Name())
			}
			if col.Type() != data.Duration {
				t.Error("Unexpected column type", col.Type())
			}
		}
	}

	// Output table, truncating durations to hours
	b := new(strings.Builder)
	if err := c.Write(b, c.OptHeader(), c.OptAscii(0, ""), c.OptDuration(time.Hour)); err != nil {
		t.Error(err)
	} else {
		t.Log(b.String())
	}
}

func Test_Table_006(t *testing.T) {
	// Transform all float output values to %.2f
	transformFloat := func(t data.Table) data.TransformFunc {
		return func(i, j int, v interface{}) (interface{}, error) {
			if v, ok := v.(float64); ok {
				return fmt.Sprintf("%.2f", v), nil
			} else {
				return nil, data.ErrSkipTransform
			}
		}
	}

	// Read table with existing columns and rows
	c := table.NewTable()
	fh, err := os.Open(DATASET_A)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()
	if err := c.Read(fh, c.OptHeader(), c.OptType(data.DefaultTypes|data.Nil)); err != nil {
		t.Error(err)
	}

	// Output table, with float transform
	b := new(strings.Builder)
	if err := c.Write(b, c.OptHeader(), c.OptAscii(0, ""), c.OptTransform(transformFloat(c))); err != nil {
		t.Error(err)
	} else {
		t.Log(b.String())
	}
}

func Test_Table_007(t *testing.T) {
	// Read larger table with existing columns and rows
	c := table.NewTable()
	fh, err := os.Open(DATASET_B)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()
	if err := c.Read(fh, c.OptHeader(), c.OptType(data.DefaultTypes|data.Nil)); err != nil {
		t.Error(err)
	}

	// Summarize column data for each column
	var col data.TableCol
	for i := 0; ; i++ {
		if col = c.Col(i); col == nil {
			break
		} else if col.Type().Is(data.Float|data.Uint|data.Int) == false {
			continue
		} else if col.Count() == 0 {
			t.Error("Expected count, got", col.Count())
		} else if col.Sum() == 0 {
			t.Error("Expected sum, got", col.Sum())
		}
	}
}

func Test_Table_008(t *testing.T) {
	tests := []struct {
		in   data.Type
		out  data.Type
		null bool
	}{
		{data.Float, data.Float, false},
		{data.Uint, data.Uint, false},
		{data.Int, data.Int, false},
		{data.Float | data.Uint, data.Float, false},
		{data.Uint | data.Int, data.Int, false},
		{data.Float | data.Int | data.Uint, data.Float, false},
		{data.Float | data.Nil, data.Float, true},
		{data.Uint | data.Int | data.Nil, data.Int, true},
		{data.Float | data.Int | data.Uint | data.Nil, data.Float, true},
		{data.Float | data.Date | data.Nil, data.String, true},
		{data.Date | data.Nil, data.Date, true},
		{data.Other | data.Nil, data.Other, true},
		{0, data.Nil, false},
		{data.Nil, data.Nil, true},
	}
	for i, test := range tests {
		if x, y := test.in.Type(); x != test.out {
			t.Error(i, "Unexpected in=", test.in, " type out=", x)
		} else if y != test.null {
			t.Error(i, "Unexpected in=", test.in, " null out=", y)
		}
	}
}

func Test_Table_009(t *testing.T) {
	// Call a row iterator, reject any rows where there is a <nil>
	// leaving only one row
	rowIterator := func(c data.Table) data.IteratorFunc {
		return func(i int, r []interface{}) error {
			for _, v := range r {
				if v == nil {
					return data.ErrSkipTransform
				}
			}
			return nil
		}
	}
	// Read table with existing columns and rows
	c := table.NewTable()
	if err := c.Read(strings.NewReader(TABLE_A), c.OptType(data.Duration|data.Nil), c.OptRowIterator(rowIterator(c))); err != nil {
		t.Error(err)
	}
	// There should only be one remaining row
	if c.Len() != 1 {
		t.Error("Unexpected number of rows, expected 1, got", c.Len())
	}
}

func Test_Table_010(t *testing.T) {
	compareFunc := func(c data.Table) data.CompareFunc {
		return func(a, b []interface{}) bool {
			a_ := a[0].(time.Time)
			b_ := b[0].(time.Time)
			return b_.After(a_)
		}
	}

	// Read table with existing columns and rows
	c := table.NewTable()
	if err := c.Read(strings.NewReader(TABLE_C)); err != nil {
		t.Error(err)
	}

	// Sort table by column zero, which is a date
	c.Sort(compareFunc(c))

	// Ensure rows are in the right order
	var d time.Time
	for i := 0; i < c.Len(); i++ {
		e := c.Row(i)[0].(time.Time)
		if e.Before(d) {
			t.Error("Unexpected order for", e, "expected to be after", d)
		}
		d = e
	}
}

func Test_Table_011(t *testing.T) {
	// Read DATASET_A
	c := table.NewTable()
	fh, err := os.Open(DATASET_A)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()
	if err := c.Read(fh, c.OptHeader(), c.OptType(data.DefaultTypes|data.Nil)); err != nil {
		t.Error(err)
	}

	// Output table with width 40 and nice lines
	b := new(strings.Builder)
	if err := c.Write(b, c.OptHeader(), c.OptAscii(40, data.BorderLines)); err != nil {
		t.Error(err)
	} else {
		t.Log("\n" + b.String())
	}
}
func Test_Table_012(t *testing.T) {
	// Transform IP addresses
	transformAddrIn := func(t data.Table) data.TransformFunc {
		return func(i, j int, v interface{}) (interface{}, error) {
			if ip := net.ParseIP(v.(string)); ip == nil {
				return nil, data.ErrSkipTransform
			} else {
				return ip, nil
			}
		}
	}
	transformAddrOut := func(t data.Table) data.TransformFunc {
		return func(i, j int, v interface{}) (interface{}, error) {
			if v_, ok := v.(net.IP); ok {
				return fmt.Sprint("[", v_, "]"), nil
			} else {
				return nil, data.ErrSkipTransform
			}
		}
	}

	// Create an empty table
	c := table.NewTable("IP Address")
	if err := c.Read(strings.NewReader(TABLE_D), c.OptTransform(transformAddrIn(c))); err != nil {
		t.Error(err)
	}

	// Output table, adding [] around IP address
	b := new(strings.Builder)
	if err := c.Write(b, c.OptHeader(), c.OptTransform(transformAddrOut(c)), c.OptAscii(0, data.BorderLines)); err != nil {
		t.Error(err)
	} else {
		t.Log("\n" + b.String())
	}
}
