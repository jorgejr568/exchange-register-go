package cfg

import (
	goenv "github.com/Netflix/go-env"
	"github.com/joho/godotenv"
	"github.com/jorgejr568/freecurrencyapi-go/v2"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

type EnvironmentVariables struct {
	DATABASE_URL             string        `env:"DATABASE_URL,required=true"`
	HTTP_PORT                string        `env:"HTTP_PORT,default=8080"`
	EXCHANGE_RATE_API_URL    string        `env:"EXCHANGE_RATE_API_URL"`
	EXCHANGE_SYNC_SLEEP      time.Duration `env:"EXCHANGE_SYNC_SLEEP,default=30m"`
	EXCHANGE_CURRENCIES_FROM string        `env:"EXCHANGE_CURRENCIES_FROM,default=USD;EUR;GBP;JPY"`
	EXCHANGE_CURRENCIES_TO   string        `env:"EXCHANGE_CURRENCIES_TO,default=BRL"`
	FREE_CURRENCY_API_KEY    string        `env:"FREE_CURRENCY_API_KEY,required=true"`
}

var _env *EnvironmentVariables

func Env() *EnvironmentVariables {
	if _env == nil {
		err := godotenv.Load()
		if err != nil {
			log.Error().Err(err).Msg("failed to load .env file")
		}

		var env EnvironmentVariables
		_, err = goenv.UnmarshalFromEnviron(&env)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to unmarshal environment variables")
		}

		_env = &env
	}

	return _env
}

func (e *EnvironmentVariables) CurrenciesFrom() []string {
	return strings.Split(e.EXCHANGE_CURRENCIES_FROM, ";")
}

func (e *EnvironmentVariables) CurrenciesTo() []string {
	return strings.Split(e.EXCHANGE_CURRENCIES_TO, ";")
}

func (e *EnvironmentVariables) FreeCurrencyAPIClient() freecurrencyapi.Client {
	return freecurrencyapi.NewClient(e.FREE_CURRENCY_API_KEY)

}
