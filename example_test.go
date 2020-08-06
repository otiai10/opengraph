package opengraph

import (
	"fmt"
	"log"
	"net/http"
)

func ExampleFetch() {
	ogp, _ := Fetch("https://github.com/otiai10/gosseract")
	fmt.Println(ogp.Title)
	// Output: otiai10/gosseract
}

func ExampleOpenGraph_Parse() {

	client := http.DefaultClient
	res, err := client.Get("https://github.com/otiai10/amesh")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	ogp := new(OpenGraph)
	ogp.Parse(res.Body)
	fmt.Println(ogp.Title)
	// Output: otiai10/amesh
}
