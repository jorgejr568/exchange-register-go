package main

import (
	"github.com/jorgejr568/exchange-register-go/cmd"
	"github.com/rs/zerolog/log"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Error().Err(err).Msg("failed to execute command")
	}
}
