package openrouter

import (
	"github.com/revrost/go-openrouter"
)

type Service struct {
	client *openrouter.Client
}

func NewService(apiKey string) *Service {
	return &Service{
		client: openrouter.NewClient(apiKey),
	}
}
