package testing

import (
	"github.com/rs/zerolog"
)

func NoopLoggerFactory() *zerolog.Logger {
	log := zerolog.Nop()

	return &log
}

func PanicLoggerFactory() *zerolog.Logger {
	panic("should not be called")

	return nil
}
