
# Data Extraction, Transformation, Processing and Visualisation

This repository contains various data extraction, transformation processing and visualization tools in __golang__. Currently it contains the following:

* [`data.Table`](doc/table.md) provides you with a way to ingest, transform and process data tables in comma-separated value format and output in CSV, ASCII and SQL formats;
* [`data.DOM`](doc/dom.md) provides a document object model which can read and write the XML format in addition to validating

  the XML;

* [`data.Canvas`](doc/canvas.md) provides a drawing canvas on which graphics primitives such as lines, circles, text and rectangles can be placed. Additionally transformation, grouping and stylizing of primitives can be applied. Canvases can currently be written in SVG format, the intention is to also allow rendering using OpenGL later.

## Documentation

I have published the documention at [data.mutablelogic.com](https://data.mutablelogic.com). You can also see the following useful sources of information:

* [pkg.go.dev documentation](https://pkg.go.dev/github.com/djthorpe/data)

## Usage

* Requires go version [1.13 and newer](https://golang.org/dl/);
* To use the Makefile, requires [GNU Make](https://www.gnu.org/software/make/).

There are various examples in the `cmd` folder. In order to build the examples, use the following command:

```bash
make
cd build/cmd
```

A temporary build folder is created on build. To run the tests or clean, use `make test` and `make clean` respectively. There is more information about the examples in the [documentation](doc/examples.md).

## Project Status

This module is currently **in development** but is mostly feature-complete.

## Contributing and Filing Issues

* [File an issue or question](http://github.com/djthorpe/graph/issues) on github.
* Feel free to fork this repository. Any pull requests are gratefully received. Licensed under Apache 2.0, please read that license about using, distribution and forking. Licensed works, modifications, and larger works may be distributed under different terms and without source code.

## License

This repository is released under the [Apache license](https://github.com/djthorpe/data/tree/7f02a4b2bcc64113cf15ee330a72d5dcbb54d60e/LICENSE/README.md):

&gt;

> Copyright 2021 David Thorpe and all other authors of the software
>
> Licensed under the Apache License, Version 2.0 \(the "License"\); you may not use this file except in compliance with the License. You may obtain a copy of the License at
>
> ```text
>   http://www.apache.org/licenses/LICENSE-2.0
> ```
>
> Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

