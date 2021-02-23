package data

import "io"

/////////////////////////////////////////////////////////////////////
// INTERFACES

type Viz interface {
	// Create graph paper with major x minor grid squares
	GraphPaper(major, minor uint) VizGraphPaper

	// Create scale on either X or Y axis for a set
	RealScale(RealSet, Orientation) VizScale

	// Create table
	// TODO Table(Table, ...TableOpt) VizTable

	// Write to data stream
	Write(Writer, io.Writer) error
}

type VizGraphPaper interface{}
type VizTable interface{}
type VizScale interface{}

/////////////////////////////////////////////////////////////////////
// TYPES

type Orientation int

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

	// Class names used for Scales
	ClassScale = "scale"
)

const (
	Horizontal Orientation = (1 << iota)
	Vertical
)
