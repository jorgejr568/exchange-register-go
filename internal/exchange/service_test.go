package exchange

import (
	"context"
	"errors"
	"github.com/jorgejr568/exchange-register-go/internal/exchange/entity"
	"github.com/jorgejr568/exchange-register-go/internal/infra"
	"github.com/jorgejr568/exchange-register-go/internal/infra/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

// mockResult implements sql.Result for testing
type mockResult struct {
	lastInsertId int64
	rowsAffected int64
}

func (m mockResult) LastInsertId() (int64, error) {
	return m.lastInsertId, nil
}

func (m mockResult) RowsAffected() (int64, error) {
	return m.rowsAffected, nil
}

func TestKsqlExchangeService_ListExchanges_NoFilters(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	service := NewKSQLExchangeService(mockDB)

	ctx := context.Background()
	expectedQuery := "SELECT * FROM exchanges"

	mockDB.EXPECT().
		Query(ctx, gomock.Any(), expectedQuery).
		DoAndReturn(func(ctx context.Context, target interface{}, query string, args ...interface{}) error {
			ptr := target.(*[]entity.Exchange)
			*ptr = []entity.Exchange{
				{ID: 1, BaseCurrency: "USD", TargetCurrency: "BRL", Rate: 5.25},
				{ID: 2, BaseCurrency: "EUR", TargetCurrency: "BRL", Rate: 5.75},
			}
			return nil
		})

	// Act
	result, err := service.ListExchanges(ctx, "", "")

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "USD", result[0].BaseCurrency)
	assert.Equal(t, "EUR", result[1].BaseCurrency)
}

func TestKsqlExchangeService_ListExchanges_WithSourceFilter(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	service := NewKSQLExchangeService(mockDB)

	ctx := context.Background()
	expectedQuery := "SELECT * FROM exchanges WHERE base_currency = $1"

	mockDB.EXPECT().
		Query(ctx, gomock.Any(), expectedQuery, "USD").
		DoAndReturn(func(ctx context.Context, target interface{}, query string, args ...interface{}) error {
			ptr := target.(*[]entity.Exchange)
			*ptr = []entity.Exchange{
				{ID: 1, BaseCurrency: "USD", TargetCurrency: "BRL", Rate: 5.25},
			}
			return nil
		})

	// Act
	result, err := service.ListExchanges(ctx, "USD", "")

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "USD", result[0].BaseCurrency)
}

func TestKsqlExchangeService_ListExchanges_WithTargetFilter(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	service := NewKSQLExchangeService(mockDB)

	ctx := context.Background()
	expectedQuery := "SELECT * FROM exchanges WHERE target_currency = $1"

	mockDB.EXPECT().
		Query(ctx, gomock.Any(), expectedQuery, "BRL").
		DoAndReturn(func(ctx context.Context, target interface{}, query string, args ...interface{}) error {
			ptr := target.(*[]entity.Exchange)
			*ptr = []entity.Exchange{
				{ID: 1, BaseCurrency: "USD", TargetCurrency: "BRL", Rate: 5.25},
				{ID: 2, BaseCurrency: "EUR", TargetCurrency: "BRL", Rate: 5.75},
			}
			return nil
		})

	// Act
	result, err := service.ListExchanges(ctx, "", "BRL")

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "BRL", result[0].TargetCurrency)
	assert.Equal(t, "BRL", result[1].TargetCurrency)
}

func TestKsqlExchangeService_ListExchanges_WithBothFilters(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	service := NewKSQLExchangeService(mockDB)

	ctx := context.Background()
	expectedQuery := "SELECT * FROM exchanges WHERE base_currency = $1 AND target_currency = $2"

	mockDB.EXPECT().
		Query(ctx, gomock.Any(), expectedQuery, "USD", "BRL").
		DoAndReturn(func(ctx context.Context, target interface{}, query string, args ...interface{}) error {
			ptr := target.(*[]entity.Exchange)
			*ptr = []entity.Exchange{
				{ID: 1, BaseCurrency: "USD", TargetCurrency: "BRL", Rate: 5.25},
			}
			return nil
		})

	// Act
	result, err := service.ListExchanges(ctx, "USD", "BRL")

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "USD", result[0].BaseCurrency)
	assert.Equal(t, "BRL", result[0].TargetCurrency)
}

