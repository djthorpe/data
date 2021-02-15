package canvas

import (
	"io"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/dom"
)

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

func (this *Canvas) readSVG(r io.Reader) error {
	if document, err := dom.Read(r, DOMOptions); err != nil {
		return err
	} else {
		this.Document = document
	}

	// TODO: Set other options for origin, size, etc.

	// Success
	return nil
}
