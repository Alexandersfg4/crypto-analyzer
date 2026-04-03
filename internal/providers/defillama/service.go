package defillama

import (
	"fmt"
	"net/http"
	"net/url"
)

const (
	baseURL = "https://api.llama.fi"
)

func NewService() *Service {
	return &Service{
		httpClient: &http.Client{
			Transport: &rt{
				baseURL: baseURL,
			},
		},
	}
}

type Service struct {
	httpClient *http.Client
}

type rt struct {
	baseURL string
}

func (rt *rt) RoundTrip(req *http.Request) (*http.Response, error) {
	baseURL, err := url.Parse(rt.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	// Resolve against the base URL to preserve query parameters.
	gotURL := baseURL.ResolveReference(req.URL)

	req.URL = gotURL
	return http.DefaultTransport.RoundTrip(req)
}
