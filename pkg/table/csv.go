package table

import (
	"encoding/csv"
	"errors"
	"io"
)

/////////////////////////////////////////////////////////////////////
// CSV

func (t *Table) readCsv(r io.Reader, fn funcRowReader) error {
	// Create a CSV reader
	csv := csv.NewReader(r)

	// Apply delimiter
	if delim := t.opts.csvDelim; delim != 0 {
		csv.Comma = delim
	}

	// Iterate through rows
	var order []int
	var num int
	for {
		if row, err := csv.Read(); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		} else if t.hasOpt(optHeader) {
			order = t.readHeader(row)
			t.setOpt(optHeader, false)
			continue
		} else if err := t.readRow(num, order, row, fn); err != nil {
			return err
		}
		num++
	}

	// Return success
	return nil
}

func (t *Table) writeCsv(w io.Writer, fn funcRowWriter) error {
	// Create a CSV reader
	csv := csv.NewWriter(w)

	// Apply delimiter
	if delim := t.opts.csvDelim; delim != 0 {
		csv.Comma = delim
	}

	// Iterate through rows
	for i, r := range t.r {
		if i == 0 && t.hasOpt(optHeader) {
			// Output the header
		}
		if row, err := fn(i, r.row(t.header.w)); err != nil {
			return err
		} else if err := csv.Write(row); err != nil {
			return err
		}
	}

	// Flush
	csv.Flush()

	// Return success
	return nil
}
