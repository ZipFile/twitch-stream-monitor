package cli

import (
	"github.com/rs/zerolog"

	tsm_app "github.com/ZipFile/twitch-stream-monitor/internal/app"
)

func emergencyLoggerFactory() *zerolog.Logger {
	log, err := tsm_app.NewLogger("trace", false, true)

	if err != nil {
		panic(err)
	}

	return log
}
