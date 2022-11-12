package exchangerate

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type httpClient struct {
	http    *http.Client
	baseUrl string
}

func (h httpClient) GetExchangeRate(ctx context.Context, request GetExchangeRateRequest) (*GetExchangeRateResponse, error) {
	url := fmt.Sprintf("%s/convert?from=%s&to=%s", h.baseUrl, request.From, request.To)
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	httpResponse, err := h.http.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", httpResponse.StatusCode)
	}

	var response GetExchangeRateResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func NewHTTPClient(http *http.Client, baseUrl string) Client {
	return &httpClient{
		http:    http,
		baseUrl: baseUrl,
	}
}
