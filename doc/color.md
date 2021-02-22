---
description: 'Colour operations'
---

# Colour

Colours in the package are represented in 24 bits (8 bits for each plane R, G and B). In order to construct a colour use a defined color from the `color` package. For example,

```go
package main

import (
    "os"

    "github.com/djthorpe/data"
    "github.com/djthorpe/data/pkg/color"
)

func main() {
    reds := color.Palette(data.ColorRed)
    for _,c := range reds {
        fmt.Println("Red %q => %s",color.Name(c),color.Hash(c))
    }
}
```

In this example all predefined "Red" colors are returned by the `Palette` function and then the name and hash of the color is output.

## The `data.Color` type

TODO

```go
// HashString returns an RGB hash for a color
func HashString(color data.Color) string {
// String returns a named color or a hash otherwise
func String(color data.Color) string {
// Name returns the name of a color or hash otherwise
func Name(color data.Color) string {
// Palette returns colors in palette which adhere to a given
// set of swatches
func Palette(data.ColorSwatch) []data.Color {
// Distance returns the distance between two colors
func Distance(x, y data.Color) float32 {
// Nearest returns a color from the provided palette which is
// nearest to the provided color. If the palette is nil then
// all colors in the palette
func Nearest(c data.Color, palette []data.Color) data.Color {
// CMYK converts a color to a CMYK quadruple
func CMYK(c data.Color) color.CMYK {
// YCbCr converts an RGB triple to a Y'CbCr triple
func YCbCr(c data.Color) color.YCbCr {
```
