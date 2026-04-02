package coinstats

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Alexandersfg4/crypto-analyzer/internal/models"
)

func (s *Service) GetFearAndGreed(ctx context.Context) (models.FearAndGreed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "insights/fear-and-greed", nil)
	if err != nil {
		return models.FearAndGreed{}, fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return models.FearAndGreed{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return models.FearAndGreed{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data models.FearAndGreed
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return data, fmt.Errorf("error decoding fear and greed data: %w", err)
	}

	return data, nil
}
