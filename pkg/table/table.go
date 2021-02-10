package table

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type optFlag uint64

type Table struct {
	optFlag
	*header

	r      []*row
	w, h   uint
	d      time.Duration
	tz     *time.Location
	border string
}

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	optHeader optFlag = (1 << iota)
	optBool
	optDate
	optDatetime
	optDuration
	optUint
	optInt
	optFloat
	optNil
	optAscii
)

const (
	nilString = "<nil>"
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewTable(size data.Size) data.Table {
	t := new(Table)

	// Set width and height of the table
	t.w, t.h = sizeFor(size)

	// Set capacity and then width of the header
	t.header = NewHeader(t.w)
	t.header.setWidth(t.w)

	// Set capacity and then height of the rows
	t.r = make([]*row, 0, int(t.h))
	t.setHeight(t.h)

	// Return initial table
	return t
}

/////////////////////////////////////////////////////////////////////
// METHODS

func (t *Table) Read(r io.Reader, opts ...data.TableOpt) error {
	// Set option flags
	t.applyOpt(opts)

	// Create a CSV reader
	csv := csv.NewReader(r)

	// Iterate through rows
	var order []int
	for {
		if row, err := csv.Read(); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		} else if t.hasOpt(optHeader) {
			if order, err = t.doReadHeader(row); err != nil {
				return err
			} else {
				t.setOpt(optHeader, false)
			}
		} else if err := t.doReadRow(order, row); err != nil {
			return err
		}
	}

	// Return success
	return nil
}

func (t *Table) Write(w io.Writer, opts ...data.TableOpt) error {
	// Set option flags
	t.applyOpt(opts)

	// Return nil if no width or height
	if t.w == 0 || t.h == 0 {
		return nil
	}

	// Write Ascii
	switch {
	case t.hasOpt(optAscii):
		return t.writeAscii(w)
	default:
		return t.writeCsv(w)
	}
}

/////////////////////////////////////////////////////////////////////
// OPTIONS

func (t *Table) OptHeader() data.TableOpt {
	return func(t data.Table) {
		t.(*Table).setOpt(optHeader, true)
	}
}

func (t *Table) OptUint() data.TableOpt {
	return func(t data.Table) {
		t.(*Table).setOpt(optUint, true)
	}
}

func (t *Table) OptInt() data.TableOpt {
	return func(t data.Table) {
		t.(*Table).setOpt(optInt, true)
	}
}

func (t *Table) OptFloat() data.TableOpt {
	return func(t data.Table) {
		t.(*Table).setOpt(optFloat, true)
	}
}

func (t *Table) OptNil() data.TableOpt {
	return func(t data.Table) {
		t.(*Table).setOpt(optNil, true)
	}
}

func (t *Table) OptDuration(trunc time.Duration) data.TableOpt {
	return func(t data.Table) {
		t.(*Table).setOpt(optDuration, true)
		t.(*Table).d = trunc
	}
}

func (t *Table) OptDate(tz *time.Location) data.TableOpt {
	return func(t data.Table) {
		t.(*Table).setOpt(optDate, true)
		if tz == nil {
			if t.(*Table).tz == nil {
				t.(*Table).tz = time.Local
			}
		} else {
			t.(*Table).tz = tz
		}
	}
}

func (t *Table) OptDatetime(tz *time.Location) data.TableOpt {
	return func(t data.Table) {
		t.(*Table).setOpt(optDatetime, true)
		if tz == nil {
			if t.(*Table).tz == nil {
				t.(*Table).tz = time.Local
			}
		} else {
			t.(*Table).tz = tz
		}
	}
}

func (t *Table) OptBool() data.TableOpt {
	return func(t data.Table) {
		t.(*Table).setOpt(optBool, true)
	}
}

