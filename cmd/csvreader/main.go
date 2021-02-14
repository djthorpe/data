package main

/*
	Read and merge CSV files and output result to terminal as
	ASCII table, or as CSV or SQL statements. Demonstrates
	reading and writing tables.
*/

import (
	"flag"
	"fmt"
	"os"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/table"
)

var (
	flagHeader    = flag.Bool("header", true, "CSV header")
	flagDelimiter = flag.String("delim", "", "CSV field delimiter")
	flagOutputCsv = flag.Bool("csv", false, "CSV output")
	flagOutputSql = flag.Bool("sql", false, "SQL output")
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "Missing filename argument")
		os.Exit(-1)
	}

	// Create table for data and set options for input and output
	// based on command line flags.
	t := table.NewTable()
	inOpts := []data.TableOpt{}
	outOpts := []data.TableOpt{}

	// Set options
	if *flagDelimiter != "" {
		d := []rune(*flagDelimiter)
		inOpts = append(inOpts, t.OptCsv(d[0]))
	}
	if *flagHeader {
		inOpts = append(inOpts, t.OptHeader())
		outOpts = append(outOpts, t.OptHeader())
	}
	switch {
	case *flagOutputCsv:
		outOpts = append(outOpts, t.OptCsv(0))
	case *flagOutputSql:
		outOpts = append(outOpts, t.OptSql("data"))
	default:
		outOpts = append(outOpts, t.OptAscii(0, data.BorderLines))
	}

	// Read CSV files
	for _, csv := range flag.Args() {
		if err := read(t, inOpts, csv); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}
	}

	// Write combined Ascii table
	if err := t.Write(os.Stdout, outOpts...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

func read(t data.Table, opts []data.TableOpt, path string) error {
	if r, err := os.Open(path); err != nil {
		return err
	} else {
		defer r.Close()
		return t.Read(r, opts...)
	}
}
