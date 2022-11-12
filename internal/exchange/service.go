package exchange

import (
	"context"
	"errors"
	"github.com/jorgejr568/exchange-register-go/internal/exchange/entity"
	"github.com/jorgejr568/exchange-register-go/internal/infra"
	"github.com/rs/zerolog/log"
)

type ksqlExchangeService struct {
	db infra.DB
}

func (k ksqlExchangeService) ReceiveExchangeRate(ctx context.Context, sourceCurrency, targetCurrency string, rate float64) error {
	exchange, err := k.getExchangeBySourceAndTarget(ctx, sourceCurrency, targetCurrency)
	if err != nil {
		if errors.Is(err, infra.ErrNotFound) {
			createdID, err := k.createExchange(ctx, sourceCurrency, targetCurrency, rate)
			if err != nil {
				log.Error().Err(err).Msg("failed to create exchange")
				return err
			}
			log.Debug().Msgf("created exchange with id %d", createdID)
			err = k.createExchangeRate(ctx, createdID, rate)
			if err != nil {
				log.Error().Err(err).Msg("failed to create exchange rate")
				return err
			}

			log.Debug().Msgf("created exchange rate for exchange with id %d", createdID)
			return nil
		}

		return err
	}

	err = k.updateExchange(ctx, exchange.ID, rate)
	if err != nil {
		log.Error().Err(err).Msgf("failed to update exchange with id %d", exchange.ID)
		return err
	}
	log.Debug().Msgf("updated exchange with id %d", exchange.ID)

	err = k.createExchangeRate(ctx, exchange.ID, rate)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create exchange rate for exchange with id %d", exchange.ID)
		return err
	}
	log.Debug().Msgf("created exchange rate for exchange with id %d", exchange.ID)

	return nil
}

func (k ksqlExchangeService) createExchangeRate(ctx context.Context, id uint64, rate float64) error {
	_, err := k.db.Exec(ctx, `INSERT INTO exchange_rates (exchange_id, rate) VALUES ($1, $2)`, id, rate)
	if err != nil {
		return err
	}

	return nil
}

func (k ksqlExchangeService) updateExchange(ctx context.Context, id uint64, rate float64) error {
	_, err := k.db.Exec(ctx, `UPDATE exchanges SET rate = $1, updated_at = (now() at TIME ZONE 'UTC') WHERE id = $2`, rate, id)
	if err != nil {
		return err
	}

	return nil
}

func (k ksqlExchangeService) createExchange(ctx context.Context, sourceCurrency, targetCurrency string, rate float64) (uint64, error) {
	var returningResult infra.ReturningID[uint64]
	err := k.db.QueryOne(ctx, &returningResult, `INSERT INTO exchanges (base_currency, target_currency, rate) VALUES ($1, $2, $3) RETURNING id`, sourceCurrency, targetCurrency, rate)
	if err != nil {
		return 0, err
	}

	return returningResult.ID, nil
}

func (k ksqlExchangeService) getExchangeBySourceAndTarget(ctx context.Context, sourceCurrency, targetCurrency string) (entity.Exchange, error) {
	var exchange entity.Exchange
	err := k.db.QueryOne(ctx, &exchange, `SELECT * FROM exchanges WHERE base_currency = $1 AND target_currency = $2 LIMIT 1`, sourceCurrency, targetCurrency)
	if err != nil {
		return entity.Exchange{}, err
	}

	return exchange, nil
}

func NewKSQLExchangeService(db infra.DB) entity.ExchangeService {
	return &ksqlExchangeService{
		db: db,
	}
}
