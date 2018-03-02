# Open Graph Parser for Golang

[![Build Status](https://travis-ci.org/otiai10/opengraph.svg?branch=master)](https://travis-ci.org/otiai10/opengraph) [![codecov](https://codecov.io/gh/otiai10/opengraph/branch/master/graph/badge.svg)](https://codecov.io/gh/otiai10/opengraph) [![GoDoc](https://godoc.org/github.com/otiai10/opengraph?status.svg)](https://godoc.org/github.com/otiai10/opengraph)

# Code Example

```go
package main

import (
	"fmt"
	"github.com/otiai10/opengraph"
)

func main() {
	og, err := opengraph.Fetch("https://www.youtube.com/watch?v=5blm22DeeHY")
	fmt.Printf("OpenGraph: %+v\nError: %v\n", og, err)
}
```

# CLI as a working example

```sh
% go get github.com/otiai10/opengraph/ogp
% ogp --help
```

For more details, see [ogp/main.go](https://github.com/otiai10/opengraph/blob/master/ogp/main.go).

# Advanced

- [`og.Parse(body *io.Reader)`](https://godoc.org/github.com/otiai10/opengraph#OpenGraph.Parse) to re-use `*http.Response`
- [`og.HTTPClient`](https://godoc.org/github.com/otiai10/opengraph#OpenGraph) to customize `*http.Client` for fetching
- [`og.ToAbsURL()`](https://godoc.org/github.com/otiai10/opengraph#OpenGraph.ToAbsURL) to restore relative URL, e.g. `og.Favicon`
- ~~[`og.Fulfill()`](https://godoc.org/github.com/otiai10/opengraph#OpenGraph.Fulfill) to fill empty fileds.~~ You ain't gonna need it
