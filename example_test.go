package opengraph

import (
	"fmt"
	"net/http"
)

func ExampleFetch() {
	ogp, _ := Fetch("https://github.com/otiai10/gosseract")
	fmt.Println(ogp.Title)
	// Output: otiai10/gosseract
}

func ExampleOpenGraph_Parse() {

	client := http.DefaultClient
	res, _ := client.Get("https://github.com/otiai10/amesh")
	defer res.Body.Close()

	ogp := new(OpenGraph)
	ogp.Parse(res.Body)
	fmt.Println(ogp.Title)
	// Output: otiai10/amesh
}
