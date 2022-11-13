package cmd

import (
	"context"
	"github.com/jorgejr568/exchange-register-go/cfg"
	"github.com/jorgejr568/exchange-register-go/internal/exchange"
	use_cases "github.com/jorgejr568/exchange-register-go/internal/exchange/use-cases"
	"github.com/jorgejr568/exchange-register-go/internal/infra"
	"github.com/jorgejr568/exchange-register-go/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Starts the http server and the sync process on the background",
	Long:  `Starts the http server and the sync process on the background`,
	Run: func(cmd *cobra.Command, args []string) {
		syncWorkerEnabled, err := cmd.Flags().GetBool("sync")
		if err != nil {
			log.Error().Err(err).Msg("failed to get sync flag")
			return
		}

		port, err := cmd.Flags().GetString("port")
		if err != nil {
			log.Error().Err(err).Msg("failed to get port flag")
			return
		}

		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		go func() {
			<-quit
			cancel()
		}()

		db, err := infra.NewKsqlPgDB(ctx, cfg.Env().DATABASE_URL)
		if err != nil {
			log.Panic().Err(err).Msg("failed to connect to database")
		}

		service := exchange.NewKSQLExchangeService(db)
		listExchangesUseCase := use_cases.NewListExchangesUseCase(service)
		s := server.NewEchoServer(listExchangesUseCase, port)

		if syncWorkerEnabled {
			go syncCmd.Run(cmd, args)
		}

		go func() {
			log.Info().Msgf("exchange-register-go service running on port %s... (press Ctrl+C to quit)", cfg.Env().HTTP_PORT)
			<-ctx.Done()
			log.Info().Msg("exchange-register-go service stopped")
		}()
		if err := s.GracefulListenAndShutdown(ctx); err != nil {
			log.Error().Err(err).Msg("exchange-register-go service stopped with error")
		}
	},
}

func init() {
	serviceCmd.Flags().StringP("port", "p", cfg.Env().HTTP_PORT, "http server port")
	serviceCmd.Flags().BoolP("sync", "s", false, "sync worker enabled")
	rootCmd.AddCommand(serviceCmd)
}
