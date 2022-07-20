package main

import (
	_ "github.com/joho/godotenv/autoload"

	tsm_app "github.com/ZipFile/twitch-stream-monitor/internal/app"
)

func main() {
	app := tsm_app.New()
	err := app.Init()
	log := app.GetLogger()

	if log == nil {
		log, _ = tsm_app.NewLogger("trace", false, true)
	}

	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize")
		return
	}

	err = app.Run()

	if err != nil {
		log.Error().Err(err).Msg("Failed to start")
	}
}
