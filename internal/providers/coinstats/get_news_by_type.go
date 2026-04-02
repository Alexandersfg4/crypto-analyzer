package coinstats

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Alexandersfg4/crypto-analyzer/internal/models"
)

func (s *Service) GetNewsByType(ctx context.Context, t models.NewsType, limit int) (models.GetNewsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("/news/type/%s", t), nil)
	if err != nil {
		return models.GetNewsResponse{}, fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return models.GetNewsResponse{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return models.GetNewsResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var news models.GetNewsResponse
	if err := json.NewDecoder(resp.Body).Decode(&news); err != nil {
		return models.GetNewsResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}
	return news, nil
}
