package cli

import (
	"context"
	"flag"

	"github.com/google/subcommands"

	tsm_app "github.com/ZipFile/twitch-stream-monitor/internal/app"
)

type unsubscribe struct {
	app            tsm_app.App
	appInitializer AppInitializer
	subIDs         []string
}

func NewUnsubscribe(app tsm_app.App) subcommands.Command {
	return &unsubscribe{
		app:            app,
		appInitializer: DefaultAppInitializer,
	}
}

func (*unsubscribe) Name() string { return "unsubscribe" }
func (*unsubscribe) Synopsis() string {
	return "Unsubscribe from stream.online by subscription id."
}
func (*unsubscribe) Usage() string          { return "unsubscribe SUB_ID ..." }
func (*unsubscribe) SetFlags(*flag.FlagSet) {}

func (u *unsubscribe) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	log, err := u.appInitializer.Init(u.app)

	if err != nil {
		return subcommands.ExitFailure
	}

	svc := u.app.GetTwitchOnlineSubscriptionService()

	if svc == nil {
		log.Error().Msg("Event subscription service was not initialized")

		return subcommands.ExitFailure
	}

	for _, subID := range f.Args() {
		err := svc.Unsubscribe(subID)
		subLog := log.With().Str("subID", subID).Logger()

		if err == nil {
			subLog.Info().Msg("Succesfully unsubscribed")
		} else {
			subLog.Error().Err(err).Msg("Failed to unsubscribe")
		}
	}

	return subcommands.ExitSuccess
}
