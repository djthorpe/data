# Data Transformation

This repository contains various data extraction, transformation processing and visualization tools. Currently it contains the following:

  * [`data.Table`](#tables) provides you with a way to ingest, transform and process data tables in comma-separated value format and output in CSV, ASCII and SQL formats;
  * [`data.DOM`](#dom) provides a document object model which can read and write the XML format in addition to validating
  the XML;
  * [`data.Canvas`](#canvas) provides a drawing canvas on which graphics primitives such as lines, circles, text and rectangles can be placed. Additionally transformation, grouping and stylizing of primitives can be applied. Canvases can currently be written in SVG format, the intention is to also allow rendering using OpenGL later.

## Tables

You can ingest tables from files and other sources by creating a table object and then reading using the `io.Reader` class. For example:

```go
package main

import (
	"os"
	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/table"
)

func main() {
	t := table.NewTable()

	// Read CSV table
	r, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)      
	}
	defer r.Close()
	if err := t.Read(r,OptHeader()); err != nil {
		panic(err)      
	}

	// Write Ascii table
	if err := t.Write(os.Stdout, t.OptHeader(), t.OptAscii(80, data.BorderLines)); err != nil {
		panic(err)
	}
}
```

Internally, cells within the table are stored in "native" format, a default transformer is applied which can interpret from text:

  * `data.Nil` values when the text is empty;
  * `data.Uint` when the text is a positive numerical value;
  * `data.Int` when the text is any whole numerical value;
  * `data.Float` when the text represents any other pure number;
  * `data.Duration` when the text represents number of hours, seconds, milli or nanoseconds;
  * `data.Date` when the text represents a date;
  * `data.Datetime` when the text represents a date and time;
  * `data.Bool` when the text represents a boolean value;
  * `data.Other` when the native value is not represented otherwise.

If the data cannot be represented, you can define your own transformation functions or it will be represented as a text (`data.String`) value otherwise. See below for more information about custom data transformation.

### Creating tables

To create a table, use the `NewTable` method. You can define column headings when creating a table:

```
package main

import (
	"github.com/djthorpe/data/pkg/table"
)

func main() {
	t := table.NewTable("A","B","C")

	t.Append(2,"Hello",true)
	t.Append(1,56,false)
	t.Sort(func(a,b []interface{}) bool {
		return a[0] < b[0]
	})

	// ...
}
```

You can append native values to the table, which may extend the width of the table as well as the length. To sort the table, provide a row comparison function which accepts two arguments `a` and `b` and returns `true` if `a < b`.

### Table reading formats and options

To read data from an external source, use the `table.Read` method with reading options:

```go
type Table interface {
	Read(io.Reader, ...TableOpt) error
	// ...
}
```

The following table options are relevant for reading:

  * `table.OptHeader()` indicates the CSV file has a header row;
  * `table.OptCsv(rune)` sets the delimiter used for separating values on a row;
  * `table.OptType(data.Type)` sets the types which can be transformed from text. Use `data.DefaultTypes` for the default set of transformations. If the text cannot be transformed into one of the listed types, the value is stored as text;
  * `table.OptDuration(time.Duration)` sets the duration units for any text. For example if setting to time.Hour then "30m" is transformed to "0h" and "5" is transformed into "5h";
  * `table.OptTimezone(tz *time.Location)` sets the timezone for any transformed dates and times which do not explicitly set the timezone;
  * `table.OptRowIterator(IteratorFunc)` sets a row iterator, which is called before the row is added to the table. The iterator function can return `data.ErrSkipTransform` in order to skip adding the row to the table.
  * `table.OptTransform(...TransformFunc)` sets one or more value transformation functions, which convert a text into native value. Any transform function can return `data.ErrSkipTransform` in order to move onto the next transform function.

If you read multiple sets of data into a single table, extending the table in both width and height as necessary.

### Table writing formats and options

To write data to an external source, use the `table.Write` method with writing options:

```go
type Table interface {
	Write(io.Writer, ...TableOpt) error
	// ...
}
```

The following table options are relevant for writing:

  * `table.OptHeader()` indicates the header of the table should be output first;
  * `table.OptCsv(rune)` sets the writing format to CSV and sets the delimiter used for separating values on a row. When argument is zero, a comma is used;
  * `table.OptAscii(int,string)` sets the writing format to ASCII and sets the maximum width of the table in characters. When the width is zero, the table width is unbounded. The second argument can be set to `data.BorderDefault` for ASCII border characters or `data.BorderLines` for UTF8 border characters.
  * `table.OptSql(string)` sets the writing format to SQL with the provided argument as the table name. When using the `table.OptHeader` option, the CREATE TABLE statement is included. Currently the idea is to be compatible with __sqlite__ rather than other SQL servers.
  * `table.OptTransform(...TransformFunc)` sets one or more value transformation functions, which convert a native value into text. Any transform function can return `data.ErrSkipTransform` in order to move onto the next transform function.

### Table introspection

The following methods return introspection on a table:

```go
type Table interface {
	// Len returns the number of rows
	Len() int

	// Row returns values for a zero-indexed table
	// row or nil if the row does not exist
	Row(int) []interface{}

	// Col returns column information for a
	// zero-indexed table column or nil if the 
	// column doesn't exist
	Col(int) TableCol
}
```

For example, to iterate over rows in the table:

```go
package main

import (
	"github.com/djthorpe/data/pkg/table"
)

func main() {
	t := table.NewTable("A","B","C")
	// ...
	for i := 0; i < t.Len(); i++ {
		fmt.Println(t.Row(i))
	}
}
```

Information about the columns is returned with the `table.Col(int)` method, or the method returns nil if a column doesn't exist. The interface for `data.TableCol` is as follows:

```go
type TableCol interface {
	Name() string 	// Column name
	Type() Type 	// Types that the column represents
	Min() float64   // Minimum value of all numbers
	Max() float64   // Maximum value of all numbers
	Sum() float64   // Sum of all numbers
	Count() uint64  // Count of all numbers
	Mean() float64 	// Mean average value of numbers or +Inf
}
```

### Parsing custom formats

You can define custom formats for parsing using the `OptTransform` option when reading a CSV file. For example, to parse IP addresses, a transformation function can be defined:

```go
func ParseIP(t data.Table) data.TransformFunc {
	return func(i, j int, v interface{}) (interface{}, error) {
		if ip := net.ParseIP(v.(string)); ip == nil {
			return nil, data.ErrSkipTransform
		} else {
			return ip, nil
		}
	}
}

func main() {
	c := table.NewTable("IP Address")
	if err := c.Read(strings.NewReader(os.Args[1]), c.OptTransform(ParseIP(c))); err != nil {
		t.Error(err)
	}
	// ...
}
```

A similar transformation function can be defined on output, for converting a native format to a string.

## DOM

The DOM package allows you to create, read and write XML formatted documents. To create a new document, use the `NewDocument` or `NewDocumentNS` method:

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

Any node can be added to another node (potentially detatching it from its' current parent) or removed with the following methods:

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

### Writing, Reading & Validating XML

The following method is provided for writing the document in XML:

```go
type Document interface {
	Write(io.Writer) error
	// ...
}
```

You can also use `xml.Marshal` from the standard library `encoding/xml` on the document or any node. To read, and parse a new document from a data stream, use the `Read` method, which returns errors if there was a parsing error:

```
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

```
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

The node argument to the callback method are always elements (as opposed to text and comment nodes), and it is assumed your validation function will validate both the attributes and children of the node.

## Canvas

You can create canvas objects by using the `NewCanvas` method and then create shape primitives, styles, transforms and groups on the canvas before writing out the canvas to render it. Currently SVG is supported for rendering. For example, the following code writes out a canvas as SVG of A4 paper size:

```go
package main

import (
	"os"
	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/canvas"
	"github.com/djthorpe/data/pkg/color"
)

func main() {
	c := canvas.NewCanvas(data.A4LandscapeSize, data.MM)

	c.Title("Canvas Document")
	c.Desc("SVG Document output")
	c.Rect(c.Origin(), c.Size()).Style(
		c.Fill(color.Ivory, 1.0),
	)
	c.Write(data.SVG, os.Stdout)
}
```

Canvas can consist of:

  * `data.CanvasElement`, which are primitive drawing elements like circles, text, lines, paths and text;
  * `data.CanvasGroup`, which are groups of elements.
  
Any element or group can have __transforms__ applied (for example, scaling, rotating and translation) and
__style__ (for example, color or opacity).

### Creating a Canvas

Canvases are created with the size of the canvas. There are several predefined sizes:

  * `canvas.NewCanvas(data.A4LandscapeSize, data.MM)` defines an A4 canvas in landscape;
  * `canvas.NewCanvas(data.A4PortraitSize, data.MM)` defines an A4 canvas in portrait;
  * `data.LetterPortraitSize` and `data.LetterLandscapeSize` can be used as the first argument for Letter sized paper;
  * `data.LegalPortraitSize` and `LegalLandscapeSize` can be used to define a canvas the size of Legal paper.

The second argument are the units:

| Constant | Unit | Description |
| -------- | ---- | ----------- |
| `data.None` | Natural | Based on screen size
| `data.PX` | Pixels | Based on screen size
| `data.CM` | Centimetre (cm) | 
| `data.MM` | Millimetre (mm) | 10mm = 1cm
| `data.IN` | Inch | 1in = 2.54cm
| `data.PC` | Picas | 6pc = 1 in
| `data.PT` | Points | 72pt ≅ 1 in
| `data.EM` | | Equivalent to the computed font-size
| `data.EX` | | Equivalent to the height of a lower-case letter

In addition to setting canvas size, it's possible to set the translation to natural units used to render elements, using `Size`, `SetSize`, `Origin` and `SetOrigin`:

```go
type Canvas interface {
	Origin() Point
	Size() Size
	SetOrigin(Point)
	SetSize(Size)
	// ...
}
```

This adjusts the top left co-ordinate to be equal to the __origin__ and the width and height to be equal to the __size__
where natural units are used. For example,

```go
	c := canvas.NewCanvas(data.LetterPortraitSize, data.MM)
	c.SetSize(geom.DivideSize(c.Size(), 0.3527))
```

This example creates a canvas of 612 units across and 792 units high, mapped onto a canvas which would fit on letter paper: `0.3527` is approximately 2.84 natural units per millimetre, which is 72 units per inch.
 
### Adding shape primitives

Shape primitives are created using the following canvas methods:

```go
type Canvas interface {
	Circle(Point, float32) CanvasElement
	Ellipse(Point, Size) CanvasElement
	Line(Point, Point) CanvasElement
	Rect(Point, Size) CanvasElement
	Text(Point, ...CanvasText) CanvasElement
	Path([]Point) CanvasPath
	// ...
}
```

For example calling `c.Circle(data.Point{ 0,0 },10 })` creates a circle at the top left origin of natural radius 10, and adds the circle to the canvas. The arguments for each method are:


| Method | Arguments | Description |
| -------- | ---- | ----------- |
| `data.Circle` | `{x-origin,y-origin},radius` | Circle  
| `data.Ellipse` | `{x-origin,y-origin},{x-radius,y-radius}` | Ellipse
| `data.Line`    | `{x-start,y-start},{x-end,y-end}` | Line
| `data.Rect`    | `{x-start,y-start},{x-end,y-end}` | Rectangle
| `data.Text`    | `{x,y},...string` | Adds text around point
| `data.Path`    | `{x-move,y-move},...{x-line,{y-line}` | Complex path. See below for bezier curves

TODO


### Styling 

TODO

### Transformation 

TODO

### Grouping 

TODO

### Rendering

TODO





