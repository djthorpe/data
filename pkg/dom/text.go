package dom

import (
	"encoding/xml"
	"fmt"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Text struct {
	Node
	cdata []byte
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewText(cdata []byte, parent *Element, document *Document) *Text {
	text := new(Text)

	// Check parameters
	if document == nil {
		return nil
	}

	// Create node
	text.cdata = make([]byte, len(cdata))
	text.parent = parent
	text.document = document

	// Copy over data
	copy(text.cdata, cdata)

	// Return success
	return text
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Text) Cdata() string {
	return string(this.cdata)
}

func (this *Text) String() string {
	if bytes, err := xml.Marshal(this); err != nil {
		return fmt.Sprintln("Error: ", err)
	} else {
		return string(bytes)
	}
}

/////////////////////////////////////////////////////////////////////
// XML ENCODING

func (this *Text) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	// Cdata
	if err := enc.EncodeElement(this.cdata, start); err != nil {
		return err
	}
	// Return success
	return nil
}
