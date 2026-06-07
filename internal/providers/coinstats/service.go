package coinstats

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	apiHeader                = "X-API-KEY"
	baseURL                  = "https://openapiv1.coinstats.app"
	defaultRateLimitInterval = 500 * time.Millisecond
)

func NewService(apiKey string) *Service {
	return &Service{
		httpClient: &http.Client{
			Transport: &rt{
				baseURL:     baseURL,
				apiKey:      apiKey,
				rateLimit:   defaultRateLimitInterval,
				lastRequest: time.Time{},
			},
		},
	}
}

type Service struct {
	httpClient *http.Client
}

type rt struct {
	baseURL   string
	apiKey    string
	rateLimit time.Duration

	mu          sync.Mutex
	lastRequest time.Time
}

func (rt *rt) RoundTrip(req *http.Request) (*http.Response, error) {
	rt.waitRateLimit()

	baseURL, err := url.Parse(rt.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	// Resolve against the base URL to preserve query parameters.
	gotURL := baseURL.ResolveReference(req.URL)

	req.URL = gotURL
	req.Header.Set(apiHeader, rt.apiKey)
	return http.DefaultTransport.RoundTrip(req)
}

func (rt *rt) waitRateLimit() {
	if rt.rateLimit <= 0 {
		return
	}

	rt.mu.Lock()
	defer rt.mu.Unlock()

	if !rt.lastRequest.IsZero() {
		if wait := rt.rateLimit - time.Since(rt.lastRequest); wait > 0 {
			time.Sleep(wait)
		}
	}

	rt.lastRequest = time.Now()
}
