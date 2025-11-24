package use_cases

import (
	"context"
	"errors"
	"github.com/jorgejr568/exchange-register-go/internal/exchange/clients/exchangerate"
	clientMocks "github.com/jorgejr568/exchange-register-go/internal/exchange/clients/exchangerate/mocks"
	"github.com/jorgejr568/exchange-register-go/internal/exchange/entity"
	entityMocks "github.com/jorgejr568/exchange-register-go/internal/exchange/entity/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestSyncExchangeRateUseCase_Execute_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := entityMocks.NewMockExchangeService(ctrl)
	mockClient := clientMocks.NewMockClient(ctrl)
	useCase := NewSyncExchangeRateUseCase(mockService, mockClient)

	ctx := context.Background()
	req := entity.SyncExchangeRateRequest{
		SourceCurrency: "USD",
		TargetCurrency: "BRL",
	}

	clientReq := exchangerate.GetExchangeRateRequest{
		From: "USD",
		To:   "BRL",
	}
	clientResp := &exchangerate.GetExchangeRateResponse{
		Rate: 5.25,
	}

	mockClient.EXPECT().
		GetExchangeRate(ctx, clientReq).
		Return(clientResp, nil)

	mockService.EXPECT().
		ReceiveExchangeRate(ctx, "USD", "BRL", 5.25).
		Return(nil)

	// Act
	result, err := useCase.Execute(ctx, req)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 5.25, result.Rate)
}

func TestSyncExchangeRateUseCase_Execute_ClientError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := entityMocks.NewMockExchangeService(ctrl)
	mockClient := clientMocks.NewMockClient(ctrl)
	useCase := NewSyncExchangeRateUseCase(mockService, mockClient)

	ctx := context.Background()
	req := entity.SyncExchangeRateRequest{
		SourceCurrency: "USD",
		TargetCurrency: "BRL",
	}

	clientReq := exchangerate.GetExchangeRateRequest{
		From: "USD",
		To:   "BRL",
	}
	expectedError := errors.New("API rate limit exceeded")

	mockClient.EXPECT().
		GetExchangeRate(ctx, clientReq).
		Return(nil, expectedError)

	// Act
	result, err := useCase.Execute(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}

func TestSyncExchangeRateUseCase_Execute_ServiceError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := entityMocks.NewMockExchangeService(ctrl)
	mockClient := clientMocks.NewMockClient(ctrl)
	useCase := NewSyncExchangeRateUseCase(mockService, mockClient)

	ctx := context.Background()
	req := entity.SyncExchangeRateRequest{
		SourceCurrency: "EUR",
		TargetCurrency: "BRL",
	}

	clientReq := exchangerate.GetExchangeRateRequest{
		From: "EUR",
		To:   "BRL",
	}
	clientResp := &exchangerate.GetExchangeRateResponse{
		Rate: 5.75,
	}
	expectedError := errors.New("database connection failed")

	mockClient.EXPECT().
		GetExchangeRate(ctx, clientReq).
		Return(clientResp, nil)

	mockService.EXPECT().
		ReceiveExchangeRate(ctx, "EUR", "BRL", 5.75).
		Return(expectedError)

	// Act
	result, err := useCase.Execute(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}

func TestSyncExchangeRateUseCase_Execute_DifferentCurrencyPairs(t *testing.T) {
	testCases := []struct {
		name string
		from string
		to   string
		rate float64
	}{
		{"USD to EUR", "USD", "EUR", 0.92},
		{"GBP to USD", "GBP", "USD", 1.27},
		{"JPY to BRL", "JPY", "BRL", 0.034},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := entityMocks.NewMockExchangeService(ctrl)
			mockClient := clientMocks.NewMockClient(ctrl)
			useCase := NewSyncExchangeRateUseCase(mockService, mockClient)

			ctx := context.Background()
			req := entity.SyncExchangeRateRequest{
				SourceCurrency: tc.from,
				TargetCurrency: tc.to,
			}

			clientReq := exchangerate.GetExchangeRateRequest{
				From: tc.from,
				To:   tc.to,
			}
			clientResp := &exchangerate.GetExchangeRateResponse{
				Rate: tc.rate,
			}

			mockClient.EXPECT().
				GetExchangeRate(ctx, clientReq).
				Return(clientResp, nil)

			mockService.EXPECT().
				ReceiveExchangeRate(ctx, tc.from, tc.to, tc.rate).
				Return(nil)

			// Act
			result, err := useCase.Execute(ctx, req)

			// Assert
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, tc.rate, result.Rate)
		})
	}
}
