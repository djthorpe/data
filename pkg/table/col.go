package table

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type col struct {
	i             int
	key           string
	value         string
	types         data.Type
	fmt           string
	min, max, sum float64
	count         uint64
}

var (
	reNotAlphanumeric = regexp.MustCompile("[^a-z0-9]+")
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewCol(i int, v string) *col {
	f := new(col)
	f.i = i
	f.key = keyForValue(i, v)
	f.value = v
	f.types = 0
	f.fmt = ""
	return f
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (c *col) Name() string {
	if c.value == "" {
		return c.key
	} else {
		return c.value
	}
}

func (c *col) Type() data.Type {
	return c.types
}

func (c *col) Min() float64 {
	return c.min
}

func (c *col) Max() float64 {
	return c.max
}

func (c *col) Sum() float64 {
	return c.sum
}

func (c *col) Count() uint64 {
	return c.count
}

func (c *col) Mean() float64 {
	if c.count == 0 {
		return math.Inf(0)
	} else {
		return c.sum / float64(c.count)
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (c *col) String() string {
	str := "<col"
	str += fmt.Sprint(" i=", c.i)
	str += fmt.Sprintf(" name=%q", c.Name())
	str += fmt.Sprintf(" key=%q", c.key)
	str += fmt.Sprintf(" types=%v", c.Type())
	if c.count > 0 {
		str += fmt.Sprintf(" min=%.0f", c.Min())
		str += fmt.Sprintf(" max=%.0f", c.Max())
		str += fmt.Sprintf(" mean=%.0f", c.Mean())
		str += fmt.Sprintf(" count=%v", c.Count())
	}
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func keyForValue(i int, v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	// Where empty column heading, use default
	if v == "" {
		return fmt.Sprintf("col_%02d", i)
	}
	// Convert non-alpha characters into underscores
	return reNotAlphanumeric.ReplaceAllString(v, "_")
}

func (c *col) validate(v interface{}) {
	if v == nil {
		c.types |= data.Nil
		return
	}
	switch v_ := v.(type) {
	case int8, int16, int32, int64, int:
		c.types |= data.Int
	case uint8, uint16, uint32, uint64, uint:
		c.types |= data.Uint
	case float32, float64:
		c.types |= data.Float
	case bool:
		c.types |= data.Bool
	case time.Duration:
		c.types |= data.Duration
	case time.Time:
		if dateValueForTime(v_) {
			c.types |= data.Date
		} else {
			c.types |= data.Datetime
		}
	case string:
		c.types |= data.String
	default:
		c.types |= data.Other
	}
}

func (c *col) asciiWidth() int {
	return 10
}

func (c *col) Fmt() string {
	if c.fmt == "" {
		if t, _ := c.types.Type(); t == data.String {
			c.fmt = fmt.Sprint("%-", c.asciiWidth(), "s")
		} else {
			c.fmt = fmt.Sprint("%", c.asciiWidth(), "s")
		}
	}
	return c.fmt
}

// group will do min,max,mean,sum calculations on uint,int and float values
func (c *col) group(v interface{}) {
	if v == nil {
		// Reset values
		c.min, c.max, c.sum, c.count = 0, 0, 0, 0
		return
	}

	// Calculate sum,min and max for value
	switch v_ := v.(type) {
	case uint8:
		c.groupFloat(float64(v_))
	case uint16:
		c.groupFloat(float64(v_))
	case uint32:
		c.groupFloat(float64(v_))
	case uint64:
		c.groupFloat(float64(v_))
	case uint:
		c.groupFloat(float64(v_))
	case int8:
		c.groupFloat(float64(v_))
	case int16:
		c.groupFloat(float64(v_))
	case int32:
		c.groupFloat(float64(v_))
	case int64:
		c.groupFloat(float64(v_))
	case int:
		c.groupFloat(float64(v_))
	case float32:
		c.groupFloat(float64(v_))
	case float64:
		c.groupFloat(float64(v_))
	}
}

func (c *col) groupFloat(v float64) {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		c.sum = math.NaN()
	} else {
		c.sum += v
	}
	if c.count == 0 {
		c.min = v
		c.max = v
	} else {
		c.min = math.Min(c.min, v)
		c.max = math.Max(c.max, v)
	}
	c.count++
}
