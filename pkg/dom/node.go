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

type Node struct {
	document *Document
	parent   *Element
}

type Element struct {
	Node
	XMLName  xml.Name
	attrs    map[string]*xml.Attr
	attro    []string
	children []interface{}
}

type Text struct {
	Node
	cdata []byte
}

type Comment struct {
	Node
	cdata []byte
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
// NODE METHODS

// These methods are default (Noop) versions

func (this *Node) Document() data.Document {
	return this.document
}

func (this *Node) Parent() data.Node {
	return this.parent
}

func (this *Node) Name() xml.Name {
	return xml.Name{}
}

func (this *Node) Attrs() []xml.Attr {
	return nil
}

func (this *Node) Attr(string) (xml.Attr, bool) {
	return xml.Attr{}, false
}

func (this *Node) AttrNS(string, string) (xml.Attr, bool) {
	return xml.Attr{}, false
}

func (this *Node) Cdata() string {
	return ""
}

func (this *Node) Children() []data.Node {
	return nil
}

func (this *Node) AddChild(child data.Node) error {
	return data.ErrInternalAppError
}

func (this *Node) RemoveChild(child data.Node) error {
	return data.ErrInternalAppError
}

func (this *Node) SetAttr(name string, value string) error {
	return data.ErrInternalAppError
}

func (this *Node) SetAttrNS(string, string, string) error {
	return data.ErrInternalAppError
}

func (this *Node) GetElementsByTagName(name string) []data.Node {
	return nil
}

func (this *Node) GetElementsByTagNameNS(string, string) []data.Node {
	return nil
}

func (this *Node) FirstChild() data.Node {
	return nil
}

func (this *Node) LastChild() data.Node {
	return nil
}

func (this *Node) InsertChildBefore(node, ref data.Node) error {
	return data.ErrInternalAppError
}

func (this *Node) RemoveChildren() error {
	return data.ErrInternalAppError
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
	switch node := child.(type) {
	case *Element:
		return this.addChildElementBefore(node, nil)
	case *Text:
		return this.addChildTextBefore(node, nil)
	case *Comment:
		return this.addChildCommentBefore(node, nil)
	default:
		return data.ErrInternalAppError.WithPrefix("AddChild")
	}
}

func (this *Element) InsertChildBefore(child, ref data.Node) error {
	switch node := child.(type) {
	case *Element:
		this.addChildElementBefore(node, ref)
	case *Text:
		this.addChildTextBefore(node, ref)
	case *Comment:
		this.addChildCommentBefore(node, ref)
	default:
		return data.ErrInternalAppError.WithPrefix("InsertChildBefore")
	}
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
	switch node := child.(type) {
	case *Element:
		return this.removeChildElement(node)
	case *Text:
		return this.removeChildText(node)
	case *Comment:
		return this.removeChildComment(node)
	default:
		return data.ErrInternalAppError.WithPrefix("RemoveChild")
	}
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
// TEXT METHODS

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
// COMMENT METHODS

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

func (this *Text) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	// Cdata
	if err := enc.EncodeElement(this.cdata, start); err != nil {
		return err
	}
	// Return success
	return nil
}

func (this *Comment) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	// Cdata
	if err := enc.EncodeElement(this.cdata, start); err != nil {
		return err
	}
	// Return success
	return nil
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Element) removeChildElement(node *Element) error {
	for i, child := range this.children {
		if node == child {
			this.children[i] = nil
			node.parent = nil
			return nil
		}
	}
	// Child not found
	return data.ErrNotFound
}

func (this *Element) removeChildText(node *Text) error {
	for i, child := range this.children {
		if node == child {
			this.children[i] = nil
			node.parent = nil
			return nil
		}
	}

	// Child not found
	return data.ErrNotFound
}

func (this *Element) removeChildComment(node *Comment) error {
	for i, child := range this.children {
		if node == child {
			this.children[i] = nil
			node.parent = nil
			return nil
		}
	}

	// Child not found
	return data.ErrNotFound
}

func (this *Element) addChildElementBefore(node *Element, ref data.Node) error {
	// Find position
	if ref != nil {
		// TODO
	}

	// Detach node from existing parent
	if node.parent != nil {
		if err := node.parent.removeChildElement(node); err != nil {
			return err
		}
	}

	// Append child
	this.children = append(this.children, node)
	node.parent = this

	// Return success
	return nil
}

func (this *Element) addChildTextBefore(node *Text, ref data.Node) error {
	// Detach node from existing parent
	if node.parent != nil {
		if err := node.parent.removeChildText(node); err != nil {
			return err
		}
	}

	// Append child
	this.children = append(this.children, node)
	node.parent = this

	// Return success
	return nil
}

func (this *Element) addChildCommentBefore(node *Comment, ref data.Node) error {
	// Detach node from existing parent
	if node.parent != nil {
		if err := node.parent.removeChildComment(node); err != nil {
			return err
		}
	}

	// Append child
	this.children = append(this.children, node)
	node.parent = this

	// Return success
	return nil
}
