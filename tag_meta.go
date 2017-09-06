package opengraph

import (
	"strings"

	"golang.org/x/net/html"
)

// Meta represents any "<meta ...>" HTML tag.
type Meta struct {
	Name     string
	Property string
	Content  string
}

// MetaTag constructs MetaTag.
func MetaTag(n *html.Node) *Meta {
	m := new(Meta)
	for _, attr := range n.Attr {
		switch attr.Key {
		case "property":
			m.Property = attr.Val
		case "content":
			m.Content = attr.Val
		case "name":
			m.Name = attr.Val
		}
	}
	return m
}

// Contribute ...
func (m *Meta) Contribute(og *OpenGraph) error {
	switch {
	case m.IsTitle():
		og.Title = m.Content
	case m.IsImage():
		og.Image = append(og.Image, &OGImage{URL: m.Content})
	case m.IsImageProperty():
		og.SiteName = m.Content
	case m.IsType():
		og.Type = m.Content
	}
	return nil
}

// IsTitle returns if it can be "title" of Oph
func (m *Meta) IsTitle() bool {
	return m.Property == "og:title"
}

// IsImage returns if it can be a root of "og:image"
func (m *Meta) IsImage() bool {
	return m.Property == "og:image"
}

// IsImageProperty retuns if it can be a property of "og:image" struct
func (m *Meta) IsImageProperty() bool {
	return strings.HasPrefix(m.Property, "og:image:")
}

// IsType returns if it can be "og:type"
func (m *Meta) IsType() bool {
	return m.Property == "og:type"
}

// IsSiteName returns if it can be "og:site_name"
func (m *Meta) IsSiteName() bool {
	return m.Property == "og:site_name"
}
