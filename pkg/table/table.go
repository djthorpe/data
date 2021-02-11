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

func NewTable(size data.Size) data.Table {
	t := new(Table)

	// Set width and height of the table, if error is returned
	// then set nil
	if w, h, err := sizeFor(size); err != nil {
		return nil
	} else {
		t.header = NewHeader(w)

		// Append any existing rows
		t.r = make([]*row, h)
		for i := 0; i < h; i++ {
			t.r[i] = NewRow(nil)
		}
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
			if err := t.rowIterator(i, row.v); err == nil {
				t.r = append(t.r, row)
			} else if errors.Is(err, data.ErrSkipTransform) == false {
				return err
			}
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

// Col returns column information for a zero-indexed table column
func (t *Table) Col(i int) data.TableCol {
	if i < 0 || i >= t.header.w {
		return nil
	} else {
		return t.header.col(i)
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

func (t *Table) Append(row ...interface{}) {
	// Extend header width as necessary
	t.header.w = maxInt(t.header.w, len(row))
	// Append the row
	t.r = append(t.r, NewRow(row))
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

// sizeFor returns w,h as integers, and returns error
// if the size does not contain positive (or zero) integers
func sizeFor(size data.Size) (int, int, error) {
	if size.W < 0 || size.H < 0 {
		return 0, 0, data.ErrBadParameter
	}
	if w, h := int(size.W), int(size.H); float32(w) != size.W || float32(h) != size.H {
		return 0, 0, data.ErrBadParameter
	} else {
		return w, h, nil
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
	// Extend header width as necessary
	if order == nil {
		t.header.w = maxInt(t.header.w, len(row))
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

func maxInt(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
