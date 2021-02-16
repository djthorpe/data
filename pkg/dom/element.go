package dom

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Element struct {
	Node
	XMLName  xml.Name
	attrs    map[string]*xml.Attr
	attro    []string
	children []interface{}
}

/////////////////////////////////////////////////////////////////////
// CONSTANTS

var (
	reAttrName = regexp.MustCompile("^[A-Za-z][A-Za-z0-9_-]*$")
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewElementNS(name, ns string, parent *Element, document *Document) *Element {
	element := new(Element)

	// Santitize parameters
	name = strings.TrimSpace(name)
	ns = strings.TrimSpace(ns)

	// Check parameters
	if document == nil || name == "" {
		return nil
	}

	// Create node
	element.XMLName = xml.Name{
		Local: name,
		Space: ns,
	}
	element.parent = parent
	element.document = document
	element.attrs = make(map[string]*xml.Attr)

	// Return success
	return element
}

/////////////////////////////////////////////////////////////////////
// ELEMENT METHODS

func (this *Element) Name() xml.Name {
	return this.XMLName
}

func (this *Element) Attrs() []xml.Attr {
	attrs := make([]xml.Attr, len(this.attro))
	for i, key := range this.attro {
		attrs[i] = *(this.attrs[key])
	}
	return attrs
}

func (this *Element) Attr(name string) (xml.Attr, bool) {
	return this.AttrNS(name, "")
}

func (this *Element) AttrNS(name, ns string) (xml.Attr, bool) {
	// Set key for attribute
	key := name
	if ns != "" {
		key = key + "," + ns
	}
	// Get attribute
	if attr, exists := this.attrs[key]; exists {
		return *attr, true
	} else {
		return xml.Attr{}, false
	}
}

func (this *Element) Children() []data.Node {
	result := make([]data.Node, 0, len(this.children))
	for _, child := range this.children {
		if child != nil {
			result = append(result, child.(data.Node))
		}
	}
	return result
}

func (this *Element) FirstChild() data.Node {
	for _, child := range this.children {
		if child != nil {
			return child.(data.Node)
		}
	}
	return nil
}

func (this *Element) LastChild() data.Node {
	for i := len(this.children) - 1; i >= 0; i-- {
		child := this.children[i]
		if child != nil {
			return child.(data.Node)
		}
	}
	return nil
}

func (this *Element) AddChild(child data.Node) error {
	return this.addChildBefore(child, nil)
}

func (this *Element) InsertChildBefore(child, ref data.Node) error {
	return this.addChildBefore(child, ref)
}

func (this *Element) RemoveAllChildren() error {
	for _, child := range this.children {
		if child != nil {
			// Detach child from parent
			switch node := child.(type) {
			case *Element:
				node.parent = nil
			case *Text:
				node.parent = nil
			case *Comment:
				node.parent = nil
			default:
				return data.ErrInternalAppError.WithPrefix("RemoveChildren")
			}
		}
	}

	// Empty array
	this.children = nil

	// Return success
	return nil
}

func (this *Element) GetElementsByTagName(name string) []data.Node {
	return this.GetElementsByTagNameNS(name, "")
}

func (this *Element) GetElementsByTagNameNS(name, ns string) []data.Node {
	xmlname := xml.Name{ns, name}
	result := make([]data.Node, 0, len(this.children))
	for _, child := range this.children {
		if child != nil {
			if element, ok := child.(*Element); ok {
				if element.XMLName == xmlname {
					result = append(result, child.(data.Node))
				}
			}
		}
	}
	return result
}

func (this *Element) RemoveChild(child data.Node) error {
	return this.removeChild(child)
}

func (this *Element) SetAttr(name string, value string) error {
	return this.SetAttrNS(name, "", value)
}

func (this *Element) SetAttrNS(name, ns string, value string) error {
	// Ensure attribute name is valid
	if reAttrName.MatchString(name) == false {
		return data.ErrBadParameter.WithPrefix("SetAttrNS ", strconv.Quote(name))
	}

	// Set key for attribute
	key := name
	if ns != "" {
		key = key + "," + ns
	}

	// Add attribute to order
	if _, exists := this.attrs[key]; exists == false {
		this.attro = append(this.attro, key)
	}

	// Set attribute for element
	this.attrs[key] = &xml.Attr{
		Name:  xml.Name{ns, name},
		Value: value,
	}

	// Attach attribute to document
	this.document.setAttr(name, ns, value, this)

	// Return success
	return nil
}

func (this *Element) String() string {
	if bytes, err := xml.Marshal(this); err != nil {
		return fmt.Sprintln("Error: ", err)
	} else {
		return string(bytes)
	}
}

/////////////////////////////////////////////////////////////////////
// XML ENCODING

func (this *Element) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {

	// Start element
	if err := enc.EncodeToken(xml.StartElement{
		Name: this.XMLName,
		Attr: this.Attrs(),
	}); err != nil {
		return err
	}

	// Children
	for _, child := range this.children {
		if child == nil {
			continue
		}
		switch node := child.(type) {
		case *Element:
			if err := enc.EncodeElement(node, xml.StartElement{Name: node.XMLName, Attr: node.Attrs()}); err != nil {
				return err
			}
		case *Text:
			if err := enc.EncodeToken(xml.CharData(node.cdata)); err != nil {
				return err
			}
		case *Comment:
			if err := enc.EncodeToken(xml.Comment(node.cdata)); err != nil {
				return err
			}
		default:
			return data.ErrBadParameter.WithPrefix("Invalid node: ", child)
		}
	}

	// End element
	if err := enc.EncodeToken(xml.EndElement{Name: this.XMLName}); err != nil {
		return err
	}

	// Flush
	if err := enc.Flush(); err != nil {
		return err
	}

	// Return success
	return nil
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func setParent(node data.Node, parent *Element) error {
	switch node := node.(type) {
	case *Element:
		node.parent = parent
	case *Text:
		node.parent = parent
	case *Comment:
		node.parent = parent
	default:
		return data.ErrBadParameter
	}
	return nil
}

func (this *Element) removeChild(node data.Node) error {
	for i, child := range this.children {
		if node == child {
			this.children[i] = nil
			return setParent(node, nil)
		}
	}
	// Child not found
	return data.ErrNotFound
}

// positionForChild returns the index in the array for a child node
// or returns -1
func (this *Element) positionForChild(ref data.Node) int {
	// Set position to end if ref is nil
	if ref == nil {
		return len(this.children)
	}
	// Determine position or return -1 if not found
	for i, child := range this.children {
		if child == ref {
			return i
		}
	}
	return -1
}

// insertChild at position performs the insert into an array of children
func (this *Element) insertChildAtPosition(pos int, child data.Node) {
	if pos == 0 {
		this.children = append([]interface{}{child}, this.children...)
	} else if pos == len(this.children) {
		this.children = append(this.children, child)
	} else {
		this.children = append(this.children[:pos+1], this.children[pos:]...)
		this.children[pos] = child
	}
}

func (this *Element) addChildBefore(node data.Node, ref data.Node) error {
	// Find position to insert child
	if refPos := this.positionForChild(ref); refPos == -1 {
		return data.ErrBadParameter.WithPrefix("InsertChildBefore")
	} else {
		// Detach node from existing parent
		if parent := node.Parent(); parent != nil {
			if err := parent.(*Element).removeChild(node); err != nil {
				return err
			}
		}

		// Insert child at position
		this.insertChildAtPosition(refPos, node)
	}

	// Set new parent
	return setParent(node, this)
}