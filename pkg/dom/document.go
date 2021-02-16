package dom

import (
	"encoding/xml"
	"io"
	"strings"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Document struct {
	*Element

	opts data.DOMOption
	id   map[string]map[string]*Element // ns->value->node
	fn   data.DOMValidateNodeFunc
}

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DefaultDOMOptions = data.DOMWriteDirective | data.DOMWriteIndentSpace2
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewDocument(name string) data.Document {
	return NewDocumentNS(name, "")
}

func NewDocumentNS(name, ns string) data.Document {
	doc := new(Document)

	doc.id = make(map[string]map[string]*Element)
	doc.opts = DefaultDOMOptions
	doc.Element = NewElementNS(name, ns, nil, doc)

	// Return success
	return doc
}

/////////////////////////////////////////////////////////////////////
// METHODS

func (this *Document) CreateElement(name string) data.Node {
	return this.CreateElementNS(name, "")
}

func (this *Document) CreateElementNS(name, ns string) data.Node {
	return NewElementNS(name, ns, nil, this)
}

func (this *Document) CreateText(cdata string) data.Node {
	return NewText([]byte(cdata), nil, this)
}

func (this *Document) CreateComment(cdata string) data.Node {
	return NewComment([]byte(cdata), nil, this)
}

func (this *Document) GetElementById(value string) data.Node {
	return this.GetElementByIdNS(value, "")
}
func (this *Document) GetElementByIdNS(value, ns string) data.Node {
	return this.getAttrId(value, ns)
}

/////////////////////////////////////////////////////////////////////
// WRITE

func (this *Document) Write(w io.Writer) error {
	return this.WriteEx(w, this.opts)
}

func (this *Document) WriteEx(w io.Writer, opts data.DOMOption) error {
	enc := xml.NewEncoder(w)
	if opts.Is(data.DOMWriteIndentSpace2) {
		enc.Indent("", "  ")
	} else if opts.Is(data.DOMWriteIndentTab) {
		enc.Indent("", "\t")
	}
	if opts.Is(data.DOMWriteDirective) {
		if _, err := w.Write([]byte(xml.Header)); err != nil {
			return err
		}
	}
	if err := enc.Encode(this.Element); err != nil {
		return err
	}
	if err := enc.Flush(); err != nil {
		return err
	}
	return nil
}

/////////////////////////////////////////////////////////////////////
// READ

func Read(r io.Reader) (data.Document, error) {
	return ReadEx(r, nil)
}

func ReadEx(r io.Reader, fn data.DOMValidateNodeFunc) (data.Document, error) {
	this := NewDocumentNS("xml", "").(*Document)
	this.fn = fn
	dec := xml.NewDecoder(r)

	// Decode document
	if err := dec.Decode(this); err != nil {
		return nil, err
	}

	// Remove validator
	this.fn = nil

	// Return success
	return this, nil
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Document) String() string {
	w := new(strings.Builder)
	if err := this.WriteEx(w, 0); err != nil {
		panic(err)
	}
	return w.String()
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// setAttr indexes elements based on id attribute
func (this *Document) setAttr(name, ns string, value string, element *Element) {
	switch name {
	case "id":
		this.setAttrId(value, ns, element)
	}
}

func (this *Document) setAttrId(value, ns string, element *Element) {
	value = strings.TrimSpace(value)
	if _, exists := this.id[ns]; exists == false {
		this.id[ns] = make(map[string]*Element)
	}
	if element == nil {
		delete(this.id[ns], value)
	} else {
		this.id[ns][value] = element
	}
}

func (this *Document) getAttrId(value, ns string) *Element {
	value = strings.TrimSpace(value)
	if _, exists := this.id[ns]; exists == false {
		return nil
	} else if element, exists := this.id[ns][value]; exists == false {
		return nil
	} else {
		return element
	}
}
