package color

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// TYPES

// SwatchColor defines a color and any palettes it belongs to
type SwatchColor struct {
	data.Color
	order int
}

// Swatch is an array of colors
type Swatch []data.Color

/////////////////////////////////////////////////////////////////////
// CONSTANTS

var (
	// https://www.rapidtables.com/web/css/css-color.html
	LightSalmon          = data.Color{255, 160, 122} // #FFA07A // Red palette
	Salmon               = data.Color{250, 128, 114} // #FA8072
	DarkSalmon           = data.Color{233, 150, 122} // #E9967A
	LightCoral           = data.Color{240, 128, 128} // #F08080
	IndianRed            = data.Color{205, 92, 92}   // #CD5C5C
	Crimson              = data.Color{220, 20, 60}   // #DC143C
	FireBrick            = data.Color{178, 34, 34}   // #B22222
	Red                  = data.Color{255, 0, 0}     // #FF0000
	DarkRed              = data.Color{139, 0, 0}     // #8B0000
	Coral                = data.Color{255, 127, 80}  // #FF7F50 // Orange palette
	Tomato               = data.Color{255, 99, 71}   // #FF6347
	OrangeRed            = data.Color{255, 69, 0}    // #FF4500
	Gold                 = data.Color{255, 215, 0}   // #FFD700
	Orange               = data.Color{255, 165, 0}   // #FFA500
	DarkOrange           = data.Color{255, 140, 0}   // #FF8C00
	LightYellow          = data.Color{255, 255, 224} // #FFFFE0 // Yellow palette
	LemonChiffon         = data.Color{255, 250, 205} // #FFFACD
	LightGoldenrodYellow = data.Color{250, 250, 210} // #FAFAD2
	PapayaWhip           = data.Color{255, 239, 213} // #FFEFD5
	Moccasin             = data.Color{255, 228, 181} // #FFE4B5
	PeachPuff            = data.Color{255, 218, 185} // #FFDAB9
	PaleGoldenrod        = data.Color{238, 232, 170} // #EEE8AA
	Khaki                = data.Color{240, 230, 140} // #F0E68C
	DarkKhaki            = data.Color{189, 183, 107} // #BDB76B
	Yellow               = data.Color{255, 255, 0}   // #FFFF00
	LawnGreen            = data.Color{124, 252, 0}   // #7CFC00 // Green palette
	Chartreuse           = data.Color{127, 255, 0}   // #7FFF00
	LimeGreen            = data.Color{50, 205, 50}   // #32CD32
	Lime                 = data.Color{0, 255, 0}     // #00FF00
	ForestGreen          = data.Color{34, 139, 34}   // #228B22
	Green                = data.Color{0, 128, 0}     // #008000
	DarkGreen            = data.Color{0, 100, 0}     // #006400
	GreenYellow          = data.Color{173, 255, 47}  // #ADFF2F
	YellowGreen          = data.Color{154, 205, 50}  // #9ACD32
	SpringGreen          = data.Color{0, 255, 127}   // #00FF7F
	MediumSpringGreen    = data.Color{0, 250, 154}   // #00FA9A
	LightGreen           = data.Color{144, 238, 144} // #90EE90
	PaleGreen            = data.Color{152, 251, 152} // #98FB98
	DarkSeaGreen         = data.Color{143, 188, 143} // #8FBC8F
	MediumSeaGreen       = data.Color{60, 179, 113}  // #3CB371
	SeaGreen             = data.Color{46, 139, 87}   // #2E8B57
	Olive                = data.Color{128, 128, 0}   // #808000
	DarkOliveGreen       = data.Color{85, 107, 47}   // #556B2F
	OliveDrab            = data.Color{107, 142, 35}  // #6B8E23
	LightCyan            = data.Color{224, 255, 255} // #E0FFFF // Cyan palette
	Cyan                 = data.Color{0, 255, 255}   // #00FFFF
	Aquamarine           = data.Color{127, 255, 212} // #7FFFD4
	MediumAquamarine     = data.Color{102, 205, 170} // #66CDAA
	PaleTurquoise        = data.Color{175, 238, 238} // #AFEEEE
	Turquoise            = data.Color{64, 224, 208}  // #40E0D0
	MediumTurquoise      = data.Color{72, 209, 204}  // #48D1CC
	DarkTurquoise        = data.Color{0, 206, 209}   // #00CED1
	LightSeagreen        = data.Color{32, 178, 170}  // #20B2AA
	CadetBlue            = data.Color{95, 158, 160}  // #5F9EA0
	DarkCyan             = data.Color{0, 139, 139}   // #008B8B
	Teal                 = data.Color{0, 128, 128}   // #008080
	PowderBlue           = data.Color{176, 224, 230} // #B0E0E6 // Blue palette
	LightBlue            = data.Color{173, 216, 230} // #ADD8E6
	LightSkyBlue         = data.Color{135, 206, 250} // #87CEFA
	SkyBlue              = data.Color{135, 206, 235} // #87CEEB
	DeepSkyBlue          = data.Color{0, 191, 255}   // #00BFFF
	LightSteelBlue       = data.Color{176, 196, 222} // #B0C4DE
	DodgerBlue           = data.Color{30, 144, 255}  // #1E90FF
	CornflowerBlue       = data.Color{100, 149, 237} // #6495ED
	SteelBlue            = data.Color{70, 130, 180}  // #4682B4
	RoyalBlue            = data.Color{65, 105, 225}  // #4169E1
	Blue                 = data.Color{0, 0, 255}     // #0000FF
	MediumBlue           = data.Color{0, 0, 205}     // #0000CD
	DarkBlue             = data.Color{0, 0, 139}     // #00008B
	Navy                 = data.Color{0, 0, 128}     // #000080
	MidnightBlue         = data.Color{25, 25, 112}   // #191970
	MediumSlateBlue      = data.Color{123, 104, 238} // #7B68EE
	SlateBlue            = data.Color{106, 90, 205}  // #6A5ACD
	DarkSlateBlue        = data.Color{72, 61, 139}   // #483D8B
	Lavender             = data.Color{230, 230, 250} // #E6E6FA // Purple palette
	Thistle              = data.Color{216, 191, 216} // #D8BFD8
	Plum                 = data.Color{221, 160, 221} // #DDA0DD
	Violet               = data.Color{238, 130, 238} // #EE82EE
	Orchid               = data.Color{218, 112, 214} // #DA70D6
	Fuchsia              = data.Color{255, 0, 255}   // #FF00FF
	Magenta              = data.Color{255, 0, 255}   // #FF00FF
	MediumOrchid         = data.Color{186, 85, 211}  // #BA55D3
	MediumPurple         = data.Color{147, 112, 219} // #9370DB
	BlueViolet           = data.Color{138, 43, 226}  // #8A2BE2
	DarkViolet           = data.Color{148, 0, 211}   // #9400D3
	DarkOrchid           = data.Color{153, 50, 204}  // #9932CC
	DarkMagenta          = data.Color{139, 0, 139}   // #8B008B
	Purple               = data.Color{128, 0, 128}   // #800080
	Indigo               = data.Color{75, 0, 130}    // #4B0082
	Pink                 = data.Color{255, 192, 203} // #FFC0CB // Pink palette
	LightPink            = data.Color{255, 182, 193} // #FFB6C1
	HotPink              = data.Color{255, 105, 180} // #FF69B4
	DeepPink             = data.Color{255, 20, 147}  // #FF1493
	PaleVioletRed        = data.Color{219, 112, 147} // #DB7093
	MediumVioletRed      = data.Color{199, 21, 133}  // #C71585
	White                = data.Color{255, 255, 255} // #FFFFFF // White palette
	Snow                 = data.Color{255, 250, 250} // #FFFAFA
	Honeydew             = data.Color{240, 255, 240} // #F0FFF0
	MintCream            = data.Color{245, 255, 250} // #F5FFFA
	Azure                = data.Color{240, 255, 255} // #F0FFFF
	AliceBlue            = data.Color{240, 248, 255} // #F0F8FF
	GhostWhite           = data.Color{248, 248, 255} // #F8F8FF
	WhiteSmoke           = data.Color{245, 245, 245} // #F5F5F5
	Seashell             = data.Color{255, 245, 238} // #FFF5EE
	Beige                = data.Color{245, 245, 220} // #F5F5DC
	OldLace              = data.Color{253, 245, 230} // #FDF5E6
	FloralWhite          = data.Color{255, 250, 240} // #FFFAF0
	Ivory                = data.Color{255, 255, 240} // #FFFFF0
	AntiqueWhite         = data.Color{250, 235, 215} // #FAEBD7
	Linen                = data.Color{250, 240, 230} // #FAF0E6
	LavenderBlush        = data.Color{255, 240, 245} // #FFF0F5
	MistyRose            = data.Color{255, 228, 225} // #FFE4E1
	Gainsboro            = data.Color{220, 220, 220} // #DCDCDC // Grey palette (inc. black)
	LightGray            = data.Color{211, 211, 211} // #D3D3D3
	Silver               = data.Color{192, 192, 192} // #C0C0C0
	DarkGray             = data.Color{169, 169, 169} // #A9A9A9
	Gray                 = data.Color{128, 128, 128} // #808080
	DimGray              = data.Color{105, 105, 105} // #696969
	LightSlateGray       = data.Color{119, 136, 153} // #778899
	SlateGray            = data.Color{112, 128, 144} // #708090
	DarkSlateGray        = data.Color{47, 79, 79}    // #2F4F4F
	Black                = data.Color{0, 0, 0}       // #000000
	Cornsilk             = data.Color{255, 248, 220} // #FFF8DC // Brown palette
	BlancheDalmond       = data.Color{255, 235, 205} // #FFEBCD
	Bisque               = data.Color{255, 228, 196} // #FFE4C4
	NavajoWhite          = data.Color{255, 222, 173} // #FFDEAD
	Wheat                = data.Color{245, 222, 179} // #F5DEB3
	Burlywood            = data.Color{222, 184, 135} // #DEB887
	Tan                  = data.Color{210, 180, 140} // #D2B48C
	RosyBrown            = data.Color{188, 143, 143} // #BC8F8F
	SandyBrown           = data.Color{244, 164, 96}  // #F4A460
	Goldenrod            = data.Color{218, 165, 32}  // #DAA520
	Peru                 = data.Color{205, 133, 63}  // #CD853F
	Chocolate            = data.Color{210, 105, 30}  // #D2691E
	SaddleBrown          = data.Color{139, 69, 19}   // #8B4513
	Sienna               = data.Color{160, 82, 45}   // #A0522D
	Brown                = data.Color{165, 42, 42}   // #A52A2A
	Maroon               = data.Color{128, 0, 0}     // #800000
)

