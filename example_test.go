package opengraph

import "fmt"

func ExampleFetch() {
	ogp, err := Fetch("https://ogp.me/")
	fmt.Println("title:", ogp.Title)
	fmt.Println("type:", ogp.Type)
	fmt.Println("error:", err)
	// Output:
	// title: Open Graph protocol
	// type: website
	// error: <nil>
}
