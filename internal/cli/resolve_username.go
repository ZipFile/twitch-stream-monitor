package cli

import (
	"context"
	"flag"

	"github.com/google/subcommands"

	tsm_app "github.com/ZipFile/twitch-stream-monitor/internal/app"
)

type resolveUsername struct {
	app            tsm_app.App
	appInitializer AppInitializer
	print          func(string) error
}

func NewResolveUsername(app tsm_app.App) subcommands.Command {
	return &resolveUsername{
		app:            app,
		appInitializer: DefaultAppInitializer,
		print:          simplePrintln,
	}
}

func (*resolveUsername) Name() string { return "resolve-username" }
func (*resolveUsername) Synopsis() string {
	return "Find twitch user id by its username."
}
func (*resolveUsername) Usage() string          { return "resolve-username USERNAME" }
func (*resolveUsername) SetFlags(*flag.FlagSet) {}

func (ru *resolveUsername) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	log, err := ru.appInitializer.Init(ru.app)

	if err != nil {
		return subcommands.ExitFailure
	}

	if f.NArg() != 1 {
		log.Error().Msg("Username is required")

		return subcommands.ExitUsageError
	}

	ur := ru.app.GetUsernameResolver()

	if ur == nil {
		log.Error().Msg("Username resolver failed to initialized")

		return subcommands.ExitFailure
	}

	username := f.Args()[0]
	broadcasterID, err := ur.Resolve(username)

	if err != nil {
		log.Error().Err(err).Msg("Failed to resolve username")

		return subcommands.ExitFailure
	}

	if broadcasterID == "" {
		return subcommands.ExitSuccess
	}

	err = ru.print(broadcasterID)

	if err != nil {
		log.Error().Err(err).Msg("Failed to print broadcaster id")

		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
