package table

import (
	"time"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type optFlag uint64

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	optHeader optFlag = (1 << iota)
	optNil
	optUint
	optInt
	optFloat
	optBool
	optDuration
	optDate
	optDatetime
	optAscii
	optCsv
)

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// OptHeader on Read options ensures the first row read is the header
// for the CSV, on Write ensures the header is output before the first
// row
func (t *Table) OptHeader() data.TableOpt {
	return func(t data.Table) {
		t.(*Table).setOpt(optHeader, true)
	}
}

func (t *Table) OptTransform(fns ...data.TransformFunc) data.TableOpt {
	return func(t data.Table) {
		t.(*Table).setTransform(fns)
	}
}

func (t *Table) OptType(types data.Type) data.TableOpt {
	return func(t data.Table) {
		if types&data.Nil != 0 {
			t.(*Table).setOpt(optNil, true)
		}
		if types&data.Uint != 0 {
			t.(*Table).setOpt(optUint, true)
		}
		if types&data.Int != 0 {
			t.(*Table).setOpt(optInt, true)
		}
		if types&data.Float != 0 {
			t.(*Table).setOpt(optFloat, true)
		}
		if types&data.Duration != 0 {
			t.(*Table).setOpt(optDuration, true)
		}
		if types&data.Date != 0 {
			t.(*Table).setOpt(optDate, true)
		}
		if types&data.Datetime != 0 {
			t.(*Table).setOpt(optDatetime, true)
		}
		if types&data.Bool != 0 {
			t.(*Table).setOpt(optBool, true)
		}
	}
}

func (t *Table) OptAscii(width uint, border string) data.TableOpt {
	return func(t data.Table) {
		t.(*Table).setOpt(optAscii, true)
		if border == "" {
			t.(*Table).opts.border = []rune(data.BorderDefault)
		} else {
			t.(*Table).opts.border = []rune(border)
		}
		t.(*Table).opts.asciiwidth = int(width)
	}
}

func (t *Table) OptCsv(delim rune) data.TableOpt {
	return func(t data.Table) {
		t.(*Table).setOpt(optCsv, true)
		if delim != 0 {
			t.(*Table).opts.csvDelim = delim
		}
	}
}

func (t *Table) OptDuration(dur time.Duration) data.TableOpt {
	return func(t data.Table) {
		t.(*Table).setOpt(optDuration, true)
		t.(*Table).opts.dur = dur
	}
}

func (t *Table) OptTimezone(tz *time.Location) data.TableOpt {
	return func(t data.Table) {
		t.(*Table).setOpt(optDate|optDatetime, true)
		t.(*Table).opts.tz = tz
	}
}

func (t *Table) OptRowIterator(fn data.IteratorFunc) data.TableOpt {
	return func(t data.Table) {
		t.(*Table).opts.iterator = fn
	}
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (t *Table) applyOpt(opts []data.TableOpt) {
	// Set default options
	t.opts.o = t.opts.d
	t.opts.tz = time.Local
	t.opts.dur = time.Second
	t.opts.border = []rune(data.BorderDefault)
	t.opts.transform = []data.TransformFunc{}
	t.opts.iterator = nil

	// Apply options
	for _, opt := range opts {
		opt(t)
	}
}

func (t *Table) hasOpt(f optFlag) bool {
	return t.opts.o.has(f)
}

func (t *Table) setOpt(f optFlag, v bool) {
	if v {
		t.opts.o |= f
	} else {
		t.opts.o ^= f
	}
}

func (t *Table) setTransform(fns []data.TransformFunc) {
	t.opts.transform = fns
}

func (f optFlag) has(v optFlag) bool {
	return f&v == v
}
