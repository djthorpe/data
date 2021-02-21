package data

import "io"

/////////////////////////////////////////////////////////////////////
// TYPES

//type SeriesIteratorFunc func(int, []interface{}) ([]interface{}, error)

/////////////////////////////////////////////////////////////////////
// INTERFACES

type Viz interface {
	// Create graph paper with major x minor grid squares
	GraphPaper(major, minor uint) VizGraphPaper

	// Write to data stream
	Write(Writer, io.Writer) error
}

type VizGraphPaper interface{}

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// Class names used for GraphPaper
	ClassGraphPaper       = "graphpaper"
	ClassGraphPaperBorder = "border"
	ClassGraphPaperXMajor = "majorx"
	ClassGraphPaperYMajor = "majory"
	ClassGraphPaperXMinor = "minorx"
	ClassGraphPaperYMinor = "minory"
)

/*
type VizGroup interface {
	CanvasGroup
}

type Series interface {
	// Read series from a table, the iterator function can return
	// ErrSkipTransform if a returned values can be discarded. The
	// values returned should be float32, data.Point or string
	// which respectively appends Values, Points and Labels
	Read(Table, SeriesIteratorFunc) error

	// Sets returns sets contained with the series
	Sets() []Set
}

type Set interface {
	// Return name associated with the set
	Name() string

	// Len returns the length of the set
	Len() int
}

type Points interface {
	Set

	// Append point to the set
	Append(Point)
}

type Values interface {
	Set

	// Append value to the set
	Append(float32)

	// Calculate a scale which can represent all values
	Scale() Scale
}

type Labels interface {
	Set

	// Append label to the set
	Append(string)
}

// Scale is an X or Y scale (currently linear) which
// can represent all values
type Scale interface {
	// Return name associated with the scale
	Name() string

	// Return minimum represented value on scale
	Min() float32

	// Return maximum represented value on scale
	Max() float32

	// Write scale to canvas
	WritePath(Canvas) CanvasGroup
}
*/
