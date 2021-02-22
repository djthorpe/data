package dtd

import "fmt"

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Document) append(state State) error {
	fmt.Println("Append: ", state)
	return nil
}
