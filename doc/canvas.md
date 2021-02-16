# Canvas

You can create canvas objects by using the `NewCanvas` method and then create shape primitives, styles, transforms and groups on the canvas before writing out the canvas to render it. Currently SVG is supported for rendering. For example, the following code writes out a canvas as SVG of A4 paper size:

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

Canvas can consist of:

* `data.CanvasElement`, which are primitive drawing elements like circles, text, lines, paths and text;
* `data.CanvasGroup`, which are groups of elements.

Any element or group can have **transforms** applied \(for example, scaling, rotating and translation\) and **style** \(for example, color or opacity\).

## Creating a Canvas

Canvases are created with the size of the canvas. There are several predefined sizes:

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

In addition to setting canvas size, it's possible to set the translation to natural units used to render elements, using `Size`, `SetSize`, `Origin` and `SetOrigin`:

```go
type Canvas interface {
    Origin() Point
    Size() Size
    SetOrigin(Point)
    SetSize(Size)
    // ...
}
```

This adjusts the top left co-ordinate to be equal to the **origin** and the width and height to be equal to the **size** where natural units are used. For example,

```go
    c := canvas.NewCanvas(data.LetterPortraitSize, data.MM)
    c.SetSize(geom.DivideSize(c.Size(), 0.3527))
```

This example creates a canvas of 612 units across and 792 units high, mapped onto a canvas which would fit on letter paper: `0.3527` is approximately 2.84 natural units per millimetre, which is 72 units per inch.

## Adding Shapes to the Canvas

Shape primitives are created using the following canvas methods:

```go
type Canvas interface {
    Circle(Point, float32) CanvasElement
    Ellipse(Point, Size) CanvasElement
    Line(Point, Point) CanvasElement
    Rect(Point, Size) CanvasElement
    Text(Point, ...CanvasText) CanvasElement
    Path([]Point) CanvasPath
    // ...
}
```

For example calling `c.Circle(data.Point{ 0,0 },10 })` creates a circle at the top left origin of natural radius 10, and adds the circle to the canvas. The arguments for each method are:

| Method | Arguments | Description |
| :--- | :--- | :--- |
| `data.Circle` | `{x-origin,y-origin},radius` | Circle |
| `data.Ellipse` | `{x-origin,y-origin},{x-radius,y-radius}` | Ellipse |
| `data.Line` | `{x-start,y-start},{x-end,y-end}` | Line |
| `data.Rect` | `{x-start,y-start},{x-end,y-end}` | Rectangle |
| `data.Text` | `{x,y},...string` | Adds text around point |
| `data.Path` | `{x-move,y-move},...{x-line,{y-line}` | Complex path. See below for bezier curves |

TODO

## Styling

TODO

## Transformation

TODO

## Grouping

TODO

## Rendering

TODO

