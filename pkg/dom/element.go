package dom

import (
	"encoding/xml"
	"fmt"
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

	// Register namespace in document if not empty
	if ns != "" {
		element.document.setTagNS(element, ns)
	}

	// Return success
	return element
}

/////////////////////////////////////////////////////////////////////
// ELEMENT METHODS

func (this *Element) Name() xml.Name {
	return this.XMLName
}

func (this *Element) Attrs() []xml.Attr {
	attrs := make([]xml.Attr, 0, len(this.attro))

	// Add namespace attributes onto the attributes for the root element
	if this.IsRootElement() {
		for ns, tag := range this.document.tag {
			if tag != "" {
				attrs = append(attrs, xml.Attr{
					Name:  xml.Name{Space: "", Local: "xmlns:" + tag},
					Value: ns,
				})
			}
		}
	}

	// Convert attribute names to shorten the namespace prefix
	for _, key := range this.attro {
		// Removed attributes have empty string
		if key != "" {
			// Get attribute
			attr := this.attrs[key]
			if name, err := this.document.getTagNS(this, attr.Name); err != nil {
				// Let's just not clean up if there are errors
				attrs = append(attrs, *attr)
			} else {
				// Or else append a new shortened tag with prefix
				attrs = append(attrs, xml.Attr{
					Name:  name,
					Value: attr.Value,
				})
			}
		}
	}

	return attrs
}

func (this *Element) Attr(name string) (xml.Attr, bool) {
	return this.AttrNS(name, "")
}

func (this *Element) AttrNS(name, ns string) (xml.Attr, bool) {
	key := attrKey(name, ns)
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

func (this *Element) PrevSibling() data.Node {
	return prevSibling(this)
}

func (this *Element) NextSibling() data.Node {
	return nextSibling(this)
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

	// Register tag for ns
	if ns != "" {
		prefix := this.document.setTagNS(this, ns)
		fmt.Println("TODO SET NS for ATTR", name, "=>", ns, "=>", prefix)
	}

	// Add attribute to order
	key := attrKey(name, ns)
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

func (this *Element) RemoveAttr(name string) error {
	return this.RemoveAttrNS(name, "")
}

func (this *Element) RemoveAttrNS(name, ns string) error {
	// Check for existence of attribute
	key := attrKey(name, ns)
	if _, exists := this.attrs[key]; exists == false {
		return data.ErrNotFound.WithPrefix("RemoveAttrNS: ", strconv.Quote(name))
	} else {
		delete(this.attrs, key)
	}

	// Remove from attribute order
	for i := range this.attro {
		if this.attro[i] == key {
			this.attro[i] = ""
		}
	}

	// Return success
	return nil
}

func (this *Element) IsRootElement() bool {
	return this == this.document.Element
}

func (this *Element) String() string {
	if bytes, err := xml.Marshal(this); err != nil {
		panic(err)
	} else {
		return string(bytes)
	}
}

/////////////////////////////////////////////////////////////////////
// XML ENCODING

func (this *Element) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	name, err := this.document.getTagNS(this, this.XMLName)
	if err != nil {
		return err
	}
	if err := enc.EncodeToken(xml.StartElement{
		Name: name,
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
			if err := enc.EncodeElement(node, xml.StartElement{Name: name, Attr: node.Attrs()}); err != nil {
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
	if err := enc.EncodeToken(xml.EndElement{Name: name}); err != nil {
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

func attrKey(name, ns string) string {
	key := name
	if ns != "" {
		key = key + "," + ns
	}
	return key
}

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
func (this *Element) positionForChild(child data.Node) int {
	// Set position to end if ref is nil
	if child == nil {
		return len(this.children)
	}
	// Determine position or return -1 if not found
	for i, node := range this.children {
		if child == node {
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
