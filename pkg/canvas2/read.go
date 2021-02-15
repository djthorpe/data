package canvas

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/dom"
)

/////////////////////////////////////////////////////////////////////
// CONSTANTS

var (
	tags = map[xml.Name]data.DOMValidateNodeFunc{
		{data.XmlNamespaceSVG, "svg"}:   tagSVG,
		{data.XmlNamespaceSVG, "title"}: tagTitle,
		{data.XmlNamespaceSVG, "path"}:  tagPath,
		{data.XmlNamespaceSVG, "g"}:     tagGroup,
		{data.XmlNamespaceSVG, "desc"}:  tagDesc,
	}
)

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func Read(fmt data.Writer, r io.Reader) (data.Canvas, error) {
	this := new(Canvas)

	switch fmt {
	case data.SVG:
		if err := this.readSVG(r); err != nil {
			return nil, err
		}
	default:
		return nil, data.ErrNotImplemented
	}

	// Return success
	return this, nil
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Canvas) readSVG(r io.Reader) error {
	if document, err := dom.ReadEx(r, DOMOptions, this.validateSVG); err != nil {
		return err
	} else {
		this.Document = document
	}

	// TODO: Set other options for origin, size, etc.

	// Success
	return nil
}

func (this *Canvas) validateSVG(node data.Node) error {
	name := node.Name()
	if fn, exists := tags[name]; exists {
		return fn(node)
	} else {
		return data.ErrBadParameter.WithPrefix("Unsupported tag: ", strconv.Quote(name.Local))
	}
}

func tagSVG(node data.Node) error {
	fmt.Println("VALIDATE SVG", node)
	return nil
}

func tagTitle(node data.Node) error {
	fmt.Println("VALIDATE TITLE", node)
	return nil
}

func tagPath(node data.Node) error {
	fmt.Println("VALIDATE PATH", node)
	return nil
}

func tagGroup(node data.Node) error {
	fmt.Println("VALIDATE GROUP", node)
	return nil
}

func tagDesc(node data.Node) error {
	fmt.Println("VALIDATE DESC", node)
	return nil
}
