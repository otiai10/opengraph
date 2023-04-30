package opengraph

import (
	"context"
	"io"
)

type HttpFetcher interface {
	Get(ctx context.Context, url string) (io.ReadCloser, error)
}

// Intent represents how to fetch, parse, and complete properties
// of this OpenGraph object.
// This SHOULD NOT have any meaning for "OpenGraph Protocol".
type Intent struct {

	// URL of this intent to fetch an OGP.
	// This does NOT mean `og:url` of the page.
	URL string

	// Context of the web request of this Intent.
	Context context.Context
	// HTTP Client to be used for this intent.
	HTTPFetcher HttpFetcher

	// Scrict is just an alias of `TrustedTags`.
	// `Strict == true` means `TrustedTags = ["meta"]`,
	// and `Strict == false` means `TrustedTags == ["meta", "title", "link"]`.
	Strict bool
	// TrustedTags specify which tags to be respected.
	TrustedTags []string
}
