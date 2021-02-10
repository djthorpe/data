package data

import (
	"io"
	"time"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type TableOpt func(Table)
type MapIterator func(int, map[string]interface{})
type ArrayIterator func(int, []interface{})

/////////////////////////////////////////////////////////////////////
// INTERFACES

type Table interface {
	// Read CSV data with table options
	Read(io.Reader, ...TableOpt) error

	// Write data with table options
	Write(io.Writer, ...TableOpt) error

	// Iterate through rows returning a map of values
	ForMap(MapIterator, ...TableOpt)

	// Iterate through rows returning an array of values
	ForArray(ArrayIterator, ...TableOpt)

	// OptHeader used on Read to indicate there is a CSV header and
	// with Write to output header in addition to data. Ignored for
	// ForMap and ForArray
	OptHeader() TableOpt

	// OptAscii used on Write to output ASCII table instead of CSV. Option
	// is ignored for Read. Arguments are maximum table width and the border
	// characters (data.BorderDefault or data.BorderLines)
	OptAscii(uint, string) TableOpt

	// OptFloat used on Read, Write and For to interpret floats
	OptFloat() TableOpt

	// OptDuration used on Read to indicate some cells are durations (h,m,s,ms,ns)
	// and truncate to the provided duration
	OptDuration(time.Duration) TableOpt

	// OptDate used on Read to indicate some cells are date (YYYY-MM-DD, DD-MM-YYYY)
	// with timezone. If timezone is nil, current local timezone is used.
	OptDate(tz *time.Location) TableOpt

	// OptDatetime used on Read to indicate some cells are datetime
	// with timezone. If timezone is nil, current local timezone is used.
	OptDatetime(tz *time.Location) TableOpt

	// OptBool used on Read to indicate some cells are booleans (t,f,T,F,true,false)
	OptBool() TableOpt

	// OptUint used on Read to indicate some cells are uint values
	OptUint() TableOpt

	// OptInt used on Read to indicate some cells are int values
	OptInt() TableOpt

	// OptNil used on Read to indicate empty cells should be interpreted
	// as nil values, and on Write to write out <nil> in cells with no value
	OptNil() TableOpt
}

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	BorderDefault = "+++++++++|-"
	BorderLines   = "┌┬┐├┼┤└┴┘│─"
)
