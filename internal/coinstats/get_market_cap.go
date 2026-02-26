package coinstats

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Alexandersfg4/crypto-analyzer/internal/models"
)

func (s *Service) GetMarketCap(ctx context.Context) (models.MarketCap, error) {
	var result models.MarketCap
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "markets", nil)
	if err != nil {
		return result, fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return result, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, fmt.Errorf("error decoding market cap data: %w", err)
	}

	return result, nil
}
