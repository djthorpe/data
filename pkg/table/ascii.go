package table

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

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

/////////////////////////////////////////////////////////////////////
// ASCII

func (t *Table) writeAscii(w io.Writer, fn funcRowWriter) error {
	cols := t.header.cols()
	fmt.Println(cols)
	// Write top line
	if err := t.asciiTop(w, cols); err != nil {
		return err
	}

	// Iterate through rows
	for i, r := range t.r {
		if t.hasOpt(optHeader) {
			header := make([]string, t.header.w)
			for i, col := range t.header.cols() {
				header[i] = col.Name()
			}
			if err := t.asciiHeader(w, cols, header); err != nil {
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

func (t *Table) asciiTop(w io.Writer, cols []*col) error {
	row := make([]string, len(cols))
	for i, col := range cols {
		row[i] = strings.Repeat("-", col.asciiWidth())
	}
	return t.asciiLine([]byte("+-"), []byte("-+-"), []byte("-+"), w, cols, row)
}

func (t *Table) asciiRow(w io.Writer, cols []*col, row []string) error {
	return t.asciiLine([]byte("| "), []byte(" | "), []byte(" |"), w, cols, row)
}

func (t *Table) asciiHeader(w io.Writer, cols []*col, row []string) error {
	if err := t.asciiLine([]byte("| "), []byte(" | "), []byte(" |"), w, cols, row); err != nil {
		return err
	}
	return t.asciiTop(w, cols)
}

func (t *Table) asciiBottom(w io.Writer, cols []*col) error {
	row := make([]string, len(cols))
	for i, col := range cols {
		row[i] = strings.Repeat("-", col.asciiWidth())
	}
	return t.asciiLine([]byte("+-"), []byte("-+-"), []byte("-+"), w, cols, row)
}

func (t *Table) asciiLine(l, m, r []byte, w io.Writer, cols []*col, row []string) error {
	row_ := make([][]byte, len(cols))
	for i, col := range cols {
		row_[i] = []byte(fmt.Sprintf(col.Fmt(), row[i])[:col.asciiWidth()])
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
