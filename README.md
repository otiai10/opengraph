# Open Graph Parser for Golang

[![Build Status](https://travis-ci.org/otiai10/opengraph.svg?branch=master)](https://travis-ci.org/otiai10/opengraph)

# in code

```go
package main

import (
	"fmt"
	"github.com/otiai10/opengraph"
)

func main() {
	og, err := opengraph.Fetch("http://github.com/otiai10/too")
	fmt.Printf("OpenGraph: %+v\nError: %v\n", og, err)
}
```

# in cli

```sh
% go get github.com/otiai10/opengraph/ogp
% ogp https://github.com/otiai10/too
```

# advanced

- [`og.Parse(body *io.Reader)`](https://godoc.org/github.com/otiai10/opengraph#OpenGraph.Parse) to re-use `*http.Response`
- [`og.HTTPClient`](https://godoc.org/github.com/otiai10/opengraph#OpenGraph) to customize `*http.Client` for fetching
- [`og.ToAbsURL()`](https://godoc.org/github.com/otiai10/opengraph#OpenGraph.ToAbsURL) to restore relative URL, e.g. `og.Favicon`
- ~~[`og.Fulfill()`](https://godoc.org/github.com/otiai10/opengraph#OpenGraph.Fulfill) to fill empty fileds.~~ You ain't gonna need it
