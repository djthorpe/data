# Data Transformation

This repository contains various data extraction, transformation
processing and visualization tools. Currently it contains the 
following:

* `data.Table` provides you with a way to ingest, transform
	and process data tables in comma-separated value format.

## Tables

You can ingest tables from files and other sources by creating
a table object and then reading using the `io.Reader` class. For
example:

```go
package main

import (
	"os"
	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/table"
)

func main() {
	t := table.NewTable(data.ZeroSize)

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

Internally, cells within the table are stored in "native" format, a default transformer is
applied which can interpret from text:

	* `data.Nil` values when the text is empty;
	* `data.Uint` when the text is a positive numerical value;
	* `data.Int` when the text is any whole numerical value;
	* `data.Float` when the text represents any other pure number;
	* `data.Duration` when the text represents number of hours, seconds, milli or nanoseconds;
	* `data.Date` when the text represents a date;
	* `data.Datetime` when the text represents a date and time;
	* `data.Bool` when the text represents a boolean value.
	
If the data cannot be represented, you can define your own transformation functions or
it will be represented as a text (`data.String`) value otherwise. See below for more
information about custom data transformation.

### Creating tables

To create a table, use the `NewTable` method
with the initial size of the table:

```
package main

import (
	"github.com/djthorpe/data/pkg/table"
)

func main() {
	t := table.NewTable(data.ZeroSize)

	t.Append(2,"Hello",true)
	t.Append(1,56,false)

	t.Sort(func(a,b []interface{}) bool {
		return a[0] < b[0]
	})

	// ...
}
```

You can append native values to the table, which may extend the width of the table as well as the length. To
sort the table, provide a row comparison function which
accepts two arguments `a` and `b` and returns `true` if
`a < b`.

### Table reading formats and options

To read data from an external source, use the `table.Read` method with reading options:

```go
type Table interface {
	Read(io.Reader, ...TableOpt) error
	// ...
}
```

The following table options are relevant for reading:

* `table.OptHeader()` indicates the CSV file has
	a header row;
* `table.OptCsv(rune)` sets the delimiter used for
	separating values on a row;
* `table.OptType(data.Type)` sets the types which can
	be transformed from text. Use `data.DefaultTypes`
	for the default set of transformations. If the text
	cannot be transformed into one of the listed types,
	the value is stored as text;
* `table.OptDuration(time.Duration)` sets the
	duration units for any text. For example if setting
	to time.Hour then "30m" is transformed to "0h" and
	"5" is transformed into "5h";
* `table.OptTimezone(tz *time.Location)` sets the
	timezone for any transformed dates and times which
	do not explicitly set the timezone;
* `table.OptRowIterator(IteratorFunc)` sets a row
	iterator, which is called before the row is added to
	the table. The iterator function can return `data.ErrSkipTransform` in order to skip adding the row
	to the table.
* `table.OptTransform(...TransformFunc)` sets one or
	more value transformation functions, which convert a text into native value. Any transform function can return `data.ErrSkipTransform` in order to move onto the next transform function.

If you read multiple sets of data into a single table, extending the table in both width
and height as necessary.

### Table writing formats and options

To write data to an external source, use the `table.Write` method with writing options:

```go
type Table interface {
	Write(io.Writer, ...TableOpt) error
	// ...
}
```

The following table options are relevant for writing:

* `table.OptHeader()` indicates the header of the table
	should be output first;
* `table.OptCsv(rune)` sets the writing format to CSV and
	sets the delimiter used for separating values on a row;
* `table.OptTransform(...TransformFunc)` sets one or
	more value transformation functions, which convert a native 
	value into text. Any transform function can return
	`data.ErrSkipTransform` in order to move onto the next 
	transform function.

### Transforming data values