func (t *Table) OptAscii(width uint, border string) data.TableOpt {
	return func(t data.Table) {
		t.(*Table).setOpt(optAscii, true)
		if border == "" {
			t.(*Table).border = data.BorderDefault
		} else {
			t.(*Table).border = border
		}
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (t *Table) String() string {
	str := "<table"
	str += fmt.Sprintf(" size={ %v,%v }", t.w, t.h)
	str += fmt.Sprint(" ", t.header)
	for i, r := range t.r {
		str += fmt.Sprint(" ", i, ":", r)
	}
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (t *Table) applyOpt(opts []data.TableOpt) {
	t.optFlag = 0
	t.d = time.Nanosecond
	t.tz = time.Local
	for _, opt := range opts {
		opt(t)
	}
}

func (t *Table) hasOpt(f optFlag) bool {
	return t.optFlag&f == f
}

func (t *Table) setOpt(f optFlag, state bool) {
	if state {
		t.optFlag |= f
	} else {
		t.optFlag ^= f
	}
}

func (t *Table) doReadHeader(row []string) ([]int, error) {
	order := make([]int, 0, len(row))
	for _, field := range row {
		fieldName := strings.TrimSpace(field)
		j, err := t.appendField(fieldName)
		if errors.Is(err, data.ErrDuplicateEntry) {
			// Ignore error
		} else if err != nil {
			return nil, err
		}
		order = append(order, j)
	}
	return order, nil
}

// doReadRow extends the width and height of the table before
// creating a row and appending to table
func (t *Table) doReadRow(order []int, row []string) error {
	width := len(row)
	for _, i := range order {
		if i+1 > width {
			width = i + 1
		}
	}
	if order != nil && len(order) != len(row) {
		return data.ErrInternalAppError
	}

	// Set header width
	t.setWidth(uint(width))

	// Create a row with "width" capacity
	r := NewRow(uint(width))

	// Set contents of row
	for i, cell := range row {
		if order == nil {
			r.set(uint(i), cell, t.valueForString(cell))
		} else {
			r.set(uint(order[i]), cell, t.valueForString(cell))
		}
	}

	// Append row
	t.appendRow(r)

	// Return success
	return nil
}

func (t *Table) appendRow(r *row) {
	t.r = append(t.r, r)
	t.h = uint(len(t.r))
}

// setHeight sets the absolute height of the table, adding and
// removing rows as necessary
func (t *Table) setHeight(height uint) {
	l := len(t.r)
	if l > int(height) {
		t.r = t.r[:height]
	}
	if l == int(height) {
		return
	}
	for len(t.r) < int(height) {
		t.r = append(t.r, NewRow(t.w))
	}
}

// setWidth sets the absolute width of the table. adding and
// removing columns as necessary
func (t *Table) setWidth(width uint) {
	if len(t.header.order) < int(width) {
		t.header.setWidth(width)
	}
	t.w = width
}

// sizeFor returns integer width and height from size
func sizeFor(size data.Size) (uint, uint) {
	return uint(math.Abs(float64(size.W))), uint(math.Abs(float64(size.H)))
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS - PARSE VALUES

// valueForString returns interpreted value from a string
func (t *Table) valueForString(v string) interface{} {
	// Check for nil values
	if t.hasOpt(optNil) && nilValueForString(v) {
		return nil
	}
	// Check for uint values
	if t.hasOpt(optUint) {
		if v_, exists := uintValueForString(v); exists {
			return v_
		}
	}
	// Check for int values
	if t.hasOpt(optInt) {
		if v_, exists := intValueForString(v); exists {
			return v_
		}
	}
	// Check for float values
	if t.hasOpt(optFloat) {
		if v_, exists := floatValueForString(v); exists {
			return v_
		}
	}
	// Check for datetime
	if t.hasOpt(optDatetime) {
		if v_, exists := datetimeValueForString(v, t.tz); exists {
			return v_
		}
	}
	// Check for date values
	if t.hasOpt(optDate) {
		if v_, exists := dateValueForString(v, t.tz); exists {
			return v_
		}
	}
	// Check for duration values
	if t.hasOpt(optDuration) {
		if v_, exists := durationValueForString(v); exists {
			return v_.Truncate(t.d)
		}
		if v_, exists := intValueForString(v); exists {
			return time.Duration(v_) * t.d
		}
	}
	// Check for bool values
	if t.hasOpt(optBool) {
		if v_, exists := boolValueForString(v); exists {
			return v_
		}
	}
	// Return string by default
	return v
}

func nilValueForString(v string) bool {
	return v == "" || strings.TrimSpace(v) == ""
}

func floatValueForString(v string) (float64, bool) {
	f, err := strconv.ParseFloat(v, 64)
	return f, err == nil
}

func uintValueForString(v string) (uint64, bool) {
	u, err := strconv.ParseUint(v, 0, 64)
	return u, err == nil
}

func intValueForString(v string) (int64, bool) {
	i, err := strconv.ParseInt(v, 0, 64)
	return i, err == nil
}

func boolValueForString(v string) (bool, bool) {
	b, err := strconv.ParseBool(v)
	return b, err == nil
}

func durationValueForString(v string) (time.Duration, bool) {
	d, err := time.ParseDuration(v)
	return d, err == nil
}

func datetimeValueForString(v string, tz *time.Location) (time.Time, bool) {
	if t, err := time.ParseInLocation(time.RFC3339, v, tz); err == nil {
		return t, true
	}
	if t, err := time.ParseInLocation(time.UnixDate, v, tz); err == nil {
		return t, true
	}
	if t, err := time.ParseInLocation(time.RFC822, v, tz); err == nil {
		return t, true
	}
	if t, err := time.ParseInLocation("2006-01-02 15:04:05", v, tz); err == nil {
		return t, true
	}
	if t, err := time.ParseInLocation("2006-01-02 15:04", v, tz); err == nil {
		return t, true
	}
	return time.Time{}, false
}

func dateValueForString(v string, tz *time.Location) (time.Time, bool) {
	if t, err := time.ParseInLocation("2006-01-02", v, tz); err == nil {
		return t, true
	}
	if t, err := time.ParseInLocation("2006/01/02", v, tz); err == nil {
		return t, true
	}
	if t, err := time.ParseInLocation("02/01/2006", v, tz); err == nil {
		return t, true
	}
	if t, err := time.ParseInLocation("02-01-2006", v, tz); err == nil {
		return t, true
	}
	if t, err := time.ParseInLocation("Jan 2 2006", v, tz); err == nil {
		return t, true
	}
	if t, err := time.ParseInLocation("2 Jan 2006", v, tz); err == nil {
		return t, true
	}
	if t, err := time.ParseInLocation("Jan 2 06", v, tz); err == nil {
		return t, true
	}
	if t, err := time.ParseInLocation("2 Jan 06", v, tz); err == nil {
		return t, true
	}
	return time.Time{}, false
}
