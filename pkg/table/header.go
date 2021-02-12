package table

import (
	"fmt"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type header struct {
	w int
	f map[string]*col
	i map[int]*col
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewHeader(cap int) *header {
	h := new(header)
	h.w = 0
	h.f = make(map[string]*col, cap)
	h.i = make(map[int]*col, cap)
	return h
}

/////////////////////////////////////////////////////////////////////
// HEADER METHODS

func (h *header) set(row []string) []int {
	order := make([]int, len(row))
	for i, value := range row {
		order[i] = h.append(value)
	}
	return order
}

func (h *header) append(value string) int {
	f := NewCol(h.w, value)
	if f_, exists := h.f[f.key]; exists == false {
		h.f[f.key], h.i[f.i] = f, f
		h.w++
		return f.i
	} else {
		return f_.i
	}
}

func (h *header) cols() []*col {
	result := make([]*col, h.w)
	for i := range result {
		if _, exists := h.i[i]; exists == false {
			f := NewCol(i, "")
			h.f[f.key], h.i[f.i] = f, f
			h.w++
		}
		result[i] = h.i[i]
	}
	return result
}

func (h *header) col(i int) *col {
	if c, exists := h.i[i]; exists {
		return c
	} else {
		return nil
	}
}

// validate will adjust width of header to accommodate
// new values and then validate each value against header
// returns true if width was changed
func (h *header) validate(row *row) bool {
	var ch bool
	// Add additional columns
	if n := len(row.v) - h.w; n > 0 {
		ch = true
		for i := 0; i < n; i++ {
			h.append("")
		}
	}
	// Validate values
	for i, v := range row.row(h.w) {
		h.i[i].validate(v)
	}
	// Return true if dimensions of table have changed
	return ch
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (h *header) String() string {
	str := "<h"
	str += " " + fmt.Sprint(h.cols())
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func maxInt(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
