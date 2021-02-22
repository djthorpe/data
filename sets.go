package data

import "time"

/////////////////////////////////////////////////////////////////////
// TYPES

type SeriesIteratorFunc func(int, []interface{}) ([]interface{}, error)

/////////////////////////////////////////////////////////////////////
// INTERFACES

// Series represents one or more ordered sets of the same size
type Series interface {
	// Read series from a table, the iterator function can return
	// ErrSkipTransform if a returned values can be discarded. The
	// values returned should be float32, data.Point or string
	// which respectively appends Values, Points and Labels to sets
	Read(Table, SeriesIteratorFunc) error

	// Sets returns sets contained with the series
	Sets() []Set
}

// Set represents an ordered set of values
type Set interface {
	// Return name associated with the set
	Name() string

	// Len returns the length of the set
	Len() int
}

// PointSet represents an ordered set of points (X,Y)
type PointSet interface {
	Set

	// Append points to the set
	Append(...Point)
}

// LabelSet represents an ordered set of labels
type LabelSet interface {
	Set

	// Append labels to the set
	Append(...string)
}

// RealSet represents an ordered set of float64 values
type RealSet interface {
	Set

	// Append float64 values to the set
	Append(...float64)
}

// TimeSet represents an ordered set of dates or datetime values
type TimeSet interface {
	Set

	// Precision returns the precision of the values
	Precision() time.Duration

	// Append datetime values to the set
	Append(...time.Time)
}
