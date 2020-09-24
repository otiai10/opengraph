# Open Graph Parser for Golang

[![Go](https://github.com/otiai10/opengraph/workflows/Go/badge.svg)](https://github.com/otiai10/opengraph/actions)
[![codecov](https://codecov.io/gh/otiai10/opengraph/branch/master/graph/badge.svg)](https://codecov.io/gh/otiai10/opengraph)
[![GoDoc](https://godoc.org/github.com/otiai10/opengraph?status.svg)](https://pkg.go.dev/github.com/otiai10/opengraph)

# Code Example

```go
package main

import (
	"fmt"
	"github.com/otiai10/opengraph/v2"
)

func main() {
	ogp, err := opengraph.Fetch("https://github.com/")
	fmt.Printf("OpenGraph: %+v\nError: %v\n", ogp, err)
}
```

# CLI as a working example

```sh
% go get github.com/otiai10/opengraph/ogp
% ogp --help
```

For more details, see [ogp/main.go](https://github.com/otiai10/opengraph/blob/master/ogp/main.go).
