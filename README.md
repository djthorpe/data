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

### Table reading formats and options

### Table writing formats and options

### Transforming data values

