package migrations

import (
	"context"
	"github.com/jorgejr568/exchange-register-go/internal/infra"
	"github.com/rs/zerolog/log"
)

func CreateExchangesTable(ctx context.Context, db infra.DB) error {
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS exchanges (
			id SERIAL PRIMARY KEY,
			base_currency VARCHAR(3) NOT NULL,
			target_currency VARCHAR(3) NOT NULL,
			rate FLOAT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),
			updated_at TIMESTAMP NULL,
			UNIQUE (base_currency, target_currency)
		);
 	`)

	if err != nil {
		log.Error().Err(err).Msg("failed to create exchanges table")
		return err
	}

	return nil
}
