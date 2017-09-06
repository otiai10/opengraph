package opengraph

import "golang.org/x/net/html"

// Link represents any "<link ...>" HTML tag
type Link struct {
	Rel  string
	Href string
}

// LinkTag constructs Link
func LinkTag(n *html.Node) *Link {
	link := new(Link)
	for _, attr := range n.Attr {
		switch attr.Key {
		case "rel":
			link.Rel = attr.Val
		case "href":
			link.Href = attr.Val
		}
	}
	return link
}

// Contribute contributes OpenGraph
func (link *Link) Contribute(og *OpenGraph) error {
	switch {
	case link.IsFavicon():
		og.Favicon = link.Href
	}
	return nil
}

// IsFavicon returns if it can be "favicon" of *opengraph.OpenGraph
func (link *Link) IsFavicon() bool {
	return link.Rel == "shortcut icon" || link.Rel == "icon"
}
