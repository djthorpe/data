package canvas

import (
	"encoding/xml"
	"fmt"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Element struct {
	XMLName xml.Name
	Attrs   []xml.Attr

	children []*Element
	cdata    string
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewElement(name string, cdata string) *Element {
	return &Element{
		XMLName: xml.Name{"", name},
		cdata:   cdata,
	}
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (e *Element) Desc(value string) data.CanvasGroup {
	e.addChild(NewElement("desc", value))
	return e
}

func (e *Element) Attr(name string, value interface{}) {
	attr := xml.Attr{xml.Name{Local: name}, fmt.Sprint(value)}
	e.Attrs = append(e.Attrs, attr)
}

func (e *Element) Group(children ...data.CanvasElement) data.CanvasGroup {
	g := NewElement("g", "")
	e.addChild(g)
	for _, node := range children {
		g.addChild(node.(*Element))
	}
	return e
}

func (e *Element) Id(value string) {
	e.Attr("id", value)
}

func (e *Element) Class(value string) {
	e.Attr("class", value)
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (e *Element) String() string {
	str := "<" + e.XMLName.Local
	for _, attr := range e.Attrs {
		str += fmt.Sprintf(" %v=%q", attr.Name.Local, attr.Value)
	}
	if len(e.children) > 0 {
		str += " <" + fmt.Sprint(e.children) + ">"
	}
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// XML ENCODING

func (e *Element) MarshalXML(x *xml.Encoder, start xml.StartElement) error {
	switch e.XMLName.Local {
	case "svg", "g":
		x.EncodeToken(xml.StartElement{
			Name: e.XMLName,
			Attr: e.Attrs,
		})
		for _, c := range e.children {
			x.EncodeElement(c, xml.StartElement{Name: c.XMLName})
		}
		x.EncodeToken(xml.EndElement{Name: e.XMLName})
		return nil
	case "circle", "rect", "path":
		x.EncodeToken(xml.StartElement{
			Name: e.XMLName,
			Attr: e.Attrs,
		})
		x.EncodeToken(xml.EndElement{Name: e.XMLName})
		return nil
	case "desc", "title":
		return x.EncodeElement(e.cdata, xml.StartElement{
			Name: e.XMLName,
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
