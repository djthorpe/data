package viz

import (
	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Set interface {
	data.Set

	SetName(string)
}
