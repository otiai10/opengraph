# Open Graph Parser for Golang

in code

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

in cli

```sh
% go get github.com/otiai10/opengraph/ogp
% ogp https://github.com/otiai10/too
```
