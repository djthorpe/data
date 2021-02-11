package table

import (
	"fmt"
	"strings"
	"time"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type col struct {
	i     int
	key   string
	value string
	types map[data.Type]bool
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewCol(i int, value string) *col {
	f := new(col)
	f.i = i
	f.key = strings.TrimSpace(value)
	f.value = value
	f.types = make(map[data.Type]bool, 5)
	return f
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (c *col) Name() string {
	return c.value
}

func (c *col) Type() data.Type {
	return data.String
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (c *col) String() string {
	str := "<col"
	str += fmt.Sprint(" i=", c.i)
	str += fmt.Sprintf(" name=%q", c.value)
	str += fmt.Sprintf(" key=%q", c.key)
	str += fmt.Sprintf(" types=%v", c.types)
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (c *col) validate(v interface{}) {
	if v == nil {
		c.types[data.Nil] = true
		return
	}
	switch v_ := v.(type) {
	case int8, int16, int32, int64, int:
		c.types[data.Int] = true
	case uint8, uint16, uint32, uint64, uint:
		c.types[data.Uint] = true
	case float32, float64:
		c.types[data.Float] = true
	case bool:
		c.types[data.Bool] = true
	case time.Duration:
		c.types[data.Duration] = true
	case time.Time:
		if v_.Hour() == 0 && v_.Minute() == 0 && v_.Second() == 0 && v_.Nanosecond() == 0 {
			c.types[data.Date] = true
		} else {
			c.types[data.Datetime] = true
		}
	}
}
