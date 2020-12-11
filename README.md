# Open Graph Parser for Golang

Yet another implementation of https://ogp.me/ by Go.

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
	fmt.Println(ogp, err)
}
```

# Advanced usage

Set an option for fetching:
```go
intent := opengraph.Intent{
	Context:     ctx,
	HTTPClient:  client,
	Strict:      true,
	TrustedTags: []string{"meta", "title"},
}
ogp, err := opengraph.Fetch("https://ogp.me", intent)
```

Use any `io.Reader` as a data source:
```go
f, _ := os.Open("my_test.html")
defer f.Close()
ogp := &opengraph.OpenGraph{}
err := ogp.Parse(f)
```

of if you already have parsed `*html.Node`:

```go
err := ogp.Walk(node)
```

Do you wanna make absolute URLs?:
```go
ogp.Image[0].URL // /logo.png
ogp.ToAbs()
ogp.Image[0].URL // https://ogp.me/logo.png
```

# CLI as a working example

```sh
% go get github.com/otiai10/opengraph/ogp
% ogp --help
% ogp -A otiai10.com
```

# Issues

- https://github.com/otiai10/opengraph/issues
