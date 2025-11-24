package server

import (
	"github.com/jorgejr568/exchange-register-go/internal/exchange/entity"
	"github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi3"
	"net/http"
)

// StatusResponse represents the health check response
type StatusResponse struct {
	Status string `json:"status" example:"ok"`
}

// ListExchangesQueryParams represents query parameters for listing exchanges
type ListExchangesQueryParams struct {
	Source string `query:"source" description:"Source currency code (e.g., USD, EUR)" example:"USD"`
	Target string `query:"target" description:"Target currency code (e.g., BRL, EUR)" example:"BRL"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"failed to list exchanges"`
}

// GenerateOpenAPISpec creates the OpenAPI 3.0 specification for the API
func GenerateOpenAPISpec() (*openapi3.Spec, error) {
	reflector := openapi3.Reflector{}
	reflector.Spec = &openapi3.Spec{Openapi: "3.0.3"}
	reflector.Spec.Info.
		WithTitle("Exchange Register API").
		WithVersion("1.0.0").
		WithDescription("Currency exchange rate tracking service that syncs exchange rates from external APIs and exposes them via REST endpoints")

	// Add server
	reflector.Spec.Servers = []openapi3.Server{
		{
			URL:         "http://localhost:8080",
			Description: stringPtr("Development server"),
		},
	}

	// Add tags
	reflector.SpecEns().WithTags(
		openapi3.Tag{
			Name:        "Health",
			Description: stringPtr("Health check endpoints"),
		},
		openapi3.Tag{
			Name:        "Exchanges",
			Description: stringPtr("Exchange rate management"),
		},
		openapi3.Tag{
			Name:        "Documentation",
			Description: stringPtr("API documentation endpoints"),
		},
	)

	// GET /status endpoint
	statusOp, err := reflector.NewOperationContext(http.MethodGet, "/status")
	if err != nil {
		return nil, err
	}
	statusOp.SetSummary("Health check")
	statusOp.SetDescription("Returns the health status of the service")
	statusOp.SetTags("Health")
	statusOp.AddRespStructure(new(StatusResponse), func(cu *openapi.ContentUnit) {
		cu.HTTPStatus = http.StatusOK
	})
	if err := reflector.AddOperation(statusOp); err != nil {
		return nil, err
	}

	// GET /exchanges endpoint
	exchangesOp, err := reflector.NewOperationContext(http.MethodGet, "/exchanges")
	if err != nil {
		return nil, err
	}
	exchangesOp.SetSummary("List exchange rates")
	exchangesOp.SetDescription("Retrieves a list of exchange rates, optionally filtered by source and/or target currency")
	exchangesOp.SetTags("Exchanges")
	exchangesOp.AddReqStructure(new(ListExchangesQueryParams))
	exchangesOp.AddRespStructure(new(entity.ListExchangesResponse), func(cu *openapi.ContentUnit) {
		cu.HTTPStatus = http.StatusOK
	})
	exchangesOp.AddRespStructure(new(ErrorResponse), func(cu *openapi.ContentUnit) {
		cu.HTTPStatus = http.StatusInternalServerError
	})
	if err := reflector.AddOperation(exchangesOp); err != nil {
		return nil, err
	}

	// GET /openapi.json endpoint (self-documenting)
	openAPIOp, err := reflector.NewOperationContext(http.MethodGet, "/openapi.json")
	if err != nil {
		return nil, err
	}
	openAPIOp.SetSummary("OpenAPI specification")
	openAPIOp.SetDescription("Returns the OpenAPI specification for this API")
	openAPIOp.SetTags("Documentation")
	openAPIOp.AddRespStructure(new(map[string]interface{}), func(cu *openapi.ContentUnit) {
		cu.HTTPStatus = http.StatusOK
		cu.Description = "OpenAPI specification"
	})
	if err := reflector.AddOperation(openAPIOp); err != nil {
		return nil, err
	}

	return reflector.Spec, nil
}

func stringPtr(s string) *string {
	return &s
}
