package cmd

import (
	"context"
	"github.com/jorgejr568/exchange-register-go/cfg"
	"github.com/jorgejr568/exchange-register-go/internal/exchange"
	"github.com/jorgejr568/exchange-register-go/internal/exchange/clients/exchangerate"
	"github.com/jorgejr568/exchange-register-go/internal/exchange/entity"
	"github.com/jorgejr568/exchange-register-go/internal/exchange/use-cases"
	"github.com/jorgejr568/exchange-register-go/internal/infra"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"time"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Syncs the exchange rates from the external API",
	Long:  `Syncs the exchange rates from the external API`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		db, err := infra.NewKsqlPgDB(ctx, cfg.Env().DATABASE_URL)
		if err != nil {
			log.Error().Err(err).Msg("failed to create db")
			return
		}

		exchangeService := exchange.NewKSQLExchangeService(
			db,
		)
		exchangeRateClient := exchangerate.NewFreeCurrencyApiClient(
			cfg.Env().FreeCurrencyAPIClient(),
		)

		useCase := use_cases.NewSyncExchangeRateUseCase(
			exchangeService,
			exchangeRateClient,
		)
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		runSync(ctx, useCase)
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		go func() {
			log.Info().Msg("exchange-register-go sync running... (press Ctrl+C to quit)")
			<-quit
			cancel()
		}()

		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("exchange-register-go sync stopped")
				err := db.Close()
				if err != nil {
					log.Error().Err(err).Msg("failed to close db")
				}
				return
			case <-time.After(cfg.Env().EXCHANGE_SYNC_SLEEP):
				runSync(ctx, useCase)
			}
		}
	},
}

func runSync(ctx context.Context, useCase entity.SyncExchangeRateUseCase) {
	currenciesFrom := cfg.Env().CurrenciesFrom()
	currenciesTo := cfg.Env().CurrenciesTo()
	for _, from := range currenciesFrom {
		for _, to := range currenciesTo {
			if from == to {
				continue
			}

			_, err := useCase.Execute(ctx, entity.SyncExchangeRateRequest{
				SourceCurrency: from,
				TargetCurrency: to,
			})
			if err != nil {
				log.Error().Err(err).Msg("failed to sync exchange rate")
			}
		}
	}
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
