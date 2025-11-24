# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Exchange Register Go is a currency exchange rate tracking service that syncs exchange rates from an external API (FreeCurrency API) and exposes them via a REST API. The application is built using Go 1.19 with a PostgreSQL database.

## Architecture

This project follows a Clean Architecture pattern with clear separation of concerns:

- **cmd/**: CLI commands using Cobra framework
  - `service`: Runs HTTP server (optionally with background sync worker)
  - `migrate`: Database migration runner
  - `sync`: Standalone sync worker for fetching exchange rates

- **internal/exchange/**: Core domain logic
  - `entity/`: Domain entities and use case interfaces
  - `use-cases/`: Business logic implementations (ListExchanges, SyncExchangeRate)
  - `clients/exchangerate/`: External API client implementations (FreeCurrency API)
  - `migrations/`: Database schema definitions
  - `service.go`: Repository implementation using KSQL

- **internal/infra/**: Infrastructure abstractions
  - `database.go`: DB interface wrapping KSQL (vingarcia/ksql with kpgx adapter)
  - Custom error handling (ErrNotFound wraps sql.ErrNoRows)

- **server/**: HTTP server implementation
  - Echo framework v4 for REST API
  - Endpoints: `/status`, `/exchanges`

- **cfg/**: Configuration management
  - Environment variables loaded via godotenv and go-env
  - Singleton pattern for config access via `cfg.Env()`

## Key Design Patterns

- **Use Case Pattern**: Business logic isolated in use case implementations
- **Repository Pattern**: Database access abstracted through `entity.ExchangeService` interface
- **Dependency Injection**: Dependencies passed explicitly to constructors
- **Interface Segregation**: Small, focused interfaces (e.g., `DB`, `ExchangeService`)

## Database

- Uses KSQL ORM (vingarcia/ksql) with PostgreSQL adapter (kpgx)
- Two main tables: `exchanges` and `exchange_rates`
- Manual migrations via `migrate` command (no migration framework)
- Connection pooling configured with max 10 connections

## Common Commands

### Build
```bash
go build -o exchange.bin
```

### Run Database (Docker)
```bash
docker-compose up -d
```

### Database Migration
```bash
go run main.go migrate
```

### Start Service (HTTP server only)
```bash
go run main.go service --port 8080
```

### Start Service with Background Sync Worker
```bash
go run main.go service --sync --port 8080
```

### Run Sync Worker Only
```bash
go run main.go sync
```

### Docker Build
```bash
docker build -t exchange-register-go .
```

## Environment Variables

Required variables (see `.env.example`):
- `DATABASE_URL`: PostgreSQL connection string
- `FREE_CURRENCY_API_KEY`: API key for FreeCurrency API
- `EXCHANGE_CURRENCIES_FROM`: Semicolon-separated list of source currencies (e.g., "USD;EUR;GBP;JPY")
- `EXCHANGE_CURRENCIES_TO`: Semicolon-separated list of target currencies (e.g., "BRL;USD")

Optional variables:
- `HTTP_PORT`: HTTP server port (default: 8080)
- `EXCHANGE_SYNC_SLEEP`: Sync interval duration (default: 30m)

## API Endpoints

### GET /status
Health check endpoint returning `{"status": "ok"}`

### GET /exchanges
List exchange rates with optional filters:
- Query params: `?source=USD&target=BRL` (both optional)
- Returns array of exchange rate objects with current rate and last acquisition time

## Key Implementation Details

- **Sync Process**: Iterates through all combinations of CURRENCIES_FROM × CURRENCIES_TO (excluding same-currency pairs)
- **Exchange Rate Updates**: Creates new exchange if not exists, otherwise updates existing and creates historical rate entry
- **Logging**: Uses zerolog throughout the application
- **Graceful Shutdown**: Both HTTP server and sync worker handle context cancellation for clean shutdown
- **Context Propagation**: Request context flows through all layers (HTTP → Use Case → Service → DB)

## External Dependencies

- **FreeCurrency API**: External exchange rate provider (github.com/jorgejr568/freecurrencyapi-go/v2)
- **KSQL**: Type-safe SQL query builder and ORM
- **Echo**: Web framework for HTTP server
- **Cobra**: CLI command framework
- **Zerolog**: Structured logging
