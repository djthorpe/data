---
description: 'XML Document ingestion, parsing, validation and transmission'
---

# The Document Object Model

The `dom` package implements some of the **Document Object Model** [DOM](https://en.wikipedia.org/wiki/Document_Object_Model). This package allows you to create, read and write XML formatted documents.

## Creating a Document

To create a new document, use the `NewDocument` or `NewDocumentNS` method:

```go
package main

import (
    "github.com/djthorpe/data"
    "github.com/djthorpe/data/pkg/dom"
)

func main() {
    d := dom.NewDocumentNS("svg", data.XmlNamespaceSVG, 0)
    if d == nil {
        panic("Unexpected nil return from NewDocument")
    }
    // ...
}
```

The third argument provides options which can affect how the document is parsed:

* `data.DOMWriteDirective` writes the XML declaration header when using `Write`;
* `data.DOMWriteIndentTab` indents XML output with tabs when using `Write`;
* `data.DOMWriteIndentSpace2` indents XML output with spaces when using `Write`.
* `data.DOMWriteDefaultNamespace` writes the default namespace attribute \(`xmlns="..."`)\ on the root element when it is set, otherwise omits it.

Elements, text and comments can be created using the following methods:

```go
type Document interface {
    CreateElement(name string) Node
    CreateElementNS(name,ns string) Node
    CreateText(cdata string) Node
    CreateComment(comment string) Node
    // ...
}
```

Introspection on a `data.Node` is provided by the following methods:

```go
type Node interface {
    Name() xml.Name
    Attrs() []xml.Attr
    Children() []Node
    Parent() Node
    Cdata() string
    // ...
}
```

Any node can be added to another node \(potentially detatching it from its' current parent\) or removed with the following methods:

```go
type Node interface {
    AddChild(Node) error
    RemoveChild(Node) error
    // ...
}
```

It is not possible to add a child to a text or comment node. Finally, an attribute can be set on an element node:

```go
type Node interface {
    SetAttr(name, value string) error
    SetAttrNS(name,ns,value string) error
    // ...
}
```

## Writing, Reading & Validating XML

The following method is provided for writing the document in XML:

```go
type Document interface {
    Write(io.Writer) error
    // ...
}
```

You can also use `xml.Marshal` from the standard library `encoding/xml` on the document or any node. To read, and parse a new document from a data stream, use the `Read` method, which returns errors if there was a parsing error:

```text
package main

import (
    "github.com/djthorpe/data"
    "github.com/djthorpe/data/pkg/dom"
)

func main() {
    d,err := dom.Read(os.Stdin, 0)
    if err != nil {
        panic(err)
    }
    // ...
}
```

Validation on child nodes can be performed using the `ReadEx` method. For example,

```text
package main

import (
    "github.com/djthorpe/data"
    "github.com/djthorpe/data/pkg/dom"
)

func main() {
    d,err := dom.ReadEx(os.Stdin, 0,func(data.Node) error {
        // Validate node and children here
        return nil
    })
    if err != nil {
        panic(err)
    }
    // ...
}
```

The node argument to the callback method are always elements \(as opposed to text and comment nodes\), and it is assumed your validation function will validate both the attributes and children of the node.

## Limitations

The following limitations apply to creating and parsing XML documents:

  * Declaring a __Namespace__ is only supported at the root element level
    when writing the namespace attribute
  * The prefix associated with a __Namespace__ cannot be user-defined;
  * Node types other than element, text or comment are not supported.
  * Validation of elements and attributes are currently not covered within
    this package. A schema package could be added later to perform validation.

