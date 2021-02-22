package dtd

import (
	"encoding/xml"
	"fmt"
	"io"
)

/////////////////////////////////////////////////////////////////////
// XML DECODING

func (this *Document) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {

	fmt.Println("start: ", start)

	// Read tokens until end
	for {
		t, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		fmt.Println("  token: ", t)
	}

	// Return success
	return nil
}
