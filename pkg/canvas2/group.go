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
	}

	// Return success
	return this
}

func (this *Canvas) Group(children ...data.CanvasElement) data.CanvasGroup {

}