var (
	colorNames = map[string]SwatchColor{
		"LightSalmon":          SwatchColor{LightSalmon, 1},
		"Salmon":               SwatchColor{Salmon, 2},
		"DarkSalmon":           SwatchColor{DarkSalmon, 3},
		"LightCoral":           SwatchColor{LightCoral, 4},
		"IndianRed":            SwatchColor{IndianRed, 5},
		"Crimson":              SwatchColor{Crimson, 6},
		"FireBrick":            SwatchColor{FireBrick, 7},
		"Red":                  SwatchColor{Red, 8},
		"DarkRed":              SwatchColor{DarkRed, 9},
		"Coral":                SwatchColor{Coral, 10},
		"Tomato":               SwatchColor{Tomato, 11},
		"OrangeRed":            SwatchColor{OrangeRed, 12},
		"Gold":                 SwatchColor{Gold, 13},
		"Orange":               SwatchColor{Orange, 14},
		"DarkOrange":           SwatchColor{DarkOrange, 15},
		"LightYellow":          SwatchColor{LightYellow, 16},
		"LemonChiffon":         SwatchColor{LemonChiffon, 17},
		"LightGoldenrodYellow": SwatchColor{LightGoldenrodYellow, 18},
		"PapayaWhip":           SwatchColor{PapayaWhip, 19},
		"Moccasin":             SwatchColor{Moccasin, 20},
		"PeachPuff":            SwatchColor{PeachPuff, 21},
		"PaleGoldenrod":        SwatchColor{PaleGoldenrod, 22},
		"Khaki":                SwatchColor{Khaki, 23},
		"DarkKhaki":            SwatchColor{DarkKhaki, 24},
		"Yellow":               SwatchColor{Yellow, 25},
		"LawnGreen":            SwatchColor{LawnGreen, 26},
		"Chartreuse":           SwatchColor{Chartreuse, 27},
		"LimeGreen":            SwatchColor{LimeGreen, 28},
		"Lime":                 SwatchColor{Lime, 29},
		"ForestGreen":          SwatchColor{ForestGreen, 30},
		"Green":                SwatchColor{Green, 31},
		"DarkGreen":            SwatchColor{DarkGreen, 32},
		"GreenYellow":          SwatchColor{GreenYellow, 33},
		"YellowGreen":          SwatchColor{YellowGreen, 34},
		"SpringGreen":          SwatchColor{SpringGreen, 35},
		"MediumSpringGreen":    SwatchColor{MediumSpringGreen, 36},
		"LightGreen":           SwatchColor{LightGreen, 37},
		"PaleGreen":            SwatchColor{PaleGreen, 38},
		"DarkSeaGreen":         SwatchColor{DarkSeaGreen, 39},
		"MediumSeaGreen":       SwatchColor{MediumSeaGreen, 40},
		"SeaGreen":             SwatchColor{SeaGreen, 41},
		"Olive":                SwatchColor{Olive, 42},
		"DarkOliveGreen":       SwatchColor{DarkOliveGreen, 43},
		"OliveDrab":            SwatchColor{OliveDrab, 44},
		"LightCyan":            SwatchColor{LightCyan, 45},
		"Cyan":                 SwatchColor{Cyan, 46},
		"Aquamarine":           SwatchColor{Aquamarine, 48},
		"MediumAquamarine":     SwatchColor{MediumAquamarine, 49},
		"PaleTurquoise":        SwatchColor{PaleTurquoise, 50},
		"Turquoise":            SwatchColor{Turquoise, 51},
		"MediumTurquoise":      SwatchColor{MediumTurquoise, 52},
		"DarkTurquoise":        SwatchColor{DarkTurquoise, 53},
		"LightSeagreen":        SwatchColor{LightSeagreen, 54},
		"CadetBlue":            SwatchColor{CadetBlue, 55},
		"DarkCyan":             SwatchColor{DarkCyan, 56},
		"Teal":                 SwatchColor{Teal, 57},
		"PowderBlue":           SwatchColor{PowderBlue, 58},
		"LightBlue":            SwatchColor{LightBlue, 59},
		"LightSkyBlue":         SwatchColor{LightSkyBlue, 60},
		"SkyBlue":              SwatchColor{SkyBlue, 61},
		"DeepSkyBlue":          SwatchColor{DeepSkyBlue, 62},
		"LightSteelBlue":       SwatchColor{LightSteelBlue, 63},
		"DodgerBlue":           SwatchColor{DodgerBlue, 64},
		"CornflowerBlue":       SwatchColor{CornflowerBlue, 65},
		"SteelBlue":            SwatchColor{SteelBlue, 66},
		"RoyalBlue":            SwatchColor{RoyalBlue, 67},
		"Blue":                 SwatchColor{Blue, 68},
		"MediumBlue":           SwatchColor{MediumBlue, 69},
		"DarkBlue":             SwatchColor{DarkBlue, 70},
		"Navy":                 SwatchColor{Navy, 71},
		"MidnightBlue":         SwatchColor{MidnightBlue, 72},
		"MediumSlateBlue":      SwatchColor{MediumSlateBlue, 73},
		"SlateBlue":            SwatchColor{SlateBlue, 74},
		"DarkSlateBlue":        SwatchColor{DarkSlateBlue, 75},
		"Lavender":             SwatchColor{Lavender, 76},
		"Thistle":              SwatchColor{Thistle, 77},
		"Plum":                 SwatchColor{Plum, 78},
		"Violet":               SwatchColor{Violet, 79},
		"Orchid":               SwatchColor{Orchid, 80},
		"Fuchsia":              SwatchColor{Fuchsia, 81},
		"Magenta":              SwatchColor{Magenta, 82},
		"MediumOrchid":         SwatchColor{MediumOrchid, 83},
		"MediumPurple":         SwatchColor{MediumPurple, 84},
		"BlueViolet":           SwatchColor{BlueViolet, 85},
		"DarkViolet":           SwatchColor{DarkViolet, 86},
		"DarkOrchid":           SwatchColor{DarkOrchid, 87},
		"DarkMagenta":          SwatchColor{DarkMagenta, 88},
		"Purple":               SwatchColor{Purple, 89},
		"Indigo":               SwatchColor{Indigo, 90},
		"Pink":                 SwatchColor{Pink, 91},
		"LightPink":            SwatchColor{LightPink, 92},
		"HotPink":              SwatchColor{HotPink, 93},
		"DeepPink":             SwatchColor{DeepPink, 94},
		"PaleVioletRed":        SwatchColor{PaleVioletRed, 95},
		"MediumVioletRed":      SwatchColor{MediumVioletRed, 96},
		"White":                SwatchColor{White, 97},
		"Snow":                 SwatchColor{Snow, 98},
		"Honeydew":             SwatchColor{Honeydew, 99},
		"MintCream":            SwatchColor{MintCream, 100},
		"Azure":                SwatchColor{Azure, 101},
		"AliceBlue":            SwatchColor{AliceBlue, 102},
		"GhostWhite":           SwatchColor{GhostWhite, 103},
		"WhiteSmoke":           SwatchColor{WhiteSmoke, 104},
		"Seashell":             SwatchColor{Seashell, 105},
		"Beige":                SwatchColor{Beige, 106},
		"OldLace":              SwatchColor{OldLace, 107},
		"FloralWhite":          SwatchColor{FloralWhite, 108},
		"Ivory":                SwatchColor{Ivory, 109},
		"AntiqueWhite":         SwatchColor{AntiqueWhite, 110},
		"Linen":                SwatchColor{Linen, 111},
		"LavenderBlush":        SwatchColor{LavenderBlush, 112},
		"MistyRose":            SwatchColor{MistyRose, 113},
		"Gainsboro":            SwatchColor{Gainsboro, 114},
		"LightGray":            SwatchColor{LightGray, 115},
		"Silver":               SwatchColor{Silver, 116},
		"DarkGray":             SwatchColor{DarkGray, 117},
		"Gray":                 SwatchColor{Gray, 118},
		"DimGray":              SwatchColor{DimGray, 119},
		"LightSlateGray":       SwatchColor{LightSlateGray, 120},
		"SlateGray":            SwatchColor{SlateGray, 121},
		"DarkSlateGray":        SwatchColor{DarkSlateGray, 122},
		"Black":                SwatchColor{Black, 123},
		"Cornsilk":             SwatchColor{Cornsilk, 124},
		"BlancheDalmond":       SwatchColor{BlancheDalmond, 125},
		"Bisque":               SwatchColor{Bisque, 126},
		"NavajoWhite":          SwatchColor{NavajoWhite, 127},
		"Wheat":                SwatchColor{Wheat, 128},
		"Burlywood":            SwatchColor{Burlywood, 129},
		"Tan":                  SwatchColor{Tan, 130},
		"RosyBrown":            SwatchColor{RosyBrown, 131},
		"SandyBrown":           SwatchColor{SandyBrown, 132},
		"Goldenrod":            SwatchColor{Goldenrod, 133},
		"Peru":                 SwatchColor{Peru, 134},
		"Chocolate":            SwatchColor{Chocolate, 135},
		"SaddleBrown":          SwatchColor{SaddleBrown, 136},
		"Sienna":               SwatchColor{Sienna, 137},
		"Brown":                SwatchColor{Brown, 138},
		"Maroon":               SwatchColor{Maroon, 139},
	}
	colorSync sync.Once
	colorHash = make(map[string]string, len(colorNames))
)

