package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/jorgejr568/exchange-register-go/internal/exchange/entity"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
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

	e.GET("/status", s.statusHandler)
	e.GET("/exchanges", s.exchangesHandler)
	e.GET("/openapi.json", s.openapiHandler)
	e.GET("/docs", s.docsHandler)

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

// Helper method to extract handler functions from the server
func (s *echoServer) statusHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (s *echoServer) exchangesHandler(c echo.Context) error {
	ctx := c.Request().Context()
	req := entity.ListExchangesRequest{
		SourceCurrency: c.QueryParam("source"),
		TargetCurrency: c.QueryParam("target"),
	}

	res, err := s.listExchangesUseCase.Execute(ctx, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to list exchanges",
		})
	}

	return c.JSON(http.StatusOK, res)
}

func (s *echoServer) docsHandler(c echo.Context) error {
	return c.File("static/docs/index.html")
}

func (s *echoServer) openapiHandler(c echo.Context) error {
	serverURL := getServerURL(c)
	spec, err := GenerateOpenAPISpec(serverURL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to generate openapi spec",
		})
	}
	return c.JSON(http.StatusOK, spec)
}
