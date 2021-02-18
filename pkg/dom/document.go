package dom

import (
	"encoding/xml"
	"fmt"
	"io"
	"regexp"
	"strconv"
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
	tag  map[string]string
}

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DefaultDOMOptions = data.DOMWriteDirective | data.DOMWriteIndentSpace2
)

var (
	reTagNSName = regexp.MustCompile("^[A-Za-z][A-Za-z0-9]*$")
	reAttrName  = regexp.MustCompile("^[A-Za-z][A-Za-z0-9_-]*$")
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
	doc.tag = make(map[string]string)

	// Create root element
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

func (this *Document) setTagNS(ns string) error {
	// Get prefix and namespace
	prefix, ns, err := parseTagNS(ns)
	if err != nil {
		return err
	}
	otherprefix, exists := this.tag[ns]
	if exists == false {
		if this.hasTagNS(prefix) {
			this.tag[ns] = this.newTagNS(ns)
		} else {
			this.tag[ns] = prefix
		}
		return nil
	} else if prefix == "" {
		return nil
	} else if otherprefix != prefix {
		return data.ErrBadParameter.WithPrefix("setTagNS: ", prefix)
	} else {
		// Probably need to fix this later!
		return data.ErrInternalAppError.WithPrefix("setTagNS: Unhandled condition")
	}
}

func (this *Document) newTagNS(ns string) string {
	var prefix = "ns"
	if tag, exists := xmlNs[ns]; exists {
		prefix = tag
	}
	if this.hasTagNS(prefix) == false {
		return prefix
	}
	for i := 1; ; i++ {
		prefix_ := fmt.Sprint(prefix, i)
		if this.hasTagNS(prefix_) == false {
			return prefix_
		}
	}
}

func (this *Document) hasTagNS(tag string) bool {
	for _, v := range this.tag {
		if tag == v {
			return true
		}
	}
	return false
}

func (this *Document) getTagNS(element *Element, name xml.Name) (xml.Name, error) {
	// Return unconverted name if no namespace
	if name.Space == "" {
		return name, nil
	}
	if tag, ns, err := parseTagNS(name.Space); err != nil {
		return name, err
	} else if othertag, exists := this.tag[ns]; exists == false {
		return name, data.ErrNotFound.WithPrefix("GetTagNS: ", strconv.Quote(name.Space))
	} else if tag != "" && othertag != tag {
		return name, data.ErrBadParameter.WithPrefix("GetTagNS: ", strconv.Quote(tag))
	} else {
		// Add prefix to the tag name
		if othertag != "" {
			name.Local = othertag + ":" + name.Local
		}
		// If not root tag, remove the name.Space
		if element.IsRootElement() == false {
			name.Space = ""
		}
	}

	// Return success
	return name, nil
}

func parseTagNS(ns string) (string, string, error) {
	if nsTag := strings.SplitN(ns, " ", 2); len(nsTag) == 1 {
		return "", nsTag[0], nil
	} else if reTagNSName.MatchString(nsTag[1]) == false {
		return "", "", data.ErrBadParameter.WithPrefix("ParseTagNS: ", strconv.Quote(ns))
	} else {
		return nsTag[1], nsTag[0], nil
	}
}
