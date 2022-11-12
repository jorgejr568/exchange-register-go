package cmd

import (
	"github.com/jorgejr568/exchange-register-go/cfg"
	migrations2 "github.com/jorgejr568/exchange-register-go/internal/exchange/migrations"
	"github.com/jorgejr568/exchange-register-go/internal/infra"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate the database",
	Long:  `Migrate the database`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("migrate called")

		ctx := cmd.Context()
		db, err := infra.NewKsqlPgDB(ctx, cfg.Env().DATABASE_URL)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create database connection")
		}

		err = migrations2.CreateExchangesTable(ctx, db)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create exchanges table")
		}

		err = migrations2.CreateExchangeRatesTable(ctx, db)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create exchange_rates table")
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
