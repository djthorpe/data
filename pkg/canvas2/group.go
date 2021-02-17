package canvas

import (
	"strings"

	"github.com/djthorpe/data"
)

func (this *Canvas) Desc(cdata string) data.CanvasGroup {
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
	g := this.Document.CreateElementNS("g", data.XmlNamespaceSVG)
	for _, child := range children {
		if err := g.AddChild(child.(*Element)); err != nil {
			return nil
		}
	}
	if err := this.Document.AddChild(g); err != nil {
		return nil
	}
	return g
}
