package opengraph

import (
	"strconv"
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
	case m.IsOGDescription():
		og.Description = m.Content
	case m.IsDescription() && og.Description == "":
		og.Description = m.Content
	case m.IsImage():
		og.Image = append(og.Image, &OGImage{URL: m.Content})
	case m.IsSiteName():
		og.SiteName = m.Content
	case m.IsImageProperty():
		if len(og.Image) == 0 {
			return nil
		}
		switch m.Property {
		case "og:image:width":
			og.Image[len(og.Image)-1].Width, _ = strconv.Atoi(m.Content)
		case "og:image:height":
			og.Image[len(og.Image)-1].Height, _ = strconv.Atoi(m.Content)
		}
	case m.IsType():
		og.Type = m.Content
	}
	return nil
}

// IsTitle returns if it can be "title" of OGP
func (m *Meta) IsTitle() bool {
	return m.Property == "og:title" && m.Content != ""
}

// IsDescription returns if it can be "description" of OGP
func (m *Meta) IsOGDescription() bool {
	return m.Property == "og:description" && m.Content != ""
}

// IsDescription returns if it can be "description" of OGP
func (m *Meta) IsDescription() bool {
	return m.Name == "description" && m.Content != ""
}

// IsImage returns if it can be a root of "og:image"
func (m *Meta) IsImage() bool {
	return m.Property == "og:image" || m.Property == "og:image:url"
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
