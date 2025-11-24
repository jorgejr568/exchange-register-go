// +build integration

package exchange

import (
	"context"
	"github.com/jorgejr568/exchange-register-go/internal/exchange/entity"
	"github.com/jorgejr568/exchange-register-go/internal/infra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const testDBURL = "postgres://exchange:secret@localhost:5432/exchange_test?sslmode=disable"

// setupTestDB creates a fresh database connection for testing
func setupTestDB(t *testing.T) infra.DB {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = testDBURL
	}

	ctx := context.Background()
	db, err := infra.NewKsqlPgDB(ctx, dbURL)
	require.NoError(t, err, "Failed to connect to test database")

	return db
}

// cleanupTestDB cleans up test data
func cleanupTestDB(t *testing.T, db infra.DB) {
	ctx := context.Background()

	// Clean up test data
	_, err := db.Exec(ctx, "DELETE FROM exchange_rates")
	require.NoError(t, err)

	_, err = db.Exec(ctx, "DELETE FROM exchanges")
	require.NoError(t, err)

	err = db.Close()
	require.NoError(t, err)
}

func TestIntegration_ReceiveExchangeRate_CreateNew(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	service := NewKSQLExchangeService(db)
	ctx := context.Background()

	// Act
	err := service.ReceiveExchangeRate(ctx, "USD", "BRL", 5.25)

	// Assert
	require.NoError(t, err)

	// Verify exchange was created
	exchanges, err := service.ListExchanges(ctx, "USD", "BRL")
	require.NoError(t, err)
	require.Len(t, exchanges, 1)
	assert.Equal(t, "USD", exchanges[0].BaseCurrency)
	assert.Equal(t, "BRL", exchanges[0].TargetCurrency)
	assert.Equal(t, 5.25, exchanges[0].Rate)
}

func TestIntegration_ReceiveExchangeRate_UpdateExisting(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	service := NewKSQLExchangeService(db)
	ctx := context.Background()

	// Create initial exchange
	err := service.ReceiveExchangeRate(ctx, "USD", "BRL", 5.25)
	require.NoError(t, err)

	// Act - Update with new rate
	err = service.ReceiveExchangeRate(ctx, "USD", "BRL", 5.50)

	// Assert
	require.NoError(t, err)

	// Verify exchange was updated
	exchanges, err := service.ListExchanges(ctx, "USD", "BRL")
	require.NoError(t, err)
	require.Len(t, exchanges, 1)
	assert.Equal(t, "USD", exchanges[0].BaseCurrency)
	assert.Equal(t, "BRL", exchanges[0].TargetCurrency)
	assert.Equal(t, 5.50, exchanges[0].Rate) // Updated rate
	assert.NotNil(t, exchanges[0].UpdatedAt)
}

func TestIntegration_ListExchanges_NoFilters(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	service := NewKSQLExchangeService(db)
	ctx := context.Background()

	// Create multiple exchanges
	err := service.ReceiveExchangeRate(ctx, "USD", "BRL", 5.25)
	require.NoError(t, err)

	err = service.ReceiveExchangeRate(ctx, "EUR", "BRL", 5.75)
	require.NoError(t, err)

	err = service.ReceiveExchangeRate(ctx, "GBP", "USD", 1.27)
	require.NoError(t, err)

	// Act
	exchanges, err := service.ListExchanges(ctx, "", "")

	// Assert
	require.NoError(t, err)
	assert.Len(t, exchanges, 3)
}

func TestIntegration_ListExchanges_WithSourceFilter(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	service := NewKSQLExchangeService(db)
	ctx := context.Background()

	// Create multiple exchanges
	err := service.ReceiveExchangeRate(ctx, "USD", "BRL", 5.25)
	require.NoError(t, err)

	err = service.ReceiveExchangeRate(ctx, "USD", "EUR", 0.92)
	require.NoError(t, err)

	err = service.ReceiveExchangeRate(ctx, "EUR", "BRL", 5.75)
	require.NoError(t, err)

	// Act
	exchanges, err := service.ListExchanges(ctx, "USD", "")

	// Assert
	require.NoError(t, err)
	assert.Len(t, exchanges, 2)
	for _, exchange := range exchanges {
		assert.Equal(t, "USD", exchange.BaseCurrency)
	}
}

