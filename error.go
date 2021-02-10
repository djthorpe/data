package data

import (
	"fmt"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Error uint

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	ErrNone Error = iota
	ErrBadParameter
	ErrDuplicateEntry
	ErrInternalAppError
)

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (e Error) WithPrefix(args ...interface{}) error {
	return fmt.Errorf("%s: %w", fmt.Sprint(args...), e)
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (e Error) Error() string {
	switch e {
	case ErrNone:
		return "ErrNone"
	case ErrBadParameter:
		return "ErrBadParameter"
	case ErrDuplicateEntry:
		return "ErrDuplicateEntry"
	case ErrInternalAppError:
		return "ErrInternalAppError"
	default:
		return "[?? Invalid Error value]"
	}
}
