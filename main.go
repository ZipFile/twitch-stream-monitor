package main

import (
	"time"

	_ "github.com/joho/godotenv/autoload"

	tsm_app "github.com/ZipFile/twitch-stream-monitor/internal/app"
	tsm_config "github.com/ZipFile/twitch-stream-monitor/internal/config"
)

func main() {
	app := tsm_app.App{
		Now:          time.Now,
		ConfigLoader: tsm_config.NewEnvironLoader(nil),
	}
	err := app.Init()

	if err != nil {
		if app.Log == nil {
			panic(err)
		} else {
			app.Log.Error().Err(err).Msg("Failed to initialize")
			return
		}
	}

	err = app.Monitor()

	if err != nil {
		app.Log.Error().Err(err).Msg("Failed to start")
	}
}
