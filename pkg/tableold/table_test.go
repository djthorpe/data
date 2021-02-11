package table_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/table"
)

const (
	DATASET_CSV = "../../etc/dataset/12-15-2020.csv"
)

func Test_Table_001(t *testing.T) {
	c := table.NewTable(data.ZeroSize)
	t.Log(c)
}

func Test_Table_002(t *testing.T) {
	c := table.NewTable(data.Size{5, 5})
	t.Log(c)
}

func Test_Table_003(t *testing.T) {
	c := table.NewTable(data.ZeroSize)
	fh, err := os.Open(DATASET_CSV)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()
	// Read a table which can include floats, datetimes and nils in addition to string
	if err := c.Read(fh, c.OptHeader(), c.OptFloat(), c.OptDatetime(nil), c.OptNil()); err != nil {
		t.Error(err)
	} else if err := c.Write(os.Stdout, c.OptHeader(), c.OptNil()); err != nil {
		t.Error(err)
	}
}

func Test_Table_004(t *testing.T) {
	c := table.NewTable(data.ZeroSize)
	r := strings.NewReader(
		`0,1
1,2
1,
,2
`,
	)
	if err := c.Read(r, c.OptFloat(), c.OptUint(), c.OptNil()); err != nil {
		t.Error(err)
	} else if err := c.Write(os.Stdout); err != nil {
		t.Error(err)
	}
}

func Test_Table_005(t *testing.T) {
	c := table.NewTable(data.ZeroSize)
	r := strings.NewReader(
		`0,1
1,2
1,
,2
`,
	)
	if err := c.Read(r, c.OptBool(), c.OptNil()); err != nil {
		t.Error(err)
	} else if err := c.Write(os.Stdout, c.OptBool(), c.OptNil()); err != nil {
		t.Error(err)
	}
}

func Test_Table_006(t *testing.T) {
	c := table.NewTable(data.ZeroSize)
	r := strings.NewReader(
		`0,1
1,2
1,
,2
`,
	)
	if err := c.Read(r, c.OptDuration(time.Millisecond), c.OptNil()); err != nil {
		t.Error(err)
	} else if err := c.Write(os.Stdout, c.OptHeader(), c.OptNil()); err != nil {
		t.Error(err)
	}
}

func Test_Table_007(t *testing.T) {
	c := table.NewTable(data.ZeroSize)
	r := strings.NewReader(
		`0,1
1,2
1,
,2
`,
	)
	if err := c.Read(r, c.OptNil()); err != nil {
		t.Error(err)
	}
	c.ForArray(func(i int, a []interface{}) {
		t.Log(i, "=>", a)
	}, c.OptNil())
}
