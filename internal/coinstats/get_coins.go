package coinstats

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Alexandersfg4/crypto-analyzer/internal/models"
)

func (s *Service) GetCoins(ctx context.Context, limit int) (models.Coins, error) {
	var coins models.Coins

	u := url.URL{
		Path: "/coins",
		RawQuery: url.Values{
			"sortBy":  {"marketCap"},
			"limit":   {strconv.Itoa(limit)},
			"sortDir": {"desc"},
		}.Encode(),
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		u.String(),
		nil,
	)
	if err != nil {
		return coins, fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return coins, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return coins, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&coins); err != nil {
		return coins, fmt.Errorf("failed to decode response: %w", err)
	}

	return coins, nil
}
