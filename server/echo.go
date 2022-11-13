package server

import (
	"context"
	"errors"
	"github.com/jorgejr568/exchange-register-go/internal/exchange/entity"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"net/http"
)

type echoServer struct {
	listExchangesUseCase entity.ListExchangesUseCase
	httpPort             string
}

func (s *echoServer) GracefulListenAndShutdown(ctx context.Context) error {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/status", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	e.GET("/exchanges", func(c echo.Context) error {
		req := entity.ListExchangesRequest{
			SourceCurrency: c.QueryParam("source"),
			TargetCurrency: c.QueryParam("target"),
		}

		res, err := s.listExchangesUseCase.Execute(ctx, req)
		if err != nil {
			log.Error().Err(err).Msg("failed to list exchanges")
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "failed to list exchanges",
			})
		}

		return c.JSON(http.StatusOK, res)
	})

	go func() {
		<-ctx.Done()
		err := e.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed to close echo server")
		}
	}()

	err := e.Start(":" + s.httpPort)
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}

	return nil
}

func NewEchoServer(listExchangesUseCase entity.ListExchangesUseCase, httpPort string) Server {
	return &echoServer{
		listExchangesUseCase: listExchangesUseCase,
		httpPort:             httpPort,
	}
}