func TestKsqlExchangeService_ListExchanges_Error(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	service := NewKSQLExchangeService(mockDB)

	ctx := context.Background()
	expectedError := errors.New("database error")

	mockDB.EXPECT().
		Query(ctx, gomock.Any(), gomock.Any()).
		Return(expectedError)

	// Act
	result, err := service.ListExchanges(ctx, "", "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}

func TestKsqlExchangeService_ReceiveExchangeRate_CreateNew(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	service := NewKSQLExchangeService(mockDB)

	ctx := context.Background()
	sourceCurrency := "USD"
	targetCurrency := "BRL"
	rate := 5.25

	// Mock getExchangeBySourceAndTarget to return not found
	mockDB.EXPECT().
		QueryOne(ctx, gomock.Any(),
			"SELECT * FROM exchanges WHERE base_currency = $1 AND target_currency = $2 LIMIT 1",
			sourceCurrency, targetCurrency).
		Return(infra.ErrNotFound)

	// Mock createExchange
	mockDB.EXPECT().
		QueryOne(ctx, gomock.Any(),
			"INSERT INTO exchanges (base_currency, target_currency, rate) VALUES ($1, $2, $3) RETURNING id",
			sourceCurrency, targetCurrency, rate).
		DoAndReturn(func(ctx context.Context, target interface{}, query string, args ...interface{}) error {
			ptr := target.(*infra.ReturningID[uint64])
			ptr.ID = 1
			return nil
		})

	// Mock createExchangeRate
	mockDB.EXPECT().
		Exec(ctx, "INSERT INTO exchange_rates (exchange_id, rate) VALUES ($1, $2)", uint64(1), rate).
		Return(mockResult{lastInsertId: 1, rowsAffected: 1}, nil)

	// Act
	err := service.ReceiveExchangeRate(ctx, sourceCurrency, targetCurrency, rate)

	// Assert
	require.NoError(t, err)
}

func TestKsqlExchangeService_ReceiveExchangeRate_UpdateExisting(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	service := NewKSQLExchangeService(mockDB)

	ctx := context.Background()
	sourceCurrency := "USD"
	targetCurrency := "BRL"
	rate := 5.50
	now := time.Now()

	existingExchange := entity.Exchange{
		ID:             1,
		BaseCurrency:   sourceCurrency,
		TargetCurrency: targetCurrency,
		Rate:           5.25,
		CreatedAt:      now.Add(-1 * time.Hour),
		UpdatedAt:      &now,
	}

	// Mock getExchangeBySourceAndTarget to return existing exchange
	mockDB.EXPECT().
		QueryOne(ctx, gomock.Any(),
			"SELECT * FROM exchanges WHERE base_currency = $1 AND target_currency = $2 LIMIT 1",
			sourceCurrency, targetCurrency).
		DoAndReturn(func(ctx context.Context, target interface{}, query string, args ...interface{}) error {
			ptr := target.(*entity.Exchange)
			*ptr = existingExchange
			return nil
		})

	// Mock updateExchange
	mockDB.EXPECT().
		Exec(ctx, "UPDATE exchanges SET rate = $1, updated_at = (now() at TIME ZONE 'UTC') WHERE id = $2",
			rate, uint64(1)).
		Return(mockResult{rowsAffected: 1}, nil)

	// Mock createExchangeRate
	mockDB.EXPECT().
		Exec(ctx, "INSERT INTO exchange_rates (exchange_id, rate) VALUES ($1, $2)", uint64(1), rate).
		Return(mockResult{lastInsertId: 1, rowsAffected: 1}, nil)

	// Act
	err := service.ReceiveExchangeRate(ctx, sourceCurrency, targetCurrency, rate)

	// Assert
	require.NoError(t, err)
}

func TestKsqlExchangeService_ReceiveExchangeRate_CreateExchangeError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	service := NewKSQLExchangeService(mockDB)

	ctx := context.Background()
	sourceCurrency := "USD"
	targetCurrency := "BRL"
	rate := 5.25
	expectedError := errors.New("insert failed")

	// Mock getExchangeBySourceAndTarget to return not found
	mockDB.EXPECT().
		QueryOne(ctx, gomock.Any(),
			"SELECT * FROM exchanges WHERE base_currency = $1 AND target_currency = $2 LIMIT 1",
			sourceCurrency, targetCurrency).
		Return(infra.ErrNotFound)

	// Mock createExchange with error
	mockDB.EXPECT().
		QueryOne(ctx, gomock.Any(),
			"INSERT INTO exchanges (base_currency, target_currency, rate) VALUES ($1, $2, $3) RETURNING id",
			sourceCurrency, targetCurrency, rate).
		Return(expectedError)

	// Act
	err := service.ReceiveExchangeRate(ctx, sourceCurrency, targetCurrency, rate)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}

func TestKsqlExchangeService_ReceiveExchangeRate_UpdateExchangeError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	service := NewKSQLExchangeService(mockDB)

	ctx := context.Background()
	sourceCurrency := "USD"
	targetCurrency := "BRL"
	rate := 5.50
	now := time.Now()
	expectedError := errors.New("update failed")

	existingExchange := entity.Exchange{
		ID:             1,
		BaseCurrency:   sourceCurrency,
		TargetCurrency: targetCurrency,
		Rate:           5.25,
		CreatedAt:      now.Add(-1 * time.Hour),
		UpdatedAt:      &now,
	}

	// Mock getExchangeBySourceAndTarget
	mockDB.EXPECT().
		QueryOne(ctx, gomock.Any(),
			"SELECT * FROM exchanges WHERE base_currency = $1 AND target_currency = $2 LIMIT 1",
			sourceCurrency, targetCurrency).
		DoAndReturn(func(ctx context.Context, target interface{}, query string, args ...interface{}) error {
			ptr := target.(*entity.Exchange)
			*ptr = existingExchange
			return nil
		})

	// Mock updateExchange with error
	mockDB.EXPECT().
		Exec(ctx, "UPDATE exchanges SET rate = $1, updated_at = (now() at TIME ZONE 'UTC') WHERE id = $2",
			rate, uint64(1)).
		Return(nil, expectedError)

	// Act
	err := service.ReceiveExchangeRate(ctx, sourceCurrency, targetCurrency, rate)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}
