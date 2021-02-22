package data

import (
	"io"
	"strings"
	"time"

	f32 "github.com/djthorpe/data/pkg/f32"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type TableOpt func(Table)
type Type uint
type TransformFunc func(int, int, interface{}) (interface{}, error)
type IteratorFunc func(int, []interface{}) error
type CompareFunc func(a, b []interface{}) bool

// type TableCellFlag uint TODO

/////////////////////////////////////////////////////////////////////
// INTERFACES

type Table interface {
	// Read CSV data with table options
	Read(io.Reader, ...TableOpt) error

	// Write data with table options
	Write(io.Writer, ...TableOpt) error

	// Output XML of the table
	// DOM(...TableOpt) Document // TODO

	// Append a row to the table
	Append(...interface{})

	// Len returns the number of rows
	Len() int

	// Row returns values for a zero-indexed table row
	// or nil if the row does not exist
	Row(int) []interface{}

	// Col returns column information for a zero-indexed table column
	// or nil if the column doesn't exist
	Col(int) TableCol

	// Cell returns cell format given a cell string
	//Cell(string, TableCellFlag) TableCell

	// Sort sorts the rows using a comparison function, which should return
	// true if the first argument is less than the second argument
	Sort(CompareFunc)

	// OptHeader used on Read to indicate there is a CSV header and
	// with Write and DOM to output header in addition to data. Ignored for
	// ForMap and ForArray
	OptHeader() TableOpt

	// OptTransform used on Read, Write or DOM to transform a value.
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

	// OptSql used to Write SQL format with the provided table name. The output
	// can be directly ingested by SQLite. Including OptHeader() option will also
	// include a statement to create the table
	OptSql(string) TableOpt

	// OptXml used to write XML format with the provided table id.  Including OptHeader()
	// option will also include the <thead> element at the top of the XML
	// OptXml(string) TableOpt TODO

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

// TableCol represents information about a table column
type TableCol interface {
	// Name returns the column name
	Name() string

	// Type returns the types that the column represents
	Type() Type

	// Min returns the minimum value of all column numbers or zero
	Min() float64

	// Max returns the maximum value of all column numbers or zero
	Max() float64

	// Sum returns the sum of all column numbers or zero
	Sum() float64

	// Count returns the count of all column numbers
	Count() uint64

	// Mean returns the mean average value of all column numbers
	// or +Inf if no numbers in the column
	Mean() float64
}

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	Nil Type = (1 << iota)
	String
	Uint
	Int
	Float
	Duration
	Date
	Datetime
	Bool
	Other
	NumberTypes  = Uint | Int | Float
	DefaultTypes = NumberTypes | Duration | Date | Datetime | Bool
	TypeMin      = Nil
	TypeMax      = Other
)

const (
	BorderDefault = "+++++++++|-"
	BorderLines   = "┌┬┐├┼┤└┴┘│─"
)

/////////////////////////////////////////////////////////////////////
// METHODS

// Is returns true if all types are represented. For example.
//
//   t.Is(Nil) returns true if type contains nil
//   t.Is(Nil|Int) returns true if type is Int or Nil
//
func (t Type) Is(f Type) bool {
	return t&f == f
}

// Type returns a single type which can represent a set
// of types, and also returns true if the columns includes
// nil values
func (t Type) Type() (Type, bool) {
	var isnil = t.Is(Nil)
	if isnil {
		t ^= Nil // Remove nil from t
	}
	// Deal with nil case
	if t == 0 {
		return Nil, isnil
	}
	// Deal with singular case
	if t.Is(Other) {
		return t, isnil
	}
	if t.FlagString() != "other" {
		return t, isnil
	}
	// Numerical cases, float,int,uint
	if t|NumberTypes == NumberTypes {
		switch {
		case t.Is(Float):
			return Float, isnil
		case t.Is(Int):
			return Int, isnil
		default:
			return Uint, isnil
		}
	}
	// Otherwise, string needs to be used to represent
	return String, isnil
}

// Float32 returns a float value from untyped or returns
// NaN otherwise
func Float32(v interface{}) float32 {
	return f32.Cast(v)
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (t Type) String() string {
	if t == 0 {
		return "none"
	}
	str := ""
	for v := TypeMin; v <= TypeMax; v <<= 1 {
		if t&v == v {
			str += "|" + v.FlagString()
		}
	}
	return strings.TrimPrefix(str, "|")
}

func (t Type) FlagString() string {
	switch t {
	case Nil:
		return "nil"
	case Uint:
		return "uint"
	case Int:
		return "int"
	case Float:
		return "float"
	case Duration:
		return "duration"
	case Date:
		return "date"
	case Datetime:
		return "datetime"
	case Bool:
		return "bool"
	case String:
		return "string"
	case Other:
		fallthrough
	default:
		return "other"
	}
}
