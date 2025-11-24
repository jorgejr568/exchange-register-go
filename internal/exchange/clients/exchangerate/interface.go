package exchangerate

//go:generate mockgen -destination=mocks/mock_client.go -package=mocks github.com/jorgejr568/exchange-register-go/internal/exchange/clients/exchangerate Client

import "context"

type Client interface {
	// GetExchangeRate returns the exchange rate between two currencies.
	GetExchangeRate(ctx context.Context, request GetExchangeRateRequest) (*GetExchangeRateResponse, error)
}
