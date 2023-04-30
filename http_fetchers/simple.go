package http_fetchers

import (
	"fmt"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"strings"
)

type SimpleHTTPFetcher struct {
	httpClient *http.Client
}

func NewSimpleHTTPFetcher(client *http.Client) *SimpleHTTPFetcher {
	return &SimpleHTTPFetcher{
		httpClient: client,
	}
}

var DefaultSimpleHTTPFetcher = NewSimpleHTTPFetcher(http.DefaultClient)

func (fetcher *SimpleHTTPFetcher) Get(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res, err := fetcher.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(res.Header.Get("Content-Type"), "text/html") {
		return nil, fmt.Errorf("content type must be text/html")
	}

	return res.Body, nil
}
