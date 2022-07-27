package cli

import (
	"context"
	"flag"

	"github.com/google/subcommands"

	tsm_app "github.com/ZipFile/twitch-stream-monitor/internal/app"
)

type monitor struct {
	app            tsm_app.App
	appInitializer AppInitializer
}

func NewMonitor(app tsm_app.App) subcommands.Command {
	return &monitor{
		app:            app,
		appInitializer: DefaultAppInitializer,
	}
}

func (*monitor) Name() string           { return "monitor" }
func (*monitor) Synopsis() string       { return "Launch Twitch Stream Monitor." }
func (*monitor) Usage() string          { return "monitor" }
func (*monitor) SetFlags(*flag.FlagSet) {}

func (m *monitor) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	log, err := m.appInitializer.Init(m.app)

	if err != nil {
		return subcommands.ExitFailure
	}

	err = m.app.Run()

	if err != nil {
		log.Error().Err(err).Msg("Failed to start")
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
