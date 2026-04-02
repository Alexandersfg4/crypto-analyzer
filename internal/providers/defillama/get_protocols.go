package defillama

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Alexandersfg4/crypto-analyzer/internal/models"
)

func (s *Service) GetProtocols(ctx context.Context) (models.GetProtocolsResponse, error) {
	var protocols models.GetProtocolsResponse

	u := url.URL{
		Path: "/protocols",
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		u.String(),
		nil,
	)
	if err != nil {
		return protocols, fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return protocols, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return protocols, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&protocols); err != nil {
		return protocols, fmt.Errorf("failed to decode response: %w", err)
	}

	return protocols, nil
}
