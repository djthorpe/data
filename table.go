package data

import (
	"io"
	"time"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type TableOpt func(Table)
type Type uint
type TableCellFlag uint
type TransformFunc func(int, int, interface{}) (interface{}, error)
type IteratorFunc func(int, []interface{}) error
type CompareFunc func(a, b []interface{}) bool

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	Nil Type = (1 << iota)
	Uint
	Int
	Float
	Duration
	Date
	Datetime
	Bool
	String       Type = 0
	DefaultTypes      = Uint | Int | Float | Duration | Date | Datetime | Bool
)

const (
	Bold TableCellFlag = (1 << iota)
)

const (
	BorderDefault = "+++++++++|-"
	BorderLines   = "┌┬┐├┼┤└┴┘│─"
)

/////////////////////////////////////////////////////////////////////
// INTERFACES

type Table interface {
	// Read CSV data with table options
	Read(io.Reader, ...TableOpt) error

	// Write data with table options
	Write(io.Writer, ...TableOpt) error

	// Append a row to the table
	Append(...interface{})

	// Len returns the number of rows
	Len() int

	// Col returns column information for a zero-indexed table column
	Col(int) TableCol

	// Cell returns cell format given a cell string
	//Cell(string, TableCellFlag) TableCell

	// Sort sorts the rows using a comparison function, which should return
	// true if the first argument is less than the second argument
	Sort(CompareFunc)

	// OptHeader used on Read to indicate there is a CSV header and
	// with Write to output header in addition to data. Ignored for
	// ForMap and ForArray
	OptHeader() TableOpt

	// OptTransform used on Read, or Write transforms a value.
	// Several transform functions can be used in series on a value.
	// Transformation functions are called in series until nil or
	// error returned. If ErrSkipTransform is returned, the next
	// transformation is tried
	OptTransform(...TransformFunc) TableOpt

	// OptType used on Read to indicate transformation of values into
	// native types. If a type cannot be transformed, string is used.
	// Use "DefaultTypes" for all types, and "Nil" to interpret empty strings
	// as nil values. For example:
	//
	//   t := data.NewTable(data.ZeroSize)
	//   t.Read(os.Stdin,data.OptType(data.DefaultTypes|data.Nil))
	//
	OptType(Type) TableOpt

	// OptAscii used on Write to output ASCII table instead of CSV. Option
	// is ignored for Read. Arguments are maximum table width and the border
	// characters (data.BorderDefault or data.BorderLines). Setting width
	// to zero makes unconstrained width, setting string to empty sets default
	// table line output
	OptAscii(uint, string) TableOpt

	// OptCsv used to Read or Write CSV files, with delimiter character, or if
	// zero then comma is used
	OptCsv(rune) TableOpt

	// OptDuration used on Read to interpret values into durations (h,m,s,ms,ns)
	// and truncate to the provided duration
	OptDuration(time.Duration) TableOpt

	// OptTimezone used on Read to set timezone for dates and times which do not
	// include timezone explicitly. If timezone is nil, current local timezone is used.
	OptTimezone(tz *time.Location) TableOpt

	// OptRowIterator called on read before appending row to table. If error returned
	// by iterator function is ErrSkipTransform then the row is not appended to the table
	OptRowIterator(IteratorFunc) TableOpt
}

// TableCol represents information about a table
// column
type TableCol interface {
	Name() string
	Type() Type
}
