package table

import (
	"fmt"
	"strings"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type row struct {
	v []interface{}
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewRow(v []interface{}) *row {
	r := new(row)
	r.v = v
	return r
}

/////////////////////////////////////////////////////////////////////
// METHODS

// row returns an array of a defined size, padding out
// with nils where necessary
func (r *row) row(sz int) []interface{} {
	if len(r.v) == sz {
		return r.v
	} else if sz < len(r.v) {
		return r.v[:sz]
	} else {
		return append(r.v, make([]interface{}, sz-len(r.v))...)
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (r *row) String() string {
	str := "<"
	for i := range r.v {
		str += fmt.Sprint(r.v[i], ",")
	}
	return strings.TrimSuffix(str, ",") + ">"
}
