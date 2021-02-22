package dtd_test

import (
	"os"
	"testing"

	"github.com/djthorpe/data/pkg/dtd"
)

const (
	DTDFILE_A = "../../etc/dtd/svg.dtd"
)

func Test_DTD_001(t *testing.T) {
	d := dtd.NewDocument("test")
	if d == nil {
		t.Fatal("Unexpected nil return from NewDTD")
	}

	// Output
	t.Log(d)
}

func Test_DTD_002(t *testing.T) {
	fh, err := os.Open(DTDFILE_A)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()
	if dtd, err := dtd.Read(fh); err != nil {
		t.Fatal(err)
	} else {
		t.Log(dtd)
	}
}