func TestIntegration_ListExchanges_WithTargetFilter(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	service := NewKSQLExchangeService(db)
	ctx := context.Background()

	// Create multiple exchanges
	err := service.ReceiveExchangeRate(ctx, "USD", "BRL", 5.25)
	require.NoError(t, err)

	err = service.ReceiveExchangeRate(ctx, "EUR", "BRL", 5.75)
	require.NoError(t, err)

	err = service.ReceiveExchangeRate(ctx, "GBP", "USD", 1.27)
	require.NoError(t, err)

	// Act
	exchanges, err := service.ListExchanges(ctx, "", "BRL")

	// Assert
	require.NoError(t, err)
	assert.Len(t, exchanges, 2)
	for _, exchange := range exchanges {
		assert.Equal(t, "BRL", exchange.TargetCurrency)
	}
}

func TestIntegration_ListExchanges_WithBothFilters(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	service := NewKSQLExchangeService(db)
	ctx := context.Background()

	// Create multiple exchanges
	err := service.ReceiveExchangeRate(ctx, "USD", "BRL", 5.25)
	require.NoError(t, err)

	err = service.ReceiveExchangeRate(ctx, "USD", "EUR", 0.92)
	require.NoError(t, err)

	err = service.ReceiveExchangeRate(ctx, "EUR", "BRL", 5.75)
	require.NoError(t, err)

	// Act
	exchanges, err := service.ListExchanges(ctx, "USD", "BRL")

	// Assert
	require.NoError(t, err)
	require.Len(t, exchanges, 1)
	assert.Equal(t, "USD", exchanges[0].BaseCurrency)
	assert.Equal(t, "BRL", exchanges[0].TargetCurrency)
	assert.Equal(t, 5.25, exchanges[0].Rate)
}

func TestIntegration_ReceiveExchangeRate_CreatesHistoricalRates(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	service := NewKSQLExchangeService(db)
	ctx := context.Background()

	// Act - Create and update exchange multiple times
	err := service.ReceiveExchangeRate(ctx, "USD", "BRL", 5.25)
	require.NoError(t, err)

	err = service.ReceiveExchangeRate(ctx, "USD", "BRL", 5.30)
	require.NoError(t, err)

	err = service.ReceiveExchangeRate(ctx, "USD", "BRL", 5.35)
	require.NoError(t, err)

	// Assert - Verify historical rates were created
	var rateCount int
	err = db.QueryOne(ctx, &rateCount, "SELECT COUNT(*) as count FROM exchange_rates WHERE exchange_id IN (SELECT id FROM exchanges WHERE base_currency = $1 AND target_currency = $2)", "USD", "BRL")
	require.NoError(t, err)
	assert.Equal(t, 3, rateCount) // Should have 3 historical rates

	// Verify current rate is the latest
	exchanges, err := service.ListExchanges(ctx, "USD", "BRL")
	require.NoError(t, err)
	require.Len(t, exchanges, 1)
	assert.Equal(t, 5.35, exchanges[0].Rate)
}

func TestIntegration_FullWorkflow(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	service := NewKSQLExchangeService(db)
	ctx := context.Background()

	// Act & Assert - Simulate a full sync and list workflow

	// 1. Initially, no exchanges
	exchanges, err := service.ListExchanges(ctx, "", "")
	require.NoError(t, err)
	assert.Len(t, exchanges, 0)

	// 2. Sync some exchanges
	err = service.ReceiveExchangeRate(ctx, "USD", "BRL", 5.25)
	require.NoError(t, err)

	err = service.ReceiveExchangeRate(ctx, "EUR", "BRL", 5.75)
	require.NoError(t, err)

	// 3. List all exchanges
	exchanges, err = service.ListExchanges(ctx, "", "")
	require.NoError(t, err)
	assert.Len(t, exchanges, 2)

	// 4. Update an existing exchange
	err = service.ReceiveExchangeRate(ctx, "USD", "BRL", 5.50)
	require.NoError(t, err)

	// 5. Verify update
	exchanges, err = service.ListExchanges(ctx, "USD", "BRL")
	require.NoError(t, err)
	require.Len(t, exchanges, 1)
	assert.Equal(t, 5.50, exchanges[0].Rate)

	// 6. Filter by target currency
	exchanges, err = service.ListExchanges(ctx, "", "BRL")
	require.NoError(t, err)
	assert.Len(t, exchanges, 2)
}
