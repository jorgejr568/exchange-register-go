package exchange

import (
	"context"
	"github.com/jorgejr568/exchange-register-go/internal/exchange/clients/exchangerate"
	"github.com/jorgejr568/exchange-register-go/internal/exchange/entity"
)

type syncExchangeRateUseCase struct {
	exchangeService    entity.ExchangeService
	exchangeRateClient exchangerate.Client
}

func (s *syncExchangeRateUseCase) Execute(ctx context.Context, req entity.GetExchangeRateRequest) (*entity.GetExchangeRateResponse, error) {
	resp, err := s.exchangeRateClient.GetExchangeRate(ctx, exchangerate.GetExchangeRateRequest{
		From: req.SourceCurrency,
		To:   req.TargetCurrency,
	})
	if err != nil {
		return nil, err
	}

	err = s.exchangeService.ReceiveExchangeRate(ctx, req.SourceCurrency, req.TargetCurrency, resp.Rate)
	if err != nil {
		return nil, err
	}

	return &entity.GetExchangeRateResponse{
		Rate: resp.Rate,
	}, nil
}

func NewSyncExchangeRateUseCase(exchangeService entity.ExchangeService, exchangeRateClient exchangerate.Client) entity.SyncExchangeRateUseCase {
	return &syncExchangeRateUseCase{
		exchangeService:    exchangeService,
		exchangeRateClient: exchangeRateClient,
	}
}
