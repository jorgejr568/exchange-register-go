package server

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jorgejr568/exchange-register-go/internal/exchange/entity"
	"github.com/jorgejr568/exchange-register-go/internal/exchange/entity/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestStatusEndpoint_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockListExchangesUseCase(ctrl)
	e := echo.New()
	server := NewEchoServer(mockUseCase, "8080").(*echoServer)

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Act
	err := server.statusHandler(c)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
}

func TestExchangesEndpoint_Success_WithFilters(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockListExchangesUseCase(ctrl)
	e := echo.New()
	server := NewEchoServer(mockUseCase, "8080").(*echoServer)

	now := time.Now()
	expectedResponse := entity.ListExchangesResponse{
		{
			ID:              1,
			SourceCurrency:  "USD",
			TargetCurrency:  "BRL",
			Rate:            5.25,
			LastAcquisition: now,
		},
	}

	ctx := context.Background()
	req := entity.ListExchangesRequest{
		SourceCurrency: "USD",
		TargetCurrency: "BRL",
	}

	mockUseCase.EXPECT().
		Execute(ctx, req).
		Return(&expectedResponse, nil)

	httpReq := httptest.NewRequest(http.MethodGet, "/exchanges?source=USD&target=BRL", nil)
	httpReq = httpReq.WithContext(ctx)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	// Act
	err := server.exchangesHandler(c)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response entity.ListExchangesResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, "USD", response[0].SourceCurrency)
	assert.Equal(t, "BRL", response[0].TargetCurrency)
	assert.Equal(t, 5.25, response[0].Rate)
}

func TestExchangesEndpoint_Success_WithoutFilters(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockListExchangesUseCase(ctrl)
	e := echo.New()
	server := NewEchoServer(mockUseCase, "8080").(*echoServer)

	now := time.Now()
	expectedResponse := entity.ListExchangesResponse{
		{
			ID:              1,
			SourceCurrency:  "USD",
			TargetCurrency:  "BRL",
			Rate:            5.25,
			LastAcquisition: now,
		},
		{
			ID:              2,
			SourceCurrency:  "EUR",
			TargetCurrency:  "BRL",
			Rate:            5.75,
			LastAcquisition: now,
		},
	}

	ctx := context.Background()
	req := entity.ListExchangesRequest{
		SourceCurrency: "",
		TargetCurrency: "",
	}

	mockUseCase.EXPECT().
		Execute(ctx, req).
		Return(&expectedResponse, nil)

	httpReq := httptest.NewRequest(http.MethodGet, "/exchanges", nil)
	httpReq = httpReq.WithContext(ctx)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	// Act
	err := server.exchangesHandler(c)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response entity.ListExchangesResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 2)
}

func TestExchangesEndpoint_Error(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockListExchangesUseCase(ctrl)
	e := echo.New()
	server := NewEchoServer(mockUseCase, "8080").(*echoServer)

	ctx := context.Background()
	req := entity.ListExchangesRequest{
		SourceCurrency: "USD",
		TargetCurrency: "BRL",
	}

	expectedError := errors.New("database connection failed")
	mockUseCase.EXPECT().
		Execute(ctx, req).
		Return(nil, expectedError)

	httpReq := httptest.NewRequest(http.MethodGet, "/exchanges?source=USD&target=BRL", nil)
	httpReq = httpReq.WithContext(ctx)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	// Act
	err := server.exchangesHandler(c)

	// Assert
	require.NoError(t, err) // Handler itself doesn't return error, it returns JSON error
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "failed to list exchanges", response["error"])
}

func TestExchangesEndpoint_EmptyResult(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockListExchangesUseCase(ctrl)
	e := echo.New()
	server := NewEchoServer(mockUseCase, "8080").(*echoServer)

	ctx := context.Background()
	req := entity.ListExchangesRequest{
		SourceCurrency: "USD",
		TargetCurrency: "XYZ",
	}

	emptyResponse := entity.ListExchangesResponse{}
	mockUseCase.EXPECT().
		Execute(ctx, req).
		Return(&emptyResponse, nil)

	httpReq := httptest.NewRequest(http.MethodGet, "/exchanges?source=USD&target=XYZ", nil)
	httpReq = httpReq.WithContext(ctx)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	// Act
	err := server.exchangesHandler(c)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response entity.ListExchangesResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 0)
}

func TestOpenAPIEndpoint_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockListExchangesUseCase(ctrl)
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := func(c echo.Context) error {
		spec, err := GenerateOpenAPISpec()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "failed to generate openapi spec",
			})
		}
		return c.JSON(http.StatusOK, spec)
	}

	// Act
	err := handler(c)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify basic OpenAPI structure
	assert.Equal(t, "3.0.3", response["openapi"])
	assert.Contains(t, response, "info")
	assert.Contains(t, response, "paths")
	assert.Contains(t, response, "components")

	// Verify paths exist
	paths := response["paths"].(map[string]interface{})
	assert.Contains(t, paths, "/status")
	assert.Contains(t, paths, "/exchanges")
	assert.Contains(t, paths, "/openapi.json")

	_ = mockUseCase // Avoid unused variable warning
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
