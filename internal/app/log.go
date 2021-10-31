package app

import (
	"io"
	"os"
	"runtime"
	"time"

	"github.com/rs/zerolog"
)

func NewLogger(level string, pretty, stdout bool) (*zerolog.Logger, error) {
	logLevel, err := zerolog.ParseLevel(level)

	if err != nil {
		return nil, err
	}

	var w io.Writer

	if stdout {
		w = os.Stdout
	} else {
		w = os.Stderr
	}

	if pretty {
		w = zerolog.ConsoleWriter{
			NoColor:    runtime.GOOS == "windows",
			Out:        w,
			TimeFormat: time.RFC3339,
		}
	}

	log := zerolog.New(w).Level(logLevel).With().Timestamp().Logger()

	return &log, nil
}
