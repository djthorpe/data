---
description: 'Organization, Reading, Transformation and Rendering of 2D Graphics primitives'
---

# Canvas

You can create canvas objects by using the `NewCanvas` method and then create shapes, styles, transforms and groups on the canvas before writing out the canvas to render it. A subset of SVG [Scalable Vector Graphics](https://developer.mozilla.org/en-US/docs/Web/SVG) format is supported for representation and presently external renderers include most web browsers, but the intention is also to provide addition
renderers for displays, bitmap and print formats.

The following code writes out a canvas as SVG of A4 paper size:

```go
package main

import (
    "os"
    "github.com/djthorpe/data"
    "github.com/djthorpe/data/pkg/canvas"
    "github.com/djthorpe/data/pkg/color"
)

func main() {
    c := canvas.NewCanvas(data.A4LandscapeSize, data.MM)

    c.Title("Canvas Document")
    c.Desc("SVG Document output")
    c.Rect(c.Origin(), c.Size()).Style(
        c.Fill(color.Ivory, 1.0),
    )
    c.Write(data.SVG, os.Stdout)
}
```

A canvas can consist of:

* `data.CanvasElement`, which are drawing elements like circles, text, lines, paths and text;
* `data.CanvasGroup`, which are groups of elements.

Any element or group can have **transforms** applied \(for example, scaling, rotating and translation\) and **style** \(for example, color or opacity\).

## Creating a Canvas

Canvases are created with the required size. There are several predefined sizes:

* `canvas.NewCanvas(data.A4LandscapeSize, data.MM)` defines an A4 canvas in landscape;
* `canvas.NewCanvas(data.A4PortraitSize, data.MM)` defines an A4 canvas in portrait;
* `data.LetterPortraitSize` and `data.LetterLandscapeSize` can be used as the first argument for Letter sized paper;
* `data.LegalPortraitSize` and `LegalLandscapeSize` can be used to define a canvas the size of Legal paper.

The second argument are the units:

| Constant | Unit | Description |
| :--- | :--- | :--- |
| `data.None` | Natural | Based on screen size |
| `data.PX` | Pixels | Based on screen size |
| `data.CM` | Centimetre \(cm\) |  |
| `data.MM` | Millimetre \(mm\) | 10mm = 1cm |
| `data.IN` | Inch | 1in = 2.54cm |
| `data.PC` | Picas | 6pc = 1 in |
| `data.PT` | Points | 72pt ≅ 1 in |
| `data.EM` |  | Equivalent to the computed font-size |
| `data.EX` |  | Equivalent to the height of a lower-case letter |

In addition to setting canvas size, it's possible to get and set the translation to natural units used to render elements, using the view box interface:

```go
type Canvas interface {
    ViewBox() (Point,Size)
    SetViewBox(origin Point,size Size)
    // ...
}
```

The `SetViewBox` method adjusts the top left co-ordinate to be equal to the **origin** and the width and height to be equal to the **size** where natural units are used. For example,

```go
    c := canvas.NewCanvas(data.LetterPortraitSize, data.MM)
    c.SetViewBox(data.ZeroPoint,geom.DivideSize(c.Size(), 0.3527))
```

This example creates a canvas of 612 units across and 792 units high, mapped onto a canvas which would fit on letter paper: `0.3527` is approximately 2.84 natural units per millimetre, which is 72 units per inch.

## Adding Shape Elements to the Canvas

Shape elements are created using the following canvas methods:

```go
type Canvas interface {
    Circle(centre Point,radius float32) CanvasElement
	Ellipse(centre Point,radius Size) CanvasElement
	Line(p1 Point,p2 Point) CanvasElement
	Rect(origin Point,size Size) CanvasElement
	Path(segments ...CanvasPath) CanvasElement
	Polyline(points ...Point) CanvasElement
	Polygon(points ...Point) CanvasElement
	Text(origin Point,rel bool,segments ...CanvasText) CanvasElement
	Image(origin Point,size Size,url string) CanvasElement
    // ...
}
```

For example calling `c.Circle(data.Point{ 0,0 },10 })` creates a circle with centre of top left origin and with radius 10. Some elements require one or more additional arguments to construct the shape:

  * `Path` requires one or more `CanvasPath` segment instructions, which can append a line or curve to the path. See below for some examples of constructing paths;
  * `Polyline` and `Polygon` require one or more points to define the shape;
  * `Text` requires one or more `TextSpan` or `TextPath` elements. See below for some examples of contructing text primitives.

## Styling

Canvas elements can be styled visually with one or more style declarations, which are arguments to the `element.Style` function. These declarations are grouped into fill, stroke, text and other. For example,

```go
package main

import (
    "os"
    "github.com/djthorpe/data"
    "github.com/djthorpe/data/pkg/canvas"
    "github.com/djthorpe/data/pkg/color"
)

func main() {
    c := canvas.NewCanvas(data.A4LandscapeSize, data.MM)
    c.Rect(c.Origin(), c.Size()).Style(
        c.Fill(color.LightGray, 1.0),
        c.Stroke(color.Black, 1.0),
        c.StrokeWidth(2.0),
    )
    c.Text(geom.CentrePoint(c.Origin(),c.Size()),false,
        c.TextSpan("Hello, world!"),
    ).Style(
        c.TextAnchor(data.Middle),
        c.FontSize(64.0,data.PT),
    )
    c.Write(data.SVG, os.Stdout)
}
```

This outputs SVG which prints "Hello, World" in 64pt at the centre of an A4 rectangle filled in light gray and with a 2mm black border.

### Fill Style Declarations

| Declaration | Arguments | Description |
| :--- | :--- | :--- |
| `canvas.NoFill` | | Do not fill the element |
| `canvas.Fill` | `color data.Color, opacity float32` | Fill the element with color and with opacity between 0.0 and 1.0 |
| `canvas.FillRule` | `data.NonZero \| data.EvenOdd` | When EvenOdd, fill crossing segments alternatively. See [here](https://developer.mozilla.org/en-US/docs/Web/SVG/Attribute/fill-rule) for more information. |

### Stroke Style Declarations

| Declaration | Arguments | Description |
| :--- | :--- | :--- |
| `canvas.NoStroke` | | Do not outline the element |
| `canvas.Stroke` | `color data.Color, opacity float32` | Draw outline with color and with opacity between 0.0 and 1.0 |
| `canvas.StrokeWidth` | `width float32` | Width of element with user unit width
| `canvas.LineCap` | `data.CapButt \| data.CapRound \| data.CapSquare` | Line endings. See [here](https://developer.mozilla.org/en-US/docs/Web/SVG/Attribute/stroke-linecap) for more information.
| `canvas.LineJoin` | `data.JoinMiter \| data.JoinMiterClip \| data.JoinArcs \| data.JoinRound \| data.JoinBevel` | Line join. See [here](https://developer.mozilla.org/en-US/docs/Web/SVG/Attribute/stroke-linejoin) for more information.
| `canvas.MiterLimit` | `ratio float32` | Miter limit when joining two lines. See [here](https://developer.mozilla.org/en-US/docs/Web/SVG/Attribute/stroke-miterlimit) for more information. |	

### Text Style Declarations

| Declaration | Arguments | Description |
| :--- | :--- | :--- |
| `canvas.FontSize` | `size float32, unit data.Unit` | Text rendering size. |
| `canvas.FontFamily` | `string` | Comma separated lists of font names.
| `canvas.FontVariant` | `data.Thin \| data.ExtraLight \| data.Light \| data.Regular \| data.Medium \| data.SemiBold \| data.Bold \| data.ExtraBold \| data.Black \| data.Bolder \| data.Lighter \| data.Italic \| data.Oblique` | One or more font variants (more than one can be specified using the OR (\|) operator).
| `canvas.TextAnchor` | `data.Start \| data.Middle \| data.End` | Anchor the origin to the start, middle or end of the text.



### Other Style Declarations

| Declaration | Arguments | Description |
| :--- | :--- | :--- |
| `canvas.UseMarker` |  `data.Start \| data.Middle \| data.End, id string` | When drawing line segments, use a specific marker referenced by *id* for the start of a series of segments, for the joins (middle) and/or end. A marker is defined by the `canvas.Marker` element and referenced by *id*. A marker can be used for more than one position using the OR (\|) operator.


## Grouping, Markers & Definitions

There are canvas declarations to group elements:

  * `canvas.Group` can group elements so that style and transformations can be applied to a set of elements;
  * `canvas.Marker` can group elements to define line start, middle and end markers;
  * `canvas.Defs` can define a set of elements to be used repeatedly on the canvas.
   
In general groups of elements are referred to using their `id` for later use. For example,

TODO

## Transformation

Elements and groups of elements can be transformed with one or more transformation declarations, which are arguments to the `element.Transform` function. Typically a transformation is a rotation, skew, scale or co-ordinate translation. Transformations usually occur one after another. For example,

```go
package main

import (
    "os"
    "github.com/djthorpe/data"
    "github.com/djthorpe/data/pkg/canvas"
    "github.com/djthorpe/data/pkg/color"
)

func main() {
    c := canvas.NewCanvas(data.A4LandscapeSize, data.MM)
    c.Rect(c.Origin(), c.Size()).Style(
        c.Fill(color.LightGray, 1.0),
    ).Transform(
        c.Rotate(45.0),
        c.Translate(geom.CentrePoint(c.Origin(),c.Size())),
    )
    c.Write(data.SVG, os.Stdout)
}
```

This will render a filled rectangle with its' origin at the centre of the canvas rotated by 45 degrees clockwise. The transformation declarations you can use are:

| Declaration | Arguments | Description |
| :--- | :--- | :--- |
| `canvas.Scale` | `data.Size` | Scale co-ordinate system by width and height. A value below 1.0 reduces the size, and above increases the size |
| `canvas.Translate` | `data.Point` | Move co-ordinate system by x and y. |
| `canvas.Rotate` | `degrees float32` | Rotate co-ordinate system around the zero point. A positive value rotates clockwise.  |
| `canvas.RotateAround` | `degrees float32, centre Point` | Rotate co-ordinate system around a centre point. This is equivalent to translating by `-centre`, rotation and then translating by `+centre`. |
| `canvas.SkewX` | `degrees float32` | Skew the X-Coordinates. |
| `canvas.SkewY` | `degrees float32` | Skew the Y-Coordinates. |

For more information on co-ordinate transformation please see [here](https://developer.mozilla.org/en-US/docs/Web/SVG/Attribute/transform).

## Rendering

TODO

## Limitations

TODO

## What's Next

TODO


