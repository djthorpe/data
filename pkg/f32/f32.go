package f32

import (
	"fmt"
	"math"
	"strings"
)

// Return absolute value
func Abs(v float32) float32 {
	if IsNaN(v) || v >= 0 {
		return v
	} else {
		return -v
	}
}

// Return minimum value
func Min(v ...float32) float32 {
	if len(v) == 0 {
		return NaN()
	}
	if len(v) == 1 {
		return v[0]
	}
	min := v[0]
	for _, v := range v {
		if min > v {
			min = v
		}
	}
	return min
}

// Return maximum value
func Max(v ...float32) float32 {
	if len(v) == 0 {
		return NaN()
	}
	if len(v) == 1 {
		return v[0]
	}
	max := v[0]
	for _, v := range v {
		if max < v {
			max = v
		}
	}
	return max
}

// Cast any numeric value to float32 or return NaN
func Cast(v interface{}) float32 {
	if v == nil {
		return NaN()
	}
	switch v_ := v.(type) {
	case float32:
		return v_
	case float64:
		return float32(v_)
	case uint:
		return float32(v_)
	case uint8:
		return float32(v_)
	case uint16:
		return float32(v_)
	case uint32:
		return float32(v_)
	case uint64:
		return float32(v_)
	case int:
		return float32(v_)
	case int8:
		return float32(v_)
	case int16:
		return float32(v_)
	case int32:
		return float32(v_)
	case int64:
		return float32(v_)
	default:
		return NaN()
	}
}

// Return true if value is NaN
func IsNaN(f float32) bool {
	return f != f
}

// Return NaN value
func NaN() float32 {
	return float32(math.NaN())
}

// Return string with values delimited by commas
func String(values ...float32) string {
	str := ""

	// Deal with 0 and 1 cases
	if len(values) == 0 {
		return str
	} else if len(values) == 1 {
		return string1(values[0])
	}

	// Concatenate values
	for _, value := range values {
		str += string1(value) + ","
	}
	return strings.TrimSuffix(str, ",")
}

// Return integer form or decimal form
func string1(value float32) string {
	if float32(int64(value)) == value {
		return fmt.Sprintf("%.0f", value)
	} else {
		return fmt.Sprintf("%f", value)
	}
}