var (
	reWords = regexp.MustCompile("[A-Z][^A-Z]*")
)

/////////////////////////////////////////////////////////////////////
// METHODS

func colorInit() {
	colorSync.Do(func() {
		// Make a reverse lookup from hash -> colorname
		for name, value := range colorNames {
			key := HashString(value.Color)
			colorHash[key] = name
		}
	})
}

// HashString returns an RGB hash for a color
func HashString(color data.Color) string {
	return fmt.Sprintf("#%02X%02X%02X", color.R, color.G, color.B)
}

// String returns a named color or a hash otherwise
func String(color data.Color) string {
	colorInit()
	colorhash := HashString(color)
	if colorname, exists := colorHash[colorhash]; exists {
		return strings.ToLower(colorname)
	} else {
		return colorhash
	}
}

// Name returns the name of a color or hash otherwise
func Name(color data.Color) string {
	colorInit()
	colorhash := HashString(color)
	if colorname, exists := colorHash[colorhash]; exists == false {
		return colorhash
	} else if submatch := reWords.FindAllString(colorname, -1); submatch != nil {
		return strings.Join(submatch, " ")
	} else {
		return colorname
	}
}

// Palette returns all colors in palette
func Palette() []data.Color {
	// TODO: Take an argument later on with different swatches
	colorInit()
	colors := make(Swatch, 0, len(colorNames))
	for _, value := range colorNames {
		colors = append(colors, value.Color)
	}
	// Sort colors
	sort.Sort(colors)
	// Return colors
	return colors
}

// Distance returns the distance between two colors
func Distance(x, y data.Color) float32 {
	// Ref: https://www.compuphase.com/cmetric.htm
	rmean := float64(x.R) + float64(y.R)/2
	r := float64(x.R) - float64(y.R)
	g := float64(x.G) - float64(y.G)
	b := float64(x.B) - float64(y.B)
	return float32(math.Sqrt((((512 + rmean) * r * r) / 256) + 4*g*g + (((767 - rmean) * b * b) / 256)))
}

/////////////////////////////////////////////////////////////////////
// SWATCH METHODS

func (arr Swatch) Len() int {
	return len(arr)
}

func (arr Swatch) Less(i, j int) bool {
	return colorSwatchOrder(arr[i]) < colorSwatchOrder(arr[j])
}

func colorSwatchOrder(color data.Color) int {
	if colorname, exists := colorHash[HashString(color)]; exists == false {
		return 0
	} else if colorswatch, exists := colorNames[colorname]; exists == false {
		return 0
	} else {
		return colorswatch.order
	}
}

func (arr Swatch) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}