package table

import (
	"fmt"
	"strings"
	"time"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type row struct {
	s []string
	v []interface{}
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewRow(cap uint) *row {
	r := new(row)
	r.s = make([]string, 0, cap)
	r.v = make([]interface{}, 0, cap)
	return r
}

/////////////////////////////////////////////////////////////////////
// ROW METHODS

func (r *row) set(i uint, s string, v interface{}) {
	if int(i) >= len(r.v) {
		r.s = append(r.s, make([]string, int(i)-len(r.v)+1)...)
		r.v = append(r.v, make([]interface{}, int(i)-len(r.v)+1)...)
	}
	r.s[i] = s
	r.v[i] = v
}

func (r *row) stringsForWidth(width uint) []string {
	result := make([]string, width)
	for i, cell := range r.s {
		if i >= int(width) {
			break
		}
		result[i] = cell
	}
	return result
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (r *row) String() string {
	str := "["
	for _, v := range r.v {
		if v == nil {
			str += "<nil>,"
			continue
		}
		switch v.(type) {
		case float64, float32:
			str += fmt.Sprint("<f>", v, ",")
		case int64, int32, int16, int8, int:
			str += fmt.Sprint("<i>", v, ",")
		case uint64, uint32, uint16, uint8, uint:
			str += fmt.Sprint("<u>", v, ",")
		case bool:
			str += fmt.Sprint("<b>", v, ",")
		case time.Duration:
			str += fmt.Sprint("<d>", v, ",")
		case time.Time:
			str += fmt.Sprint("<dt>", v, ",")
		case string:
			str += fmt.Sprintf("<s>%q,", v)
		default:
			str += fmt.Sprint("<?>", v, ",")
		}
	}
	return strings.TrimSuffix(str, ",") + "]"
}
