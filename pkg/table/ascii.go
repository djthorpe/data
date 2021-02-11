package table

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type asciicol struct {
	w   int // Col width
	a   int // Col alignment
	fmt string
}

const (
	borderNW = iota
	borderN
	borderNE
	borderW
	borderO
	borderE
	borderSW
	borderS
	borderSE
	borderV
	borderH
)

const (
	alignLeft = iota
	alignRight
)

/////////////////////////////////////////////////////////////////////
// ASCII

func (t *Table) writeAscii(w io.Writer, fn funcRowWriter) error {
	cols := make([]*asciicol, t.header.w)
	for i := range cols {
		cols[i] = &asciicol{15, alignRight, ""}
	}

	// Write top line
	if err := t.asciiTop(w, cols); err != nil {
		return err
	}

	// Iterate through rows
	for i, r := range t.r {
		if t.hasOpt(optHeader) {
			if err := t.asciiHeader(w, cols, t.header.row()); err != nil {
				return err
			}
			t.setOpt(optHeader, false)
		}
		if row, err := fn(i, r.row(t.header.w)); err != nil {
			return err
		} else if err := t.asciiRow(w, cols, row); err != nil {
			return err
		}
	}

	// Write bottom line
	if err := t.asciiBottom(w, cols); err != nil {
		return err
	}

	// Return success
	return nil
}

func (t *Table) asciiTop(w io.Writer, cols []*asciicol) error {
	row := make([]string, len(cols))
	for i, col := range cols {
		row[i] = strings.Repeat("-", col.w)
	}
	return t.asciiLine([]byte("+-"), []byte("-+-"), []byte("-+"), w, cols, row)
}

func (t *Table) asciiRow(w io.Writer, cols []*asciicol, row []string) error {
	return t.asciiLine([]byte("| "), []byte(" | "), []byte(" |"), w, cols, row)
}

func (t *Table) asciiHeader(w io.Writer, cols []*asciicol, row []string) error {
	if err := t.asciiLine([]byte("| "), []byte(" | "), []byte(" |"), w, cols, row); err != nil {
		return err
	}
	return t.asciiTop(w, cols)
}

func (t *Table) asciiBottom(w io.Writer, cols []*asciicol) error {
	row := make([]string, len(cols))
	for i, col := range cols {
		row[i] = strings.Repeat("-", col.w)
	}
	return t.asciiLine([]byte("+-"), []byte("-+-"), []byte("-+"), w, cols, row)
}

func (t *Table) asciiLine(l, m, r []byte, w io.Writer, cols []*asciicol, row []string) error {
	row_ := make([][]byte, len(cols))
	for i, col := range cols {
		row_[i] = []byte(fmt.Sprintf(col.Fmt(), row[i])[:col.w])
	}
	if _, err := w.Write(l); err != nil {
		return err
	}
	if _, err := w.Write(bytes.Join(row_, m)); err != nil {
		return err
	}
	if _, err := w.Write(r); err != nil {
		return err
	}
	if _, err := w.Write([]byte("\n")); err != nil {
		return err
	}
	return nil
}

func (c *asciicol) Fmt() string {
	if c.fmt == "" {
		if c.w == 0 {
			return "%s"
		}
		switch c.a {
		case alignLeft:
			c.fmt = fmt.Sprint("%-", c.w, "s")
		case alignRight:
			c.fmt = fmt.Sprint("%+", c.w, "s")
		default:
			c.fmt = fmt.Sprint("%*", c.w, "s")
		}
	}
	return c.fmt
}
