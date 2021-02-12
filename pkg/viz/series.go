package viz

import (
	"errors"
	"fmt"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type series struct {
	sets []data.Set
}

/////////////////////////////////////////////////////////////////////
// CONSTANTS

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

// NewSeries returns an empty array of series
func NewSeries() data.Series {
	this := new(series)
	return this
}

/////////////////////////////////////////////////////////////////////
// METHODS

func (s *series) Read(t data.Table, fn data.SeriesIteratorFunc) error {
	if t == nil || fn == nil {
		return data.ErrBadParameter
	}
	series := []data.Set{}
	for i := 0; i < t.Len(); i++ {
		values, err := fn(i, t.Row(i))
		// Skip values
		if errors.Is(err, data.ErrSkipTransform) {
			continue
		} else if err != nil {
			return err
		} else if len(values) == 0 {
			return data.ErrBadParameter.WithPrefix("Invalid iterator return value")
		}
		// If i==0 then create series and set names
		if i == 0 {
			if series, err = create(s.Len(), values); err != nil {
				return err
			}
			if names, err := fn(-1, nil); len(names) > 0 && err == nil {
				if err := name(series, names); err != nil {
					return err
				}
			} else if errors.Is(err, data.ErrSkipTransform) == false {
				return err
			}
		}
		// Check to make sure there are the right number of values
		if len(values) != len(series) {
			return data.ErrBadParameter.WithPrefix("Invalid iterator return value")
		}
		// Now append values to sets
		for i, s := range series {
			if err := push(s, values[i]); err != nil {
				return err
			}
		}
	}

	// Append series
	s.sets = append(s.sets, series...)

	// Return success
	return nil
}

func (s *series) Len() int {
	return len(s.sets)
}

func (s *series) Sets() []data.Set {
	return s.sets
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (s *series) String() string {
	str := "<series"
	for _, set := range s.sets {
		str += fmt.Sprint(" ", set)
	}
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func create(j int, values []interface{}) ([]data.Set, error) {
	sets := make([]data.Set, len(values))
	for i, v := range values {
		switch v.(type) {
		case data.Point:
			sets[i] = NewPoints(fmt.Sprintf("points_%02d", i+j))
		case string:
			sets[i] = NewLabels(fmt.Sprintf("labels_%02d", i+j))
		case float32:
			sets[i] = NewValues(fmt.Sprintf("values_%02d", i+j))
		default:
			return nil, data.ErrBadParameter.WithPrefix("Invalid iterator return value: ", v)
		}
	}
	return sets, nil
}

func name(sets []data.Set, names []interface{}) error {
	for i, set := range sets {
		if len(names) <= i {
			continue
		}
		if name, ok := names[i].(string); ok == false {
			return data.ErrBadParameter.WithPrefix("Expected string iterator value: ", names[i])
		} else if name != "" {
			set.(Set).SetName(name)
		}
	}
	// Return success
	return nil
}

func push(set data.Set, v interface{}) error {
	switch p := v.(type) {
	case data.Point:
		if set_, ok := set.(data.Points); ok {
			set_.Append(p)
			return nil
		} else {
			return data.ErrBadParameter.WithPrefix("Invalid iterator return value: ", v)
		}
	case string:
		if set_, ok := set.(data.Labels); ok {
			set_.Append(p)
			return nil
		} else {
			return data.ErrBadParameter.WithPrefix("Invalid iterator return value: ", v)
		}
	case float32:
		if set_, ok := set.(data.Values); ok {
			set_.Append(p)
			return nil
		} else {
			return data.ErrBadParameter.WithPrefix("Invalid iterator return value: ", v)
		}
	default:
		return data.ErrBadParameter.WithPrefix("Invalid iterator return value: ", v)
	}
}
