package data

import (
	"encoding/xml"
	"io"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type DOMOption uint8
type DOMValidateNodeFunc func(Node) error

/////////////////////////////////////////////////////////////////////
// INTERFACES

type Document interface {
	Node

	// Return nodes
	GetElementById(string) Node
	GetElementByIdNS(string, string) Node

	// Create nodes
	CreateElement(string) Node
	CreateElementNS(string, string) Node
	CreateText(string) Node
	CreateComment(string) Node

	// Write XML
	Write(io.Writer) error
}

type Node interface {
	Name() xml.Name
	Cdata() string
	Parent() Node
	GetElementsByTagName(string) []Node
	GetElementsByTagNameNS(string, string) []Node

	Children() []Node
	AddChild(Node) error
	RemoveChild(Node) error
	FirstChild() Node
	LastChild() Node
	InsertChildBefore(Node, Node) error
	RemoveAllChildren() error

	Attrs() []xml.Attr
	Attr(string) (xml.Attr, bool)
	AttrNS(string, string) (xml.Attr, bool)
	SetAttr(string, string) error
	SetAttrNS(string, string, string) error
	RemoveAttr(string) error
	RemoveAttrNS(string, string) error
}

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DOMWriteDirective    DOMOption = (1 << iota) // Write <?xml?> at top
	DOMWriteIndentTab                            // Indent output with tabs
	DOMWriteIndentSpace2                         // Indent output with two spaces
)

const (
	XmlNamespaceSVG   = "http://www.w3.org/2000/svg"   // Scalable Vector Graphics
	XmlNamespaceXLink = "http://www.w3.org/1999/xlink" // XLink
	XmlNamespaceXHTML = "http://www.w3.org/1999/xhtml" // XHTML
)

/////////////////////////////////////////////////////////////////////
// METHODS

func (o DOMOption) Is(f DOMOption) bool {
	return o&f == f
}
