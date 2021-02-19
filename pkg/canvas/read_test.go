package canvas_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	data "github.com/djthorpe/data"
	canvas "github.com/djthorpe/data/pkg/canvas"
)

const (
	SVGFILE_A = "../../etc/xml/aa.svg"
	SVGFILE_B = "../../etc/xml/acid.svg"
	SVGFILE_C = "../../etc/xml/adobe.svg"
	SVGFILE_D = "../../etc/xml/android.svg"
	SVGFILE_E = "../../etc/xml/car.svg"
	SVGFILE_F = "../../etc/xml/tiger.svg"
)

func Test_Read_001(t *testing.T) {
	c1 := canvas.NewCanvas(data.Size{16, 16}, data.PX).Version("1.1").Title("Hello, World")

	// Write SVG
	b := new(strings.Builder)
	if err := c1.Write(data.SVG, b); err != nil {
		t.Fatal(err)
	}

	// Read SVG
	if c2, err := canvas.Read(data.SVG, strings.NewReader(b.String())); err != nil {
		t.Error(err)
	} else if str := fmt.Sprint(c2); str != `<svg xmlns="http://www.w3.org/2000/svg" width="16px" height="16px" viewBox="0 0 16 16" version="1.1"><title>Hello, World</title></svg>` {
		t.Error("Unexpected return: ", str)
	} else {
		t.Log(c2)
	}
}

func Test_Read_002(t *testing.T) {
	fh, err := os.Open(SVGFILE_A)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()

	if c1, err := canvas.Read(data.SVG, fh); err != nil {
		t.Fatal(err)
	} else {
		t.Log(c1)
	}
}

func Test_Read_003(t *testing.T) {
	fh, err := os.Open(SVGFILE_B)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()

	if c1, err := canvas.Read(data.SVG, fh); err != nil {
		t.Fatal(err)
	} else {
		t.Log(c1)
	}
}

func Test_Read_004(t *testing.T) {
	fh, err := os.Open(SVGFILE_C)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()

	if c1, err := canvas.Read(data.SVG, fh); err != nil {
		t.Fatal(err)
	} else {
		t.Log(c1)
	}
}

func Test_Read_005(t *testing.T) {
	fh, err := os.Open(SVGFILE_D)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()

	if c1, err := canvas.Read(data.SVG, fh); err != nil {
		t.Fatal(err)
	} else {
		t.Log(c1)
	}
}

func Test_Read_006(t *testing.T) {
	fh, err := os.Open(SVGFILE_E)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()

	if c1, err := canvas.Read(data.SVG, fh); err != nil {
		t.Fatal(err)
	} else {
		t.Log(c1)
	}
}

func Test_Read_007(t *testing.T) {
	fh, err := os.Open(SVGFILE_F)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()

	if c1, err := canvas.Read(data.SVG, fh); err != nil {
		t.Fatal(err)
	} else {
		t.Log(c1)
	}
}
