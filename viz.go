package data

import "io"

/////////////////////////////////////////////////////////////////////
// INTERFACES

type Viz interface {
	// Create graph paper with major x minor grid squares
	GraphPaper(major, minor uint) VizGraphPaper

	// Create table
	// TODO Table(Table, ...TableOpt) VizTable

	// Write to data stream
	Write(Writer, io.Writer) error
}

type VizGraphPaper interface{}
type VizTable interface{}

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// Class names used for GraphPaper
	ClassGraphPaper       = "graphpaper"
	ClassGraphPaperBorder = "border"
	ClassGraphPaperXMajor = "majorx"
	ClassGraphPaperYMajor = "majory"
	ClassGraphPaperXMinor = "minorx"
	ClassGraphPaperYMinor = "minory"
)
