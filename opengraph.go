// Package opengraph implements and parses "The Open Graph Protocol" of web pages.
// See http://ogp.me/ for more information.
package opengraph

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

// OpenGraph represents web page information according to OGP <ogp.me>,
// and some more additional informations like URL.Host and so.
type OpenGraph struct {

	// Basics
	Title    string
	Type     string
	URL      URL
	SiteName string

	// Structures
	Image []*OGImage
	Video []*OGVideo
	Audio []*OGAudio

	// Optionals
	Description string
	Determiner  string // TODO: enum?
	Locale      string
	LocaleAlt   []string

	// Additionals
	Favicon string

	// Utils
	HTTPClient *http.Client `json:"-"`
	Error      error        `json:"-"`
}

// URL includes *url.URL
type URL struct {
	Source string
	*url.URL
}

// New creates new OpenGraph struct with specified URL.
func New(rawurl string) *OpenGraph {
	og := new(OpenGraph)
	og.HTTPClient = http.DefaultClient
	og.Image = []*OGImage{}
	og.Video = []*OGVideo{}
	og.Audio = []*OGAudio{}
	og.LocaleAlt = []string{}
	og.Favicon = "/favicon.ico"
	u, err := url.Parse(rawurl)
	if err != nil {
		og.Error = err
		return og
	}
	og.URL = URL{Source: u.String(), URL: u}
	return og
}

// Fetch creates and parses OpenGraph with specified URL.
func Fetch(rawurl string, customHTTPClient ...*http.Client) (*OpenGraph, error) {
	return FetchWithContext(context.Background(), rawurl, customHTTPClient...)
}

// FetchWithContext creates and parses OpenGraph with specified URL.
// Timeout can be handled with provided context.
func FetchWithContext(ctx context.Context, rawurl string, customHTTPClient ...*http.Client) (*OpenGraph, error) {
	og := New(rawurl)
	if og.Error != nil {
		return og, og.Error
	}

	// Use custom http client if given
	if len(customHTTPClient) != 0 {
		og.HTTPClient = customHTTPClient[0]
	}

	req, err := http.NewRequest("GET", og.URL.String(), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res, err := og.HTTPClient.Do(req)
	if err != nil {
		return og, err
	}
	defer res.Body.Close()

	if !strings.HasPrefix(res.Header.Get("Content-Type"), "text/html") {
		return og, fmt.Errorf("Content type must be text/html")
	}

	if err = og.Parse(res.Body); err != nil {
		return og, err
	}

	return og, err
}

// Parse parses http.Response.Body and construct OpenGraph informations.
// Caller should close body after it get parsed.
func (og *OpenGraph) Parse(body io.Reader) error {
	if og.Error != nil {
		return og.Error
	}
	node, err := html.Parse(body)
	if err != nil {
		return err
	}
	og.walk(node)
	return nil
}

func (og *OpenGraph) satisfied() bool {
	return false
}

func (og *OpenGraph) walk(n *html.Node) error {

	if og.satisfied() {
		return nil
	}

	if n.Type == html.ElementNode {
		switch n.Data {
		case "title":
			return TitleTag(n).Contribute(og)
		case "meta":
			return MetaTag(n).Contribute(og)
		case "link":
			return LinkTag(n).Contribute(og)
		}
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		og.walk(child)
	}

	return nil
}

// ToAbsURL make og.Image and og.Favicon absolute URL if relative.
func (og *OpenGraph) ToAbsURL() *OpenGraph {
	for _, img := range og.Image {
		img.URL = og.abs(img.URL)
	}
	og.Favicon = og.abs(og.Favicon)
	return og
}

// abs make given URL absolute.
func (og *OpenGraph) abs(raw string) string {
	u, _ := url.Parse(raw)
	if u.IsAbs() {
		return raw
	}
	if u.Scheme == "" {
		u.Scheme = og.URL.Scheme
	}
	if u.Host == "" {
		u.Host = og.URL.Host
	}
	if !filepath.IsAbs(raw) {
		u.Path = path.Join(filepath.Dir(og.URL.Path), u.Path)
	}
	return u.String()
}

// Fulfill fulfills OG informations with some expectations.
func (og *OpenGraph) Fulfill() error {
	if og.SiteName == "" {
		og.SiteName = og.URL.Host
	}
	return nil
}
