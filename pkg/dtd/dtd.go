package dtd

import (
	"io"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Document struct {
	name string
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewDocument(name string) *Document {
	this := new(Document)
	this.name = name
	return this
}

func Read(r io.Reader) (*Document, error) {
	this := NewDocument("dtd")
	dec := NewDecoder(r)

	// Decode document
	if err := dec.Decode(this); err != nil {
		return nil, err
	}

	// Return success
	return this, nil
}
