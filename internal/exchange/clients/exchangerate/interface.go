package exchangerate

import "context"

type Client interface {
	// GetExchangeRate returns the exchange rate between two currencies.
	GetExchangeRate(ctx context.Context, request GetExchangeRateRequest) (*GetExchangeRateResponse, error)
}
