package openrouter

import (
	"github.com/revrost/go-openrouter"
)

func New(apiKey string) *Service {
	return &Service{
		client: openrouter.NewClient(apiKey),
	}
}
