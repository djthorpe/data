package table

import (
	"fmt"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/dom"
)

/////////////////////////////////////////////////////////////////////
// ASCII

func (t *Table) DOM(opts ...data.TableOpt) data.Document {
	// Set option flags
	t.applyOpt(opts)

	// Create a document
	dom := dom.NewDocumentNS("table", t.opts.ns)
	if dom == nil {
		return nil
	} else if t.opts.name != "" {
		if err := dom.SetAttr("id", t.opts.name); err != nil {
			return nil
		}
	}

	// Set header
	if t.hasOpt(optHeader) {
		head := dom.CreateElement("thead")
		row := dom.CreateElement("tr")
		if err := dom.AddChild(head); err != nil {
			return nil
		}
		if err := head.AddChild(row); err != nil {
			return nil
		}
		for _, name := range t.header.names() {
			th := dom.CreateElement("th")
			if err := th.AddChild(dom.CreateText(name)); err != nil {
				return nil
			}
			if err := row.AddChild(th); err != nil {
				return nil
			}
		}
	}

	// Return if no rows
	if len(t.r) == 0 {
		return dom
	}

	body := dom.CreateElement("tbody")
	result := make([]string, t.header.w)

	// Iterate through rows
	for i, r := range t.r {
		// Generate string form for row
		for j, v := range r.row(t.header.w) {
			if v_, err := t.outValue(i, j, v); err != nil {
				return nil
			} else if v__, ok := v_.(string); ok {
				result[j] = v__
			} else {
				result[j] = fmt.Sprint(v_)
			}
		}

		// Create row
		row := dom.CreateElement("tr")
		for _, cell := range result {
			td := dom.CreateElement("td")
			if err := td.AddChild(dom.CreateText(cell)); err != nil {
				return nil
			} else if err := row.AddChild(td); err != nil {
				return nil
			}
		}

		// Append row to body
		if err := body.AddChild(row); err != nil {
			return nil
		}
	}

	// Append body to table
	if err := dom.AddChild(body); err != nil {
		return nil
	}

	// Return the document
	return dom
}
