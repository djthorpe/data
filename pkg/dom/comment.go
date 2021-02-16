package dom

import (
	"encoding/xml"
	"fmt"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Comment struct {
	Node
	cdata []byte
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewComment(cdata []byte, parent *Element, document *Document) *Comment {
	comment := new(Comment)

	// Check parameters
	if document == nil {
		return nil
	}

	// Create node
	comment.cdata = make([]byte, len(cdata))
	comment.parent = parent
	comment.document = document

	// Copy over data
	copy(comment.cdata, cdata)

	// Return success
	return comment
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Comment) Cdata() string {
	return string(this.cdata)
}

func (this *Comment) String() string {
	if bytes, err := xml.Marshal(this); err != nil {
		return fmt.Sprintln("Error: ", err)
	} else {
		return "<!-- " + string(bytes) + " -->"
	}
}

/////////////////////////////////////////////////////////////////////
// XML ENCODING

func (this *Comment) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	// Cdata
	if err := enc.EncodeElement(this.cdata, start); err != nil {
		return err
	}
	// Return success
	return nil
}
