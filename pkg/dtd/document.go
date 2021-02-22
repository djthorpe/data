package dtd

import (
	"fmt"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Document) Validate(node data.Node) error {
	// Validate children for a node
	return data.ErrInternalAppError
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Document) append(state State) error {
	fmt.Println("Append: ", state)
	return nil
}
