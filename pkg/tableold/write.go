package table

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"
)

/////////////////////////////////////////////////////////////////////
// METHODS

func (t *Table) writeAscii(io.Writer) error {
	return nil
}

func (t *Table) writeCsv(w io.Writer) error {
	enc := csv.NewWriter(w)

	if t.hasOpt(optHeader) {
		if err := enc.Write(t.header.stringsForWidth(t.w)); err != nil {
			return err
		}
	}
	for _, row := range t.r {
		if err := enc.Write(t.stringsForRow(row)); err != nil {
			return err
		}
	}
	enc.Flush()

	// Return success
	return nil
}

func (t *Table) stringsForRow(r *row) []string {
	result := make([]string, t.w)
	for i := range result {
		if i < len(r.v) {
			result[i] = t.stringForValue(r.v[i], r.s[i])
		} else {
			result[i] = t.stringForValue(nil, "")
		}
	}
	return result
}

/////////////////////////////////////////////////////////////////////
// CONVERT VALUE TO STRING

func (t *Table) stringForValue(v interface{}, s string) string {
	if t.hasOpt(optNil) {
		if v == nil {
			return nilString
		}
	}
	if t.hasOpt(optFloat) {
		switch v.(type) {
		case float32, float64:
			return fmt.Sprint(v)
		}
	}
	if t.hasOpt(optUint) {
		switch v.(type) {
		case uint8, uint16, uint32, uint64, uint:
			return fmt.Sprint(v)
		}
	}
	if t.hasOpt(optInt) {
		switch v.(type) {
		case int8, int16, int32, int64, int:
			return fmt.Sprint(v)
		}
	}
	if t.hasOpt(optBool) {
		switch v.(type) {
		case bool:
			return fmt.Sprint(v)
		}
	}
	if t.hasOpt(optDuration) {
		switch v.(type) {
		case time.Duration:
			return fmt.Sprint(v)
		}
	}
	return s
}
