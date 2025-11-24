package use_cases

import (
	"context"
	"github.com/jorgejr568/exchange-register-go/internal/exchange/entity"
	"github.com/jorgejr568/exchange-register-go/internal/exchange/entity/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestListExchangesUseCase_Execute_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockExchangeService(ctrl)
	useCase := NewListExchangesUseCase(mockService)

	ctx := context.Background()
	req := entity.ListExchangesRequest{
		SourceCurrency: "USD",
		TargetCurrency: "BRL",
	}

	now := time.Now()
	expectedExchanges := []entity.Exchange{
		{
			ID:             1,
			BaseCurrency:   "USD",
			TargetCurrency: "BRL",
			Rate:           5.25,
			CreatedAt:      now,
			UpdatedAt:      &now,
		},
	}

	mockService.EXPECT().
		ListExchanges(ctx, "USD", "BRL").
		Return(expectedExchanges, nil)

	// Act
	result, err := useCase.Execute(ctx, req)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, *result, 1)

	response := (*result)[0]
	assert.Equal(t, uint64(1), response.ID)
	assert.Equal(t, "USD", response.SourceCurrency)
	assert.Equal(t, "BRL", response.TargetCurrency)
	assert.Equal(t, 5.25, response.Rate)
	assert.Equal(t, now, response.LastAcquisition)
}

func TestListExchangesUseCase_Execute_WithoutFilters(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockExchangeService(ctrl)
	useCase := NewListExchangesUseCase(mockService)

	ctx := context.Background()
	req := entity.ListExchangesRequest{
		SourceCurrency: "",
		TargetCurrency: "",
	}

	now := time.Now()
	expectedExchanges := []entity.Exchange{
		{
			ID:             1,
			BaseCurrency:   "USD",
			TargetCurrency: "BRL",
			Rate:           5.25,
			CreatedAt:      now,
			UpdatedAt:      nil,
		},
		{
			ID:             2,
			BaseCurrency:   "EUR",
			TargetCurrency: "BRL",
			Rate:           5.75,
			CreatedAt:      now,
			UpdatedAt:      nil,
		},
	}

	mockService.EXPECT().
		ListExchanges(ctx, "", "").
		Return(expectedExchanges, nil)

	// Act
	result, err := useCase.Execute(ctx, req)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, *result, 2)

	// Verify first exchange
	firstExchange := (*result)[0]
	assert.Equal(t, "USD", firstExchange.SourceCurrency)
	assert.Equal(t, "BRL", firstExchange.TargetCurrency)
	assert.Equal(t, 5.25, firstExchange.Rate)
	assert.Equal(t, now, firstExchange.LastAcquisition) // Should use CreatedAt when UpdatedAt is nil
}

func TestListExchangesUseCase_Execute_Error(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockExchangeService(ctrl)
	useCase := NewListExchangesUseCase(mockService)

	ctx := context.Background()
	req := entity.ListExchangesRequest{
		SourceCurrency: "USD",
		TargetCurrency: "BRL",
	}

	mockService.EXPECT().
		ListExchanges(ctx, "USD", "BRL").
		Return(nil, assert.AnError)

	// Act
	result, err := useCase.Execute(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestListExchangesUseCase_Execute_EmptyResult(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockExchangeService(ctrl)
	useCase := NewListExchangesUseCase(mockService)

	ctx := context.Background()
	req := entity.ListExchangesRequest{
		SourceCurrency: "USD",
		TargetCurrency: "XYZ",
	}

	emptyExchanges := []entity.Exchange{}
	mockService.EXPECT().
		ListExchanges(ctx, "USD", "XYZ").
		Return(emptyExchanges, nil)

	// Act
	result, err := useCase.Execute(ctx, req)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, *result, 0)
}

func TestListExchangesUseCase_Execute_UsesUpdatedAtWhenAvailable(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockExchangeService(ctrl)
	useCase := NewListExchangesUseCase(mockService)

	ctx := context.Background()
	req := entity.ListExchangesRequest{
		SourceCurrency: "USD",
		TargetCurrency: "BRL",
	}

	createdAt := time.Now().Add(-1 * time.Hour)
	updatedAt := time.Now()
	expectedExchanges := []entity.Exchange{
		{
			ID:             1,
			BaseCurrency:   "USD",
			TargetCurrency: "BRL",
			Rate:           5.25,
			CreatedAt:      createdAt,
			UpdatedAt:      &updatedAt,
		},
	}

	mockService.EXPECT().
		ListExchanges(ctx, "USD", "BRL").
		Return(expectedExchanges, nil)

	// Act
	result, err := useCase.Execute(ctx, req)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, *result, 1)

	response := (*result)[0]
	assert.Equal(t, updatedAt, response.LastAcquisition) // Should use UpdatedAt when available
}
