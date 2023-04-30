// Package opengraph implements and parses "The Open Graph Protocol" of web pages.
// See http://ogp.me/ for more information.
package opengraph

import (
	"context"
	"fmt"
	"github.com/otiai10/opengraph/v2/http_fetchers"
	"io"
	"net/url"
	"path"
	"strings"

	"golang.org/x/net/html"
)

const (
	// HTMLMetaTag is a tag name of <meta>
	HTMLMetaTag string = "meta"
	// HTMLLinkTag is a tag name of <link>
	HTMLLinkTag string = "link"
	// HTMLTitleTag is a tag name of <title>
	HTMLTitleTag string = "title"
)

// OpenGraph represents web page information according to OGP <ogp.me>,
// and some more additional informations like URL.Host and so.
type OpenGraph struct {

	// Basic Metadata
	// https://ogp.me/#metadata
	Title string  `json:"title"`
	Type  string  `json:"type"`
	Image []Image `json:"image"` // could be multiple
	URL   string  `json:"url"`

	// Optional Metadata
	// https://ogp.me/#optional
	Audio       []Audio  `json:"audio"` // could be multiple
	Description string   `json:"description"`
	Determiner  string   `json:"determiner"` // TODO: enum of (a, an, the, "", auto)
	Locale      string   `json:"locale"`
	LocaleAlt   []string `json:"locale_alternate"`
	SiteName    string   `json:"site_name"`
	Video       []Video  `json:"video"`

	// Additional (unofficial)
	Favicon Favicon `json:"favicon"`

	// Intent represents how to fetch, parse, and complete properties
	// of this OpenGraph object.
	// This SHOULD NOT have any meaning for "OpenGraph Protocol".
	Intent Intent `json:"-"`
}

// New constructs new OpenGraph struct and fill nullable fields.
func New(rawurl string) *OpenGraph {
	return &OpenGraph{
		Image:     []Image{},
		Audio:     []Audio{},
		Video:     []Video{},
		LocaleAlt: []string{},
		Intent: Intent{
			URL: rawurl,
		},
	}
}

// Fetch creates and parses OpenGraph with specified URL.
func Fetch(url string, intent ...Intent) (*OpenGraph, error) {
	ogp := New(url)
	if len(intent) > 0 {
		ogp.Intent = intent[0]
	}
	ogp.Intent.URL = url
	err := ogp.Fetch()
	return ogp, err
}

// Fetch ...
func (og *OpenGraph) Fetch() error {
	if og.Intent.URL == "" {
		return fmt.Errorf("no URL given yet")
	}

	if og.Intent.HTTPFetcher == nil {
		og.Intent.HTTPFetcher = http_fetchers.DefaultSimpleHTTPFetcher
	}

	if og.Intent.Context == nil {
		og.Intent.Context = context.Background()
	}

	body, err := og.Intent.HTTPFetcher.Get(og.Intent.Context, og.Intent.URL)
	if err != nil {
		return err
	}
	defer body.Close()

	if err = og.Parse(body); err != nil {
		return err
	}

	if !og.Intent.Strict && og.Favicon.URL == "" {
		og.Favicon.URL = "/favicon.ico"
	}

	return nil
}

// Parse parses http.Response.Body and construct OpenGraph informations.
// Caller should close body after it gets parsed.
func (og *OpenGraph) Parse(body io.Reader) error {
	node, err := html.Parse(body)
	if err != nil {
		return err
	}
	return og.Walk(node)
}

// Walk scans HTML nodes to pick up meaningful OGP data.
func (og *OpenGraph) Walk(node *html.Node) error {
	if len(og.Intent.TrustedTags) == 0 {
		if og.Intent.Strict {
			og.Intent.TrustedTags = []string{HTMLMetaTag}
		} else {
			og.Intent.TrustedTags = []string{HTMLMetaTag, HTMLTitleTag, HTMLLinkTag}
		}
	}
	return og.walk(node)
}

func (og *OpenGraph) walk(node *html.Node) error {
	if node.Type == html.ElementNode {
		switch {
		case node.Data == HTMLMetaTag && og.trust(node.Data):
			return MetaTag(node).Contribute(og)
		case node.Data == HTMLTitleTag && og.trust(node.Data):
			return TitleTag(node).Contribute(og)
		case node.Data == HTMLLinkTag && og.trust(node.Data):
			return LinkTag(node).Contribute(og)
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		og.walk(child)
	}

	return nil
}

func (og *OpenGraph) trust(tagName string) bool {
	for _, name := range og.Intent.TrustedTags {
		if name == tagName {
			return true
		}
	}
	return false
}

// ToAbs makes all relative URLs to absolute URLs
// by applying hostname of ogp.URL or Intent.URL.
func (og *OpenGraph) ToAbs() error {
	raw := og.URL
	if raw == "" {
		raw = og.Intent.URL
	}
	base, err := url.Parse(raw)
	if err != nil {
		return err
	}
	// For og:image.
	for i, img := range og.Image {
		og.Image[i].URL = og.joinToAbsolute(base, img.URL)
	}
	// For og:audio
	for i, audio := range og.Audio {
		og.Audio[i].URL = og.joinToAbsolute(base, audio.URL)
	}
	// For og:video
	for i, video := range og.Video {
		og.Video[i].URL = og.joinToAbsolute(base, video.URL)
	}
	// For favicon
	if og.Favicon.URL != "" {
		og.Favicon.URL = og.joinToAbsolute(base, og.Favicon.URL)
	}
	return nil
}

func (og *OpenGraph) joinToAbsolute(base *url.URL, relpath string) string {
	src, err := url.Parse(relpath)
	if err == nil && src.IsAbs() {
		return src.String()
	}
	if strings.HasPrefix(relpath, "//") {
		return fmt.Sprintf("%s:%s", base.Scheme, relpath)
	}
	if strings.HasPrefix(relpath, "/") {
		return fmt.Sprintf("%s://%s%s", base.Scheme, base.Host, relpath)
	}
	return fmt.Sprintf("%s://%s%s", base.Scheme, base.Host, path.Join(base.Path, relpath))
}
