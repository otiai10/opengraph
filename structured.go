package opengraph

/**
 * Structured Properties.
 * https://ogp.me/#structured
 */

// Image represents a structure of "og:image".
// "og:image" might have following properties:
//   - og:image:url
//   - og:image:secure_url
//   - og:image:type
//   - og:image:width
//   - og:image:height
//   - og:image:alt
type Image struct {
	URL       string `json:"url"`
	SecureURL string `json:"secure_url"`
	Type      string `json:"type"` // Content-Type
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Alt       string `json:"alt"`
}

// Video represents a structure of "og:video".
// "og:video" might have following properties:
//   - og:video:url
//   - og:video:secure_url
//   - og:video:type
//   - og:video:width
//   - og:video:height
type Video struct {
	URL       string `json:"url"`
	SecureURL string `json:"secure_url"`
	Type      string `json:"type"` // Content-Type
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

// Audio represents a structure of "og:audio".
// "og:audio" might have following properties:
//   - og:audio:url
//   - og:audio:secure_url
//   - og:audio:type
type Audio struct {
	URL       string `json:"url"`
	SecureURL string `json:"secure_url"`
	Type      string `json:"type"` // Content-Type
}

// Favicon represents an extra structure for "shortcut icon".
type Favicon struct {
	URL string `json:"url"`
}
