package dom

import (
	"encoding/xml"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Node struct {
	document *Document
	parent   *Element
}

/////////////////////////////////////////////////////////////////////
// NODE METHODS

// These methods are default (Noop) versions

func (this *Node) Document() data.Document {
	return this.document
}

func (this *Node) Parent() data.Node {
	if this.parent == nil {
		return nil
	} else {
		return this.parent
	}
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

func (this *Node) RemoveAllChildren() error {
	return data.ErrInternalAppError
}
