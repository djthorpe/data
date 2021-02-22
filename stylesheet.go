package data

import (
	"io"
)

/////////////////////////////////////////////////////////////////////
// INTERFACES

type StyleSheet interface {
	// Read styles from a data stream
	Read(r io.Reader) error

	// Write styles to a data stream
	Write(r io.Writer) error

	// Return style rules by tag, class, id and state. Where any element is
	// empty any rule can match
	Rules(tag, class, id, state string) []StyleRule

	// Add a style rule by combination of tag, class, id and state
	// followed by one or more name/value pairs "<name>:<value>; <name>:<value>"
	Add(string, string) error
}

type StyleRule interface {
	// Return tag, class, id and state for the rule
	Tag() string
	Class() string
	Id() string
	State() string

	// Return the style name and value
	Name() string
	Value() interface{}
}
