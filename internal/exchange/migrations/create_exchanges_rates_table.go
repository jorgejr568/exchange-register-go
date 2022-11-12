package migrations

import (
	"context"
	"github.com/jorgejr568/exchange-register-go/internal/infra"
	"github.com/rs/zerolog/log"
)

func CreateExchangeRatesTable(ctx context.Context, db infra.DB) error {
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS exchange_rates (
			id SERIAL PRIMARY KEY,
			exchange_id BIGINT NOT NULL,
			rate FLOAT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),
			FOREIGN KEY (exchange_id) REFERENCES exchanges(id) ON DELETE CASCADE
		);
 	`)

	if err != nil {
		log.Error().Err(err).Msg("failed to create exchange_rates table")
		return err
	}

	return nil
}
