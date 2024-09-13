package main

import (
	"Audio2TextService/cmd/log"
	"Audio2TextService/internal/pkg/app"
)

func main() {

	logger := log.New()

	app := app.New(&logger)

	logger.Info().Msg("Starting application...")
	err := app.Run()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to start application")
	}
}
