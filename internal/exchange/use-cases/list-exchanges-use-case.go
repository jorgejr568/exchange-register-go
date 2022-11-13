package use_cases

import (
	"context"
	"github.com/jorgejr568/exchange-register-go/internal/exchange/entity"
)

type listExchangesUseCase struct {
	exchangeService entity.ExchangeService
}

func (s *listExchangesUseCase) Execute(ctx context.Context, req entity.ListExchangesRequest) (*entity.ListExchangesResponse, error) {
	exchanges, err := s.exchangeService.ListExchanges(ctx, req.SourceCurrency, req.TargetCurrency)
	if err != nil {
		return nil, err
	}

	exchangesResponse := make(entity.ListExchangesResponse, len(exchanges))
	for i, exchange := range exchanges {
		lastAcquisition := exchange.CreatedAt
		if exchange.UpdatedAt != nil {
			lastAcquisition = *exchange.UpdatedAt
		}

		exchangesResponse[i] = entity.ExchangeResponse{
			ID:              exchange.ID,
			SourceCurrency:  exchange.BaseCurrency,
			TargetCurrency:  exchange.TargetCurrency,
			Rate:            exchange.Rate,
			LastAcquisition: lastAcquisition,
		}
	}

	return &exchangesResponse, nil
}

func NewListExchangesUseCase(exchangeService entity.ExchangeService) entity.ListExchangesUseCase {
	return &listExchangesUseCase{
		exchangeService: exchangeService,
	}
}
