package coinmarketcap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Alexandersfg4/crypto-analyzer/internal/models"
)

func (s *Service) GetListingsLatest(ctx context.Context, start, limit int) (models.ListingsLatestResponse, error) {
	var listingsLatest models.ListingsLatestResponse

	u := url.URL{
		Path: "/v3/cryptocurrency/listings/latest",
		RawQuery: url.Values{
			"limit": {strconv.Itoa(limit)},
			"start": {strconv.Itoa(start)},
		}.Encode(),
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		u.String(),
		nil,
	)
	if err != nil {
		return listingsLatest, fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return listingsLatest, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return listingsLatest, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&listingsLatest); err != nil {
		return listingsLatest, fmt.Errorf("failed to decode response: %w", err)
	}

	return listingsLatest, nil
}
