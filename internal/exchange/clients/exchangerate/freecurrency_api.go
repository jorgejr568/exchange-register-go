package exchangerate

import (
	"context"
	"github.com/jorgejr568/freecurrencyapi-go/v2"
)

type freeCurrencyApiClient struct {
	client freecurrencyapi.Client
}

func (f freeCurrencyApiClient) GetExchangeRate(ctx context.Context, request GetExchangeRateRequest) (*GetExchangeRateResponse, error) {
	rate, err := f.client.Latest(ctx, freecurrencyapi.LatestRequest{
		BaseCurrency: request.From,
		Currencies:   []string{request.To},
	})
	if err != nil {
		return nil, err
	}

	return &GetExchangeRateResponse{
		Rate: rate.Rates[request.To],
	}, nil
}

func NewFreeCurrencyApiClient(client freecurrencyapi.Client) Client {
	return &freeCurrencyApiClient{
		client: client,
	}
}
