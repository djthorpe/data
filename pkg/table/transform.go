package table

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type transformInFunc func(*Table, string) (interface{}, error)
type transformOutFunc func(*Table, interface{}) (string, error)

/////////////////////////////////////////////////////////////////////
// CONSTANTS

var (
	defaultInTransforms = []transformInFunc{
		transformInNil,
		transformInUint,
		transformInInt,
		transformInFloat,
		transformInDuration,
		transformInBool,
		transformInDate,
		transformInDatetime,
	}
	defaultOutTransforms = []transformOutFunc{
		transformOutNil,
		transformOutDuration,
		transformOutDate,
	}
)

/////////////////////////////////////////////////////////////////////
// VALUE TRANSFORMS

// inValue converts string to native type
func (t *Table) inValue(i, j int, value interface{}) (interface{}, error) {
	if str, ok := value.(string); ok == false {
		return nil, data.ErrInternalAppError.WithPrefix("inValue")
	} else {
		for _, fn := range t.opts.transform {
			if fn == nil {
				continue
			} else if v_, err := fn(i, j, str); err == data.ErrSkipTransform {
				continue
			} else if err != nil {
				return nil, err
			} else {
				return v_, nil
			}
		}
		// Use default transformation
		return t.defaultInTransform(str)
	}
}

// outValue converts from native type to string
func (t *Table) outValue(i, j int, value interface{}) (interface{}, error) {
	for _, fn := range t.opts.transform {
		if fn == nil {
			continue
		} else if v_, err := fn(i, j, value); err == data.ErrSkipTransform {
			continue
		} else if err != nil {
			return nil, err
		} else {
			return v_, nil
		}
	}
	// Use default transformation
	return t.defaultOutTransform(value)
}

// rowIterator calls a row iterator
func (t *Table) rowIterator(i int, row []interface{}) error {
	if t.opts.iterator != nil {
		if err := t.opts.iterator(i, row); err != nil {
			return err
		}
	}
	// Return success
	return nil
}

func (t *Table) defaultInTransform(str string) (interface{}, error) {
	for _, fn := range defaultInTransforms {
		if value, err := fn(t, str); err == nil {
			return value, nil
		} else if errors.Is(err, data.ErrSkipTransform) == false {
			return nil, err
		}
	}
	// By default, no transformation is done
	return str, nil
}

func (t *Table) defaultOutTransform(value interface{}) (string, error) {
	for _, fn := range defaultOutTransforms {
		if value, err := fn(t, value); err == nil {
			return value, nil
		} else if errors.Is(err, data.ErrSkipTransform) == false {
			return "", err
		}
	}
	// By default, no transformation is done
	return fmt.Sprint(value), nil
}

func transformInNil(t *Table, str string) (interface{}, error) {
	// Check for nil values
	if t.hasOpt(optNil) && nilValueForString(str) {
		return nil, nil
	} else {
		return nil, data.ErrSkipTransform
	}
}

func transformOutNil(t *Table, value interface{}) (string, error) {
	if value != nil {
		return "", data.ErrSkipTransform
	}
	if t.hasOpt(optNil) {
		return "<nil>", nil
	} else {
		return "", nil
	}
}

func transformInUint(t *Table, str string) (interface{}, error) {
	if t.hasOpt(optUint) == false {
		return nil, data.ErrSkipTransform
	} else if v, err := strconv.ParseUint(str, 0, 64); err == nil {
		return v, nil
	} else {
		return nil, data.ErrSkipTransform
	}
}

func transformInInt(t *Table, str string) (interface{}, error) {
	if t.hasOpt(optInt) == false {
		return nil, data.ErrSkipTransform
	} else if v, err := strconv.ParseInt(str, 0, 64); err == nil {
		return v, nil
	} else {
		return nil, data.ErrSkipTransform
	}
}

func transformInFloat(t *Table, str string) (interface{}, error) {
	if t.hasOpt(optFloat) == false {
		return nil, data.ErrSkipTransform
	} else if v, err := strconv.ParseFloat(str, 64); err == nil {
		return v, nil
	} else {
		return nil, data.ErrSkipTransform
	}
}

