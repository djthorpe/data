package set_test

import (
	"testing"
	"time"

	"github.com/djthorpe/data/pkg/set"
)

func Test_TimeSet_001(t *testing.T) {
	s := set.NewTimeSet("timeset")
	if s == nil {
		t.Fatal("Expected non-nil from NewTimeSet")
	}
	now := time.Now()
	s.Append(now)
	if s.Min() != now {
		t.Fatal("Unexpected return from Min", s.Min())
	}
	if s.Max() != now {
		t.Fatal("Unexpected return from Max", s.Max())
	}
	if s.Len() != 1 {
		t.Fatal("Unexpected return from Len", s.Len())
	}
	if s.Duration() != 0 {
		t.Fatal("Unexpected return from Duration", s.Duration())
	}
}

func Test_TimeSet_002(t *testing.T) {
	s := set.NewTimeSet("timeset")
	if s == nil {
		t.Fatal("Expected non-nil from NewTimeSet")
	}
	now := time.Now()

	s.Append(now.Truncate(time.Hour * 24))
	if p := s.Precision(); p != time.Hour*24 {
		t.Error("Unexpected return from Precision", s.Precision())
	}

	s.Append(now.Truncate(time.Hour))
	if p := s.Precision(); p != time.Hour {
		t.Error("Unexpected return from Precision", s.Precision())
	}

	t.Log(s)
}
