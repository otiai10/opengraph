package opengraph

import (
	"net/http"
)

// Intent represents how to fetch, parse, and complete properties
// of this OpenGraph object.
// This SHOULD NOT have any meaning for "OpenGraph Protocol".
type Intent struct {
	URL        string       // Target URL of this intent.
	HTTPClient *http.Client // HTTP Client to be used for this intent.
	Strict     bool
}
