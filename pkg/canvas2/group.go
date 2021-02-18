package canvas

import (
	"strings"

	"github.com/djthorpe/data"
)

func (this *Element) Desc(cdata string) data.CanvasGroup {
	cdata = strings.TrimSpace(cdata)

	// Remove existing desc tags
	if desc := this.Document.GetElementsByTagNameNS("desc", data.XmlNamespaceSVG); len(desc) != 0 {
		for _, child := range desc {
			this.Document.RemoveChild(child)
		}
	}

	// Create a new desc tag - put at top of element
	if cdata != "" {
		desc := this.Document.CreateElementNS("desc", data.XmlNamespaceSVG)
		if err := desc.AddChild(this.Document.CreateText(cdata)); err != nil {
			return nil
		} else if err := this.Document.InsertChildBefore(desc, this.Document.FirstChild()); err != nil {
			return nil
		}
		// If there is a title tag, then put the desc tag after the title tag
		if title := this.Document.GetElementsByTagNameNS("title", data.XmlNamespaceSVG); len(title) != 0 {
			this.Document.InsertChildBefore(desc, title[0].NextSibling())
		}
	}

	// Return success
	return this
}

func (this *Canvas) Group(children ...data.CanvasElement) data.CanvasGroup {
	g, err := this.NewElement("g")
	if err != nil {
		return nil
	}
	// Append children. If any children are nil, then return nil to bubble up
	// any errors
	for _, child := range children {
		if child == nil {
			return nil
		} else if err := g.AddChild(child.(*Element).Node); err != nil {
			return nil
		}
	}
	return g
}
