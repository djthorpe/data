package table

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Table struct {
	opts struct {
		d, o       optFlag
		tz         *time.Location
		dur        time.Duration
		border     []rune
		asciiwidth int
		csvDelim   rune
		transform  []data.TransformFunc
		iterator   data.IteratorFunc
		compare    data.CompareFunc
	}
	*header
	r []*row
}

type funcRowReader func(int, []string) error
type funcRowWriter func(int, []interface{}) ([]string, error)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewTable(cols ...string) data.Table {
	t := new(Table)
	t.header = NewHeader(len(cols))
	t.opts.d = optNil | optUint | optInt | optFloat | optBool | optDuration | optDate | optDatetime

	// Append columns
	for _, col := range cols {
		t.header.append(col)
	}

	// Return initial table
	return t
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (t *Table) Read(r io.Reader, opts ...data.TableOpt) error {
	// Set option flags
	t.applyOpt(opts)

	// Perform read
	switch {
	case t.hasOpt(optCsv):
		fallthrough
	default:
		return t.readCsv(r, func(i int, values []string) error {
			row := NewRow(make([]interface{}, len(values)))
			for j, v := range values {
				if v_, err := t.inValue(i, j, v); err != nil {
					return err
				} else {
					row.v[j] = v_
				}
			}
			// Call row iterator
			if err := t.rowIterator(i, row.v); errors.Is(err, data.ErrSkipTransform) {
				return nil
			} else if err != nil {
				return err
			}
			// Validate values and re-scan if the width of the table has changed
			if rescan := t.header.validate(row); rescan {
				t.validate()
			}
			// Append the row
			t.r = append(t.r, row)
			// Return success
			return nil
		})
	}
}

// Write data with table options
func (t *Table) Write(w io.Writer, opts ...data.TableOpt) error {
	// Set option flags
	t.applyOpt(opts)

	// Return nil if no width or height
	if len(t.r) == 0 || t.header.w == 0 {
		return nil
	}

	// Perform write
	switch {
	case t.hasOpt(optAscii):
		return t.writeAscii(w, func(i int, row []interface{}) ([]string, error) {
			result := make([]string, len(row))
			for j, v := range row {
				if v_, err := t.outValue(i, j, v); err != nil {
					return nil, err
				} else if v__, ok := v_.(string); ok {
					result[j] = v__
				} else {
					result[j] = fmt.Sprint(v_)
				}
			}
			return result, nil
		})
	case t.hasOpt(optCsv):
		fallthrough
	default:
		return t.writeCsv(w, func(i int, row []interface{}) ([]string, error) {
			result := make([]string, len(row))
			for j, v := range row {
				if v_, err := t.outValue(i, j, v); err != nil {
					return nil, err
				} else if v__, ok := v_.(string); ok {
					result[j] = v__
				} else {
					result[j] = fmt.Sprint(v_)
				}
			}
			return result, nil
		})
	}
}

// Stream data with table options
/*
func (t *Table) Stream(w io.Writer, r io.Reader, opts ...data.TableOpt) error {
	// Set option flags
	t.applyOpt(opts)

	// Return nil if no width or height
	if len(t.r) == 0 || t.header.w == 0 {
		return nil
	}

	// Not yet implemented
	return nil
}
*/

// Col returns column information for a zero-indexed table column
func (t *Table) Col(i int) data.TableCol {
	c := t.header.col(i)
	if c == nil {
		return nil
	}
	// Summarize data for column
	c.group(nil)
	for _, row := range t.r {
		if len(row.v) < c.i || row.v[c.i] == nil {
			continue
		} else {
			c.group(row.v[c.i])
		}
	}

	// Return column
	return c
}

// Row returns values for a zero-indexed table row
func (t *Table) Row(i int) []interface{} {
	if i < 0 || i >= len(t.r) {
		return nil
	} else {
		return t.r[i].row(t.header.w)
	}
}

// Sort does a quicksort on the table using a comparison function
func (t *Table) Sort(fn data.CompareFunc) {
	if fn != nil {
		t.opts.compare = fn
		sort.Sort(t)
	}
}

func (t *Table) Len() int {
	return len(t.r)
}

func (t *Table) Less(i, j int) bool {
	return t.opts.compare(t.r[i].row(t.header.w), t.r[j].row(t.header.w))
}

func (t *Table) Swap(i, j int) {
	t.r[i], t.r[j] = t.r[j], t.r[i]
}

func (t *Table) Append(values ...interface{}) {
	row := NewRow(values)
	if rescan := t.header.validate(row); rescan {
		t.validate()
	}
	t.r = append(t.r, row)
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (t *Table) String() string {
	str := "<table"
	str += " " + fmt.Sprint(t.header)
	for i, r := range t.r {
		str += fmt.Sprint(" ", i, ":", r)
	}
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// validate rescans the table for types
func (t *Table) validate() {
	for _, r := range t.r {
		t.header.validate(r)
	}
}

// readHeader adds columns to the table and returns the order of the columns
// for subsequent reads
func (t *Table) readHeader(row []string) []int {
	return t.header.set(row)
}

// readRow reorders the row in the correct order and uses callback
// to either stream the row out or store in the table
func (t *Table) readRow(i int, order []int, row []string, fn funcRowReader) error {
	// Check incoming parameters
	if fn == nil {
		return data.ErrInternalAppError
	} else if len(order) != 0 && len(order) != len(row) {
		return data.ErrInternalAppError
	}
	// Re-order row as necessary
	if len(order) == 0 {
		return fn(i, row)
	} else {
		r := make([]string, t.header.w)
		for i, value := range row {
			r[order[i]] = value
		}
		return fn(i, r)
	}
}
