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
func MetaTag(node *html.Node) *Meta {
	meta := new(Meta)
	for _, attr := range node.Attr {
		switch attr.Key {
		case "property":
			meta.Property = attr.Val
		case "content":
			meta.Content = attr.Val
		case "name":
			meta.Name = attr.Val
		}
	}
	return meta
}

// Contribute ...
func (meta *Meta) Contribute(og *OpenGraph) (err error) {
	switch {
	case meta.IsTitle():
		og.Title = meta.Content
	case meta.IsOGDescription():
		og.Description = meta.Content
	case meta.IsDescription() && !og.Intent.Strict && og.Description == "":
		og.Description = meta.Content
	case meta.IsSiteName():
		og.SiteName = meta.Content
	case meta.IsImage():
		og.Image = append(og.Image, Image{URL: meta.Content})
	case meta.IsImageProperty():
		if len(og.Image) == 0 {
			return nil
		}
		switch meta.Property {
		case "og:image:width":
			og.Image[len(og.Image)-1].Width, err = strconv.Atoi(meta.Content)
		case "og:image:height":
			og.Image[len(og.Image)-1].Height, err = strconv.Atoi(meta.Content)
		}
	case meta.IsType():
		og.Type = meta.Content
	case meta.IsURL():
		og.URL = meta.Content
	}
	return err
}

// IsTitle returns if it can be "title" of OGP
func (meta *Meta) IsTitle() bool {
	return meta.Property == "og:title" && meta.Content != ""
}

// IsOGDescription returns if it can be "description" of OGP
func (meta *Meta) IsOGDescription() bool {
	return meta.Property == "og:description" && meta.Content != ""
}

// IsDescription returns if it can be "description" of OGP.
// CAUTION: This property SHOULD NOT be used when Intent.Strict == true.
func (meta *Meta) IsDescription() bool {
	return meta.Name == "description" && meta.Content != ""
}

// IsImage returns if it can be a root of "og:image"
func (meta *Meta) IsImage() bool {
	return meta.Property == "og:image" || meta.Property == "og:image:url"
}

// IsImageProperty returns if it can be a property of "og:image" struct
func (meta *Meta) IsImageProperty() bool {
	return strings.HasPrefix(meta.Property, "og:image:")
}

// IsType returns if it can be "og:type"
func (meta *Meta) IsType() bool {
	return meta.Property == "og:type"
}

// IsSiteName returns if it can be "og:site_name"
func (meta *Meta) IsSiteName() bool {
	return meta.Property == "og:site_name"
}

// IsURL returns if it can be "og:url"
func (meta *Meta) IsURL() bool {
	return meta.Property == "og:url"
}
