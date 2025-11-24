package entity

//go:generate mockgen -destination=mocks/mock_use_case.go -package=mocks github.com/jorgejr568/exchange-register-go/internal/exchange/entity SyncExchangeRateUseCase,ListExchangesUseCase,ExchangeService

import "context"

type SyncExchangeRateRequest struct {
	SourceCurrency string
	TargetCurrency string
}

type SyncExchangeRateResponse struct {
	Rate float64
}

type ListExchangesRequest struct {
	SourceCurrency string `json:"source_currency"`
	TargetCurrency string `json:"target_currency"`
}

type ListExchangesResponse []ExchangeResponse

type SyncExchangeRateUseCase interface {
	Execute(ctx context.Context, req SyncExchangeRateRequest) (*SyncExchangeRateResponse, error)
}

type ListExchangesUseCase interface {
	Execute(ctx context.Context, req ListExchangesRequest) (*ListExchangesResponse, error)
}

type ExchangeService interface {
	// ReceiveExchangeRate creates a new exchange rate in the database.
	ReceiveExchangeRate(ctx context.Context, sourceCurrency, targetCurrency string, rate float64) error

	// ListExchanges returns a list of exchanges.
	ListExchanges(ctx context.Context, sourceCurrency, targetCurrency string) ([]Exchange, error)
}
