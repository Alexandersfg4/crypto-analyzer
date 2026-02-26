package coinstats

import (
	"fmt"
	"net/http"
	"net/url"
)

const apiHeader = "X-API-KEY"

func NewService(baseURL, apiKey string) *Service {
	return &Service{
		httpClient: &http.Client{
			Transport: &rt{
				baseURL: baseURL,
				apiKey:  apiKey,
			},
		},
	}
}

type Service struct {
	httpClient *http.Client
}

type rt struct {
	baseURL string
	apiKey  string
}

func (rt *rt) RoundTrip(req *http.Request) (*http.Response, error) {
	fullPath, err := url.JoinPath(rt.baseURL, req.URL.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to join path: %w", err)
	}
	gotUrl, err := url.Parse(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	req.URL = gotUrl
	req.Header.Set(apiHeader, rt.apiKey)
	return http.DefaultTransport.RoundTrip(req)
}
