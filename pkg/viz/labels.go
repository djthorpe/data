package viz

import (
	"fmt"
	"strings"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type labels struct {
	name string
	v    []string
}

/////////////////////////////////////////////////////////////////////
// CONSTANTS

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

// NewValues returns an empty array of values
func NewLabels(name string) data.Labels {
	this := new(labels)
	this.name = name
	return this
}

/////////////////////////////////////////////////////////////////////
// METHODS

func (p *labels) Name() string {
	return p.name
}

func (p *labels) SetName(name string) {
	p.name = name
}

func (p *labels) Len() int {
	return len(p.v)
}

func (p *labels) Append(label string) {
	p.v = append(p.v, label)
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (p *labels) String() string {
	str := "<labels"
	str += fmt.Sprintf(" name=%q", p.Name())
	if len(p.v) > 0 {
		str += " <"
		for _, v := range p.v {
			str += fmt.Sprintf("%q,", v)
		}
		str = strings.TrimSuffix(str, ",") + ">"
	}
	return str + ">"
}
