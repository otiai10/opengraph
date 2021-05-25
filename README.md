# Open Graph Parser for Golang

Go implementation of https://ogp.me/

[![reference](https://pkg.go.dev/badge/github.com/otiai10/opengraph)](https://pkg.go.dev/github.com/otiai10/opengraph/v2)
[![Go](https://github.com/otiai10/opengraph/workflows/Go/badge.svg)](https://github.com/otiai10/opengraph/actions)
[![codecov](https://codecov.io/gh/otiai10/opengraph/branch/main/graph/badge.svg?token=D4mPKqi9fH)](https://codecov.io/gh/otiai10/opengraph)
[![Go Report Card](https://goreportcard.com/badge/github.com/otiai10/opengraph)](https://goreportcard.com/report/github.com/otiai10/opengraph)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://github.com/otiai10/opengraph/blob/main/LICENSE)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/otiai10/opengraph?sort=semver)](https://pkg.go.dev/github.com/otiai10/opengraph/v2)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fotiai10%2Fopengraph.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fotiai10%2Fopengraph?ref=badge_shield)

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

# You can try CLI as a working example

```sh
% go get github.com/otiai10/opengraph/ogp
% ogp --help
% ogp -A otiai10.com
```

Just for fun 😉

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

or if you already have parsed `*html.Node`:

```go
err := ogp.Walk(node)
```

Do you wanna make absolute URLs?:
```go
ogp.Image[0].URL // /logo.png
ogp.ToAbs()
ogp.Image[0].URL // https://ogp.me/logo.png
```

# Issues

- https://github.com/otiai10/opengraph/issues


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fotiai10%2Fopengraph.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fotiai10%2Fopengraph?ref=badge_large)