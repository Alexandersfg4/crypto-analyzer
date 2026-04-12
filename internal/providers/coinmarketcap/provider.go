package coinmarketcap

import (
	"net/http"
)

func New(apiKey string) *Service {
	return &Service{
		httpClient: &http.Client{
			Transport: &rt{
				baseURL: baseURL,
				apiKey:  apiKey,
			},
		},
	}
}
