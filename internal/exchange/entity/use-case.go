package entity

import "context"

type GetExchangeRateRequest struct {
	SourceCurrency string
	TargetCurrency string
}

type GetExchangeRateResponse struct {
	Rate float64
}

type SyncExchangeRateUseCase interface {
	Execute(ctx context.Context, req GetExchangeRateRequest) (*GetExchangeRateResponse, error)
}

type ExchangeService interface {
	// ReceiveExchangeRate creates a new exchange rate in the database.
	ReceiveExchangeRate(ctx context.Context, sourceCurrency, targetCurrency string, rate float64) error
}
