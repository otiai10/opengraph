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
	case meta.IsPropertyOf("og:image"):
		if len(og.Image) == 0 {
			return nil
		}
		switch meta.Property {
		case "og:image:width":
			og.Image[len(og.Image)-1].Width, err = strconv.Atoi(meta.Content)
		case "og:image:height":
			og.Image[len(og.Image)-1].Height, err = strconv.Atoi(meta.Content)
		}
	case meta.IsAudio():
		og.Audio = append(og.Audio, Audio{URL: meta.Content})
	case meta.IsVideo():
		og.Video = append(og.Video, Video{URL: meta.Content})
	case meta.IsPropertyOf("og:video"):
		if len(og.Video) == 0 {
			return nil
		}
		switch meta.Property {
		case "og:video:type":
			og.Video[len(og.Video)-1].Type = meta.Content
		case "og:video:secure_url":
			og.Video[len(og.Video)-1].SecureURL = meta.Content
		case "og:video:width":
			og.Video[len(og.Video)-1].Width, err = strconv.Atoi(meta.Content)
		case "og:video:height":
			og.Video[len(og.Video)-1].Height, err = strconv.Atoi(meta.Content)
		case "og:video:duration":
			og.Video[len(og.Video)-1].Duration, err = strconv.Atoi(meta.Content)
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

// IsPropertyOf returns if it can be a property of specified struct
func (meta *Meta) IsPropertyOf(name string) bool {
	return strings.HasPrefix(meta.Property, name+":")
}

// IsAudio reeturns if it can be a root of "og:audio"
func (meta *Meta) IsAudio() bool {
	return meta.Property == "og:audio" || meta.Property == "og:audio:url"
}

// IsVideo returns if it can be a root of "og:video"
func (meta *Meta) IsVideo() bool {
	return meta.Property == "og:video" || meta.Property == "og:video:url"
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
