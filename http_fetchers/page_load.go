package http_fetchers

import (
	"bytes"
	"github.com/chromedp/chromedp"
	"golang.org/x/net/context"
	"io"
	"time"
)

type PageLoadHTTPFetcher struct {
	fetchSeconds int
}

func NewPageLoadHTTPFetcher(fetchSeconds int) *PageLoadHTTPFetcher {
	return &PageLoadHTTPFetcher{
		fetchSeconds: fetchSeconds,
	}
}

func (fetcher *PageLoadHTTPFetcher) Get(ctx context.Context, url string) (io.ReadCloser, error) {
	var htmlContent string

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(time.Duration(fetcher.fetchSeconds)*time.Second),
		chromedp.InnerHTML("html", &htmlContent),
	)

	if err != nil {
		return nil, err
	}

	return io.NopCloser(bytes.NewReader([]byte(htmlContent))), nil
}
