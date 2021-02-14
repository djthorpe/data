package canvas

import (
	"encoding/xml"
	"fmt"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/f32"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Element struct {
	XMLName   xml.Name
	XMLAttrs  []*xml.Attr
	root      *Element
	children  []*Element
	cdata     string
	transform *Transform
	style     *Style

	// Canvas properties
	origin data.Point
	size   data.Size
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewElement(name string, cdata string, root *Element) *Element {
	return &Element{
		XMLName: xml.Name{"", name},
		cdata:   cdata,
		root:    root,
	}
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (e *Element) Desc(value string) data.CanvasGroup {
	e.addChild(NewElement("desc", value, e.root))
	return e
}

func (e *Element) getAttr(name string) *xml.Attr {
	for _, attr := range e.XMLAttrs {
		if attr.Name.Local == name {
			return attr
		}
	}
	return nil
}

func (e *Element) Attr(name, value string) {
	if attr := e.getAttr(name); attr == nil {
		e.XMLAttrs = append(e.XMLAttrs, &xml.Attr{
			Name:  xml.Name{Local: name},
			Value: value,
		})
	} else {
		attr.Value += " " + value
	}
}

func (e *Element) Attrs() []xml.Attr {
	attrs := make([]xml.Attr, len(e.XMLAttrs))
	for i, attr := range e.XMLAttrs {
		attrs[i] = *attr
	}
	if e.transform != nil {
		if value := e.transform.String(); value != "" {
			attrs = append(attrs, xml.Attr{
				Name:  xml.Name{Local: "transform"},
				Value: value,
			})
		}
	}
	if e.style != nil {
		if value := e.style.String(); value != "" {
			attrs = append(attrs, xml.Attr{
				Name:  xml.Name{Local: "style"},
				Value: value,
			})
		}
	}
	if e.size != data.ZeroSize {
		attrs = append(attrs, xml.Attr{
			Name: xml.Name{Local: "viewBox"},
			Value: fmt.Sprint(
				f32.String(e.origin.X), " ",
				f32.String(e.origin.Y), " ",
				f32.String(f32.Abs(e.size.W)), " ",
				f32.String(f32.Abs(e.size.H)),
			),
		})
	}
	return attrs
}

func (e *Element) Id(value string) data.CanvasElement {
	e.Attr("id", value)
	return e
}

func (e *Element) Class(value string) data.CanvasElement {
	e.Attr("class", value)
	return e
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (e *Element) String() string {
	str := "<" + e.XMLName.Local
	for _, attr := range e.XMLAttrs {
		str += fmt.Sprintf(" %v=%q", attr.Name.Local, attr.Value)
	}
	if e.transform != nil {
		str += fmt.Sprintf(" transform=%q", e.transform.String())
	}
	if e.style != nil {
		str += fmt.Sprintf(" style=%q", e.style.String())
	}
	if len(e.children) > 0 {
		str += " <"
		for _, c := range e.children {
			if c != nil {
				str += c.String()
			}
		}
		str += ">"
	}
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// XML ENCODING

func (e *Element) MarshalXML(x *xml.Encoder, start xml.StartElement) error {
	switch e.XMLName.Local {
	case "svg", "g", "text": // Cases with children
		x.EncodeToken(xml.StartElement{
			Name: e.XMLName,
			Attr: e.Attrs(),
		})
		for _, c := range e.children {
			if c != nil {
				x.EncodeElement(c, xml.StartElement{Name: c.XMLName})
			}
		}
		x.EncodeToken(xml.EndElement{Name: e.XMLName})
		return nil
	case "circle", "rect", "path", "line": // Cases without children
		x.EncodeToken(xml.StartElement{
			Name: e.XMLName,
			Attr: e.Attrs(),
		})
		x.EncodeToken(xml.EndElement{Name: e.XMLName})
		return nil
	case "desc", "title", "tspan": // Cases with cdata
		return x.EncodeElement(e.cdata, xml.StartElement{
			Name: e.XMLName,
			Attr: e.Attrs(),
		})
	default:
		return data.ErrBadParameter.WithPrefix(e.XMLName.Local)
	}
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (e *Element) addChild(c *Element) {
	e.children = append(e.children, c)
}

func (e *Element) removeChild(c *Element) {
	for i, child := range e.children {
		if child == c {
			e.children[i] = nil
		}
	}
}
