package canvas

import (
	"encoding/xml"
	"fmt"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Element struct {
	XMLName   xml.Name
	XMLAttrs  []xml.Attr
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

func (e *Element) Attr(name string, value interface{}) {
	attr := xml.Attr{xml.Name{Local: name}, fmt.Sprint(value)}
	e.XMLAttrs = append(e.XMLAttrs, attr)
}

func (e *Element) Attrs() []xml.Attr {
	attrs := make([]xml.Attr, len(e.XMLAttrs))
	for i, attr := range e.XMLAttrs {
		attrs[i] = attr
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
	return attrs
}

func (e *Element) Group(children ...data.CanvasElement) data.CanvasGroup {
	g := NewElement("g", "", e.root)
	e.addChild(g)
	for _, node := range children {
		e.removeChild(node.(*Element))
		g.addChild(node.(*Element))
	}
	return g
}

func (e *Element) Id(value string) data.CanvasElement {
	e.Attr("id", value)
	return e
}

func (e *Element) Class(value string) data.CanvasElement {
	e.Attr("class", value)
	return e
}

func (e *Element) Transform(op ...data.CanvasTransform) data.CanvasElement {
	if e.transform == nil {
		e.transform = NewTransform(op)
	} else {
		e.transform.op = append(e.transform.op, op...)
	}
	return e
}

func (e *Element) Style(styles ...data.CanvasStyle) data.CanvasElement {
	if e.style == nil {
		e.style = NewStyle(styles)
	} else {
		e.style.Append(styles)
	}
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
	case "svg", "g": // Cases with children
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
	case "desc", "title": // Cases with cdata
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
