package color_test

import (
	"math/rand"
	"testing"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/color"
)

func Test_Color_001(t *testing.T) {
	// Return all colors
	all := color.Palette(data.ColorAll)
	for _, c := range all {
		if name := color.Name(c); name == "" {
			t.Error("Missing name for color", c)
		} else {
			t.Logf("color=%v name=%q", c, name)
		}
	}
}

func Test_Color_002(t *testing.T) {
	// Distance to self
	all := color.Palette(data.ColorAll)
	for _, c := range all {
		if distance := color.Distance(c, c); distance != 0.0 {
			t.Errorf("Non-zero distance for color %q: %v", color.Name(c), distance)
		}
	}
}

func Test_Color_003(t *testing.T) {
	// Hash for color
	all := color.Palette(data.ColorAll)
	for _, c := range all {
		if hash := color.HashString(c); hash == "" {
			t.Error("Unexpected return from Hash function")
		} else {
			t.Logf("color=%v hash=%q", c, hash)
		}
	}
}

func Test_Color_004(t *testing.T) {
	// Nearest named color for random colors
	for i := 0; i < 100; i++ {
		random := data.Color{uint8(rand.Uint32()), uint8(rand.Uint32()), uint8(rand.Uint32())}
		c := color.Nearest(random, nil)
		if name := color.Name(c); name == "" {
			t.Error("Unexpected return from Hash function")
		} else {
			t.Logf("random=%v color=%v name=%q", random, c, name)
		}
	}
}
