package dom

import (
	"github.com/djthorpe/data"
)

var (
	// Ref: https://www.informit.com/articles/article.aspx?p=31837&seqNum=10
	xmlNs = map[string]string{
		data.XmlNamespaceSVG:                         "svg",
		data.XmlNamespaceXLink:                       "xlink",
		data.XmlNamespaceXHTML:                       "xhtml",
		"http://www.w3.org/XML/1998/namespace":       "xml",
		"http://www.w3.org/2001/XInclude":            "xi",
		"http://www.w3.org/1999/XSL/Format":          "fo",
		"http://www.w3.org/1999/XSL/Transform":       "xsl",
		"http://icl.com/saxon":                       "saxon",
		"http://xml.apache.org/xslt":                 "xalan",
		"http://www.w3.org/2001/XMLSchema":           "xsd",
		"http://www.w3.org/2001/XMLSchema-datatypes": "dt",
		"http://www.w3.org/2001/XMLSchema-instance":  "xsi",
		"http://www.w3.org/2000/01/rdf-schema":       "rdf",
		"http://www.w3.org/2001/SMIL20":              "smil",
		"http://www.w3.org/1998/Math/MathML":         "m",
	}
)
