package color

import (
	"fmt"
	"strings"

	"github.com/djthorpe/data"
)

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
	Aqua                 = data.Color{0, 255, 255}   // #00FFFF
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
	colorNames = map[string]data.Color{
		"LightSalmon":          LightSalmon,
		"Salmon":               Salmon,
		"DarkSalmon":           DarkSalmon,
		"LightCoral":           LightCoral,
		"IndianRed":            IndianRed,
		"Crimson":              Crimson,
		"FireBrick":            FireBrick,
		"Red":                  Red,
		"DarkRed":              DarkRed,
		"Coral":                Coral,
		"Tomato":               Tomato,
		"OrangeRed":            OrangeRed,
		"Gold":                 Gold,
		"Orange":               Orange,
		"DarkOrange":           DarkOrange,
		"LightYellow":          LightYellow,
		"LemonChiffon":         LemonChiffon,
		"LightGoldenrodYellow": LightGoldenrodYellow,
		"PapayaWhip":           PapayaWhip,
		"Moccasin":             Moccasin,
		"PeachPuff":            PeachPuff,
		"PaleGoldenrod":        PaleGoldenrod,
		"Khaki":                Khaki,
		"DarkKhaki":            DarkKhaki,
		"Yellow":               Yellow,
		"LawnGreen":            LawnGreen,
		"Chartreuse":           Chartreuse,
		"LimeGreen":            LimeGreen,
		"Lime":                 Lime,
		"ForestGreen":          ForestGreen,
		"Green":                Green,
		"DarkGreen":            DarkGreen,
		"GreenYellow":          GreenYellow,
		"YellowGreen":          YellowGreen,
		"SpringGreen":          SpringGreen,
		"MediumSpringGreen":    MediumSpringGreen,
		"LightGreen":           LightGreen,
		"PaleGreen":            PaleGreen,
		"DarkSeaGreen":         DarkSeaGreen,
		"MediumSeaGreen":       MediumSeaGreen,
		"SeaGreen":             SeaGreen,
		"Olive":                Olive,
		"DarkOliveGreen":       DarkOliveGreen,
		"OliveDrab":            OliveDrab,
		"LightCyan":            LightCyan,
		"Cyan":                 Cyan,
		"Aqua":                 Aqua,
		"Aquamarine":           Aquamarine,
		"MediumAquamarine":     MediumAquamarine,
		"PaleTurquoise":        PaleTurquoise,
		"Turquoise":            Turquoise,
		"MediumTurquoise":      MediumTurquoise,
		"DarkTurquoise":        DarkTurquoise,
		"LightSeagreen":        LightSeagreen,
		"CadetBlue":            CadetBlue,
		"DarkCyan":             DarkCyan,
		"Teal":                 Teal,
		"PowderBlue":           PowderBlue,
		"LightBlue":            LightBlue,
		"LightSkyBlue":         LightSkyBlue,
		"SkyBlue":              SkyBlue,
		"DeepSkyBlue":          DeepSkyBlue,
		"LightSteelBlue":       LightSteelBlue,
		"DodgerBlue":           DodgerBlue,
		"CornflowerBlue":       CornflowerBlue,
		"SteelBlue":            SteelBlue,
		"RoyalBlue":            RoyalBlue,
		"Blue":                 Blue,
		"MediumBlue":           MediumBlue,
		"DarkBlue":             DarkBlue,
		"Navy":                 Navy,
		"MidnightBlue":         MidnightBlue,
		"MediumSlateBlue":      MediumSlateBlue,
		"SlateBlue":            SlateBlue,
		"DarkSlateBlue":        DarkSlateBlue,
		"Lavender":             Lavender,
		"Thistle":              Thistle,
		"Plum":                 Plum,
		"Violet":               Violet,
		"Orchid":               Orchid,
		"Fuchsia":              Fuchsia,
		"Magenta":              Magenta,
		"MediumOrchid":         MediumOrchid,
		"MediumPurple":         MediumPurple,
		"BlueViolet":           BlueViolet,
		"DarkViolet":           DarkViolet,
		"DarkOrchid":           DarkOrchid,
		"DarkMagenta":          DarkMagenta,
		"Purple":               Purple,
		"Indigo":               Indigo,
		"Pink":                 Pink,
		"LightPink":            LightPink,
		"HotPink":              HotPink,
		"DeepPink":             DeepPink,
		"PaleVioletRed":        PaleVioletRed,
		"MediumVioletRed":      MediumVioletRed,
		"White":                White,
		"Snow":                 Snow,
		"Honeydew":             Honeydew,
		"MintCream":            MintCream,
		"Azure":                Azure,
		"AliceBlue":            AliceBlue,
		"GhostWhite":           GhostWhite,
		"WhiteSmoke":           WhiteSmoke,
		"Seashell":             Seashell,
		"Beige":                Beige,
		"OldLace":              OldLace,
		"FloralWhite":          FloralWhite,
		"Ivory":                Ivory,
		"AntiqueWhite":         AntiqueWhite,
		"Linen":                Linen,
		"LavenderBlush":        LavenderBlush,
		"MistyRose":            MistyRose,
		"Gainsboro":            Gainsboro,
		"LightGray":            LightGray,
		"Silver":               Silver,
		"DarkGray":             DarkGray,
		"Gray":                 Gray,
		"DimGray":              DimGray,
		"LightSlateGray":       LightSlateGray,
		"SlateGray":            SlateGray,
		"DarkSlateGray":        DarkSlateGray,
		"Black":                Black,
		"Cornsilk":             Cornsilk,
		"BlancheDalmond":       BlancheDalmond,
		"Bisque":               Bisque,
		"NavajoWhite":          NavajoWhite,
		"Wheat":                Wheat,
		"Burlywood":            Burlywood,
		"Tan":                  Tan,
		"RosyBrown":            RosyBrown,
		"SandyBrown":           SandyBrown,
		"Goldenrod":            Goldenrod,
		"Peru":                 Peru,
		"Chocolate":            Chocolate,
		"SaddleBrown":          SaddleBrown,
		"Sienna":               Sienna,
		"Brown":                Brown,
		"Maroon":               Maroon,
	}
	colorHash = map[string]string{}
)

func init() {
	// Make a reverse lookup from hash -> colorname
	for name, value := range colorNames {
		key := HashString(value)
		colorHash[key] = name
	}
}

// Return hashed color
func HashString(color data.Color) string {
	return fmt.Sprint("#%02X%02X%02X", color.R, color.G, color.B)
}

// Return named colours, or hashed colours
func String(color data.Color) string {
	colorhash := HashString(color)
	if colorname, exists := colorHash[colorhash]; exists {
		return strings.ToLower(colorname)
	} else {
		return colorhash
	}
}
