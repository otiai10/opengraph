package opengraph

import "golang.org/x/net/html"

// Title represents any "<title ...>" HTML tag.
type Title struct {
	Text string
}

// TitleTag constructs Title.
func TitleTag(n *html.Node) *Title {
	t := new(Title)
	if n.FirstChild != nil {
		t.Text = n.FirstChild.Data
	}
	return t
}

// Contribute contributes to OpenGraph
func (t *Title) Contribute(og *OpenGraph) error {
	if og.Title == "" && t.Text != "" {
		og.Title = t.Text
	}
	return nil
}
