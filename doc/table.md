---
description: 'Data extraction, transformation and transmission'
---

# Tables

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

If the data cannot be represented, you can define your own transformation functions or it will be represented as a text \(`data.String`\) value otherwise. See below for more information about custom data transformation.

## Creating Tables

To create a table, use the `NewTable` method. You can define column headings when creating a table:

```text
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

## Reading Tables

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

## Writing Tables

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
* `table.OptSql(string)` sets the writing format to SQL with the provided argument as the table name. When using the `table.OptHeader` option, the CREATE TABLE statement is included. Currently the idea is to be compatible with **sqlite** rather than other SQL servers.
* `table.OptTransform(...TransformFunc)` sets one or more value transformation functions, which convert a native value into text. Any transform function can return `data.ErrSkipTransform` in order to move onto the next transform function.

## Introspection

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
    Name() string     // Column name
    Type() Type     // Types that the column represents
    Min() float64   // Minimum value of all numbers
    Max() float64   // Maximum value of all numbers
    Sum() float64   // Sum of all numbers
    Count() uint64  // Count of all numbers
    Mean() float64     // Mean average value of numbers or +Inf
}
```

## Parsing Custom Data

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

