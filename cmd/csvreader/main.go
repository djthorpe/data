package main

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
)

func main() {
	opts := []data.TableOpt{}

	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "Missing filename argument")
		os.Exit(-1)
	}

	t := table.NewTable(data.ZeroSize)

	// Set options
	if *flagDelimiter != "" {
		d := []rune(*flagDelimiter)
		opts = append(opts, t.OptCsv(d[0]))
	}
	if *flagHeader {
		opts = append(opts, t.OptHeader())
	}

	// Read CSV files
	for _, csv := range flag.Args() {
		if err := read(t, opts, csv); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}
	}

	// Write combined Ascii table
	if err := t.Write(os.Stdout, t.OptHeader(), t.OptAscii(80, data.BorderLines)); err != nil {
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
