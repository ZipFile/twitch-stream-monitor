package cli

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	"github.com/rs/zerolog"

	tsm_app "github.com/ZipFile/twitch-stream-monitor/internal/app"
)

type monitor struct {
	app           tsm_app.App
	loggerFactory func() *zerolog.Logger
}

func NewMonitor(app tsm_app.App) subcommands.Command {
	return &monitor{
		app:           app,
		loggerFactory: emergencyLoggerFactory,
	}
}

func (*monitor) Name() string               { return "monitor" }
func (*monitor) Synopsis() string           { return "Launch Twitch Stream Monitor." }
func (*monitor) Usage() string              { return "monitor" }
func (m *monitor) SetFlags(f *flag.FlagSet) {}

func (m *monitor) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	err := m.app.Init()
	log := m.app.GetLogger()

	if log == nil {
		log = m.loggerFactory()
	}

	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize")
		return subcommands.ExitFailure
	}

	err = m.app.Run()

	if err != nil {
		log.Error().Err(err).Msg("Failed to start")
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
