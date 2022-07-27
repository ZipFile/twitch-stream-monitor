package cli

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"

	tsm_app "github.com/ZipFile/twitch-stream-monitor/internal/app"
)

func emergencyLoggerFactory() *zerolog.Logger {
	log := zerolog.New(os.Stderr).Level(zerolog.TraceLevel).With().Timestamp().Logger()

	return &log
}

func simplePrintln(s string) error {
	_, err := fmt.Println(s)

	return err
}

type AppInitializer interface {
	Init(tsm_app.App) (zerolog.Logger, error)
}

type appInitializer struct {
	loggerFactory func() *zerolog.Logger
}

var DefaultAppInitializer AppInitializer = &appInitializer{
	loggerFactory: emergencyLoggerFactory,
}

func (ai *appInitializer) Init(app tsm_app.App) (zerolog.Logger, error) {
	err := app.Init()
	log := app.GetLogger()

	if log == nil {
		log = ai.loggerFactory()
	}

	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize")
	}

	return *log, err
}
