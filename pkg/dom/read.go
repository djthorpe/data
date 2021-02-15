package dom

import (
	"bytes"
	"encoding/xml"
	"io"
	"strconv"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPE

type context struct {
	document *Document
	parent   []*Element
}

/////////////////////////////////////////////////////////////////////
// XML DECODING

func (this *Document) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	ctx := new(context)
	ctx.document = this

	// Set root tag and attributes
	if err := ctx.StartElement(start, this.Element); err != nil {
		return err
	}

	// Read tokens until end
	for {
		t, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		switch tok := t.(type) {
		case xml.StartElement:
			if err := ctx.StartElement(tok, nil); err != nil {
				return err
			}
		case xml.CharData:
			if err := ctx.Text(tok); err != nil {
				return err
			}
		case xml.EndElement:
			if node, err := ctx.EndElement(tok); err != nil {
				return err
			} else if this.fn == nil {
				// no-op
			} else if err := this.fn(node); err != nil {
				return err
			}
		case xml.Comment:
			if err := ctx.Comment(tok); err != nil {
				return err
			}
		default:
			return data.ErrBadParameter.WithPrefix(t)
		}
	}

	// Return success
	return nil
}

/////////////////////////////////////////////////////////////////////
// DECODING CONTEXT

func (ctx *context) IsRoot() bool {
	return len(ctx.parent) == 0
}

func (ctx *context) Parent() *Element {
	if ctx.IsRoot() {
		return nil
	} else {
		return ctx.parent[len(ctx.parent)-1]
	}
}

func (ctx *context) Pop() *Element {
	if parent := ctx.Parent(); parent == nil {
		return nil
	} else {
		ctx.parent = ctx.parent[:len(ctx.parent)-1]
		return parent
	}
}

func (ctx *context) StartElement(start xml.StartElement, node *Element) error {
	// Create element
	if node == nil {
		node = NewElementNS(start.Name.Local, start.Name.Space, nil, ctx.document)
	} else {
		node.XMLName = start.Name
	}

	// Copy attributes
	for _, attr := range start.Attr {
		node.SetAttrNS(attr.Name.Local, attr.Name.Space, attr.Value)
	}

	// Add element to parent
	if parent := ctx.Parent(); parent != nil {
		if err := parent.addChildElement(node); err != nil {
			return err
		}
	}

	// Set element as new parent
	ctx.parent = append(ctx.parent, node)

	// Return success
	return nil
}

func (ctx *context) EndElement(end xml.EndElement) (*Element, error) {
	if parent := ctx.Pop(); parent == nil {
		return nil, data.ErrBadParameter.WithPrefix("Invalid end tag outside document: ", strconv.Quote(end.Name.Local))
	} else if parent.XMLName != end.Name {
		return nil, data.ErrBadParameter.WithPrefix("Non-matching end tag: ", strconv.Quote(end.Name.Local))
	} else {
		return parent, nil
	}
}

func (ctx *context) Text(cdata xml.CharData) error {
	if ctx.IsRoot() {
		// Ignore whitespace when outside document
		if cdata := bytes.TrimSpace(cdata); len(cdata) != 0 {
			return data.ErrBadParameter.WithPrefix("Invalid cdata outside document: ", strconv.Quote(string(cdata)))
		}
	}

	// Do not preserve whitespace
	if len(bytes.TrimSpace(cdata)) == 0 {
		return nil
	}

	// Create node, add to parent
	node := NewText(cdata, nil, ctx.document)
	if parent := ctx.Parent(); parent != nil {
		if err := parent.addChildText(node); err != nil {
			return err
		}
	}

	// Return success
	return nil
}

func (ctx *context) Comment(cdata xml.Comment) error {
	// Create node, add to parent
	node := NewComment(cdata, nil, ctx.document)
	if parent := ctx.Parent(); parent != nil {
		if err := parent.addChildComment(node); err != nil {
			return err
		}
	}

	// Return success
	return nil
}
