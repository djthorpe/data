package table

import (
	"fmt"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type header struct {
	*types
	fields map[string]*field
	order  []*field
}

type field struct {
	name string
	i    int
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewHeader(cap uint) *header {
	h := new(header)
	h.types = NewTypes(cap)
	h.fields = make(map[string]*field, int(cap))
	h.order = make([]*field, 0, int(cap))
	return h
}

func NewField(name string, i int) *field {
	return &field{
		name, i,
	}
}

/////////////////////////////////////////////////////////////////////
// HEADER METHODS

func (h *header) hasField(key string) (int, bool) {
	if f, exists := h.fields[key]; exists {
		return f.i, exists
	} else {
		return -1, exists
	}
}

func (h *header) appendField(field string) (int, error) {
	l := len(h.order)
	f := NewField(field, l)
	k := f.Key()
	if i, exists := h.hasField(k); exists {
		return i, data.ErrDuplicateEntry
	} else {
		h.fields[k] = f
		h.order = append(h.order, f)
	}

	// Return success
	return l, nil
}

func (h *header) setWidth(width uint) {
	// If width is over, remove fields
	l := len(h.order)
	if l > int(width) {
		for i := int(width); i < l; i++ {
			k := h.order[i].Key()
			delete(h.fields, k)
		}
		h.order = h.order[:width]
	}
	// If width is identical, return
	if len(h.order) == int(width) {
		return
	}
	// We need to add rows on
	for len(h.order) < int(width) {
		k := fmt.Sprintf("Col_%02d", len(h.order))
		h.appendField(k)
	}
}

func (h *header) stringsForWidth(width uint) []string {
	result := make([]string, width)
	for i, f := range h.order {
		if i >= int(width) {
			break
		}
		result[i] = f.name
	}
	return result

}

/////////////////////////////////////////////////////////////////////
// FIELD METHODS

func (f *field) Key() string {
	return f.name
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (h *header) String() string {
	str := "<h"
	str += fmt.Sprint(" ", h.order)
	return str + ">"
}

func (f *field) String() string {
	str := "<field"
	str += fmt.Sprintf(" name=%q key=%q", f.name, f.Key())
	return str + ">"
}