func transformInDuration(t *Table, str string) (interface{}, error) {
	if t.hasOpt(optDuration) == false {
		return nil, data.ErrSkipTransform
	} else if v, err := time.ParseDuration(str); err == nil {
		if t.opts.dur != 0 {
			return v.Truncate(t.opts.dur), nil
		} else {
			return v, nil
		}
	} else if v, err := strconv.ParseFloat(str, 64); err == nil {
		if t.opts.dur == 0 {
			return time.Duration(v * float64(time.Second)), nil
		} else {
			return time.Duration(v * float64(t.opts.dur)), nil
		}
	} else {
		return nil, data.ErrSkipTransform
	}
}

func transformOutDuration(t *Table, v interface{}) (string, error) {
	if t.hasOpt(optDuration) == false {
		return "", data.ErrSkipTransform
	}
	if v_, ok := v.(time.Duration); ok == false {
		return "", data.ErrSkipTransform
	} else if t.opts.dur == 0 {
		return v_.String(), nil
	} else {
		return v_.Truncate(t.opts.dur).String(), nil
	}
}

func transformInBool(t *Table, str string) (interface{}, error) {
	if t.hasOpt(optBool) == false {
		return nil, data.ErrSkipTransform
	} else if v, err := strconv.ParseBool(str); err == nil {
		return v, nil
	} else {
		return nil, data.ErrSkipTransform
	}
}

func transformInDate(t *Table, str string) (interface{}, error) {
	if t.hasOpt(optDate) == false {
		return nil, data.ErrSkipTransform
	} else if v, err := time.ParseInLocation("2006-1-2", str, t.opts.tz); err == nil {
		return v, nil
	} else if v, err := time.ParseInLocation("2006/1/2", str, t.opts.tz); err == nil {
		return v, nil
	} else if v, err := time.ParseInLocation("2/1/2006", str, t.opts.tz); err == nil {
		return v, nil
	} else if v, err := time.ParseInLocation("2-1-2006", str, t.opts.tz); err == nil {
		return v, nil
	} else if v, err := time.ParseInLocation("Jan 2 2006", str, t.opts.tz); err == nil {
		return v, nil
	} else if v, err := time.ParseInLocation("2 Jan 2006", str, t.opts.tz); err == nil {
		return v, nil
	} else if v, err := time.ParseInLocation("Jan 2 06", str, t.opts.tz); err == nil {
		return v, nil
	} else if v, err := time.ParseInLocation("2 Jan 06", str, t.opts.tz); err == nil {
		return v, nil
	} else {
		return nil, data.ErrSkipTransform
	}
}

func transformOutDate(t *Table, v interface{}) (string, error) {
	if t.hasOpt(optDate) == false {
		return "", data.ErrSkipTransform
	}
	if v_, ok := v.(time.Time); ok == false {
		return "", data.ErrSkipTransform
	} else if dateValueForTime(v_) == false {
		return "", data.ErrSkipTransform
	} else {
		return v_.Format("2006-01-02"), nil
	}
}

func transformInDatetime(t *Table, str string) (interface{}, error) {
	if t.hasOpt(optDatetime) == false {
		return nil, data.ErrSkipTransform
	} else if v, err := time.ParseInLocation(time.RFC3339, str, t.opts.tz); err == nil {
		return v, nil
	} else if v, err := time.ParseInLocation(time.UnixDate, str, t.opts.tz); err == nil {
		return v, nil
	} else if v, err := time.ParseInLocation(time.RFC822, str, t.opts.tz); err == nil {
		return v, nil
	} else if v, err := time.ParseInLocation("2006-01-02 15:04:05", str, t.opts.tz); err == nil {
		return v, nil
	} else if v, err := time.ParseInLocation("2006-01-02 15:04", str, t.opts.tz); err == nil {
		return v, nil
	} else {
		return nil, data.ErrSkipTransform
	}
}

func nilValueForString(v string) bool {
	return v == "" || strings.TrimSpace(v) == ""
}

func dateValueForTime(v time.Time) bool {
	return v.Hour() == 0 && v.Minute() == 0 && v.Second() == 0 && v.Nanosecond() == 0
}
