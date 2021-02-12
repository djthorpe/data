package data

/////////////////////////////////////////////////////////////////////
// TYPES

type PointIteratorFunc func(int, []interface{}) (Point, error)

/////////////////////////////////////////////////////////////////////
// INTERFACES

type Points interface {
	// Return name associated with set of points
	Series() string

	// Read points from a table, the iterator function can return
	// ErrSkipTransform if a returned point can be discarded
	Read(Table, PointIteratorFunc) error

	// Write points as path to canvas and return the group
	WritePath(Canvas) CanvasGroup
}
