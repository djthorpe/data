package data

/////////////////////////////////////////////////////////////////////
// TYPES

// Color represents an 24-bit RGB colour without opacity
type Color struct {
	R, G, B uint8
}

type ColorSwatch uint32

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	ColorAll ColorSwatch = 0
	ColorRed ColorSwatch = (1 << iota)
	ColorOrange
	ColorYellow
	ColorGreen
	ColorCyan
	ColorBlue
	ColorPurple
	ColorPink
	ColorWhite
	ColorGray
	ColorBrown
	ColorPrimary
)

/////////////////////////////////////////////////////////////////////
// METHODS

// RGBA returns fully opaque color values
func (c Color) RGBA() (r, g, b, a uint32) {
	return uint32(c.R) * 257, uint32(c.G) * 257, uint32(c.B) * 257, 0xFFFF
}
