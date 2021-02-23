package set_test

import (
	"testing"

	"github.com/djthorpe/data/pkg/set"
)

func Test_RealSet_001(t *testing.T) {
	s := set.NewRealSet("realset")
	if s == nil {
		t.Fatal("Expected non-nil from NewRealSet")
	}
	s.Append(0.0)
	if s.Min() != 0.0 {
		t.Fatal("Unexpected return from Min", s.Min())
	}
	if s.Max() != 0.0 {
		t.Fatal("Unexpected return from Max", s.Max())
	}
	if s.Len() != 1 {
		t.Fatal("Unexpected return from Len", s.Len())
	}
}
