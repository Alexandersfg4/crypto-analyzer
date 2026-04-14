package defillama

import (
	"net/http"
)

func New() *Service {
	return &Service{
		httpClient: &http.Client{
			Transport: &rt{
				baseURL: baseURL,
			},
		},
	}
}
