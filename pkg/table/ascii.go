package table

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

/////////////////////////////////////////////////////////////////////
// TYPES

const (
	borderNW = iota
	borderN
	borderNE
	borderW
	borderC
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

func (t *Table) asciiRow(w io.Writer, cols []*col, row []string) error {
	return t.asciiLine(borderV, borderV, borderV, w, cols, row)
}

func (t *Table) asciiHeader(w io.Writer, cols []*col, row []string) error {
	if err := t.asciiLine(borderV, borderV, borderV, w, cols, row); err != nil {
		return err
	}
	return t.asciiDivider(borderW, borderC, borderE, w, cols)
}

func (t *Table) asciiTop(w io.Writer, cols []*col) error {
	return t.asciiDivider(borderNW, borderN, borderNE, w, cols)
}

func (t *Table) asciiBottom(w io.Writer, cols []*col) error {
	return t.asciiDivider(borderSW, borderS, borderSE, w, cols)
}

func (t *Table) asciiDivider(l, m, r int, w io.Writer, cols []*col) error {
	row := make([]string, len(cols))
	fill := t.asciiChar(borderH)
	for i, col := range cols {
		row[i] = strings.Repeat(string(fill), col.asciiWidth())
	}
	return t.asciiLine(l, m, r, w, cols, row)
}

func (t *Table) asciiLine(l, m, r int, w io.Writer, cols []*col, row []string) error {
	for i, v := range row {
		col := cols[i]
		if utf8.RuneCountInString(v) != col.asciiWidth() {
			v = fmt.Sprintf(col.Fmt(), v)[:col.asciiWidth()]
		}
		left := t.asciiChar(m)
		right := []byte{}
		if i == 0 {
			left = t.asciiChar(l)
		}
		if i == len(row)-1 {
			right = append(t.asciiChar(r), byte('\n'))
		}
		if _, err := w.Write(left); err != nil {
			return err
		} else if _, err := w.Write([]byte(v)); err != nil {
			return err
		} else if _, err := w.Write(right); err != nil {
			return err
		}
	}
	return nil
}

func (t *Table) asciiChar(i int) []byte {
	// Rune for character
	r := t.opts.border[i]
	return []byte(string(r))
}
