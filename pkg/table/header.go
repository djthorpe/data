package table

import (
	"fmt"
	"strings"
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
	h.w = cap
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

func (h *header) row() []string {
	result := make([]string, h.w)
	for i := range result {
		if f, exists := h.i[i]; exists {
			result[i] = f.value
		} else {
			result[i] = fmt.Sprintf("Col_%02d", i)
		}
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

func (h *header) validate(row []interface{}) {
	for i, v := range row {
		if c, exists := h.i[i]; exists {
			c.validate(v)
		}
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (h *header) String() string {
	str := "<h"
	str += " " + strings.Join(h.row(), ",")
	return str + ">"
}
