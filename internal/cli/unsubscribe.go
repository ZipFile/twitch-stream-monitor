package cli

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	"github.com/rs/zerolog"

	tsm_app "github.com/ZipFile/twitch-stream-monitor/internal/app"
)

type unsubscribe struct {
	app           tsm_app.App
	loggerFactory func() *zerolog.Logger
	subIDs        []string
}

func NewUnsubscribe(app tsm_app.App) subcommands.Command {
	return &unsubscribe{
		app:           app,
		loggerFactory: emergencyLoggerFactory,
	}
}

func (*unsubscribe) Name() string { return "unsubscribe" }
func (*unsubscribe) Synopsis() string {
	return "Unsubscribe from stream.online by subscription id."
}
func (*unsubscribe) Usage() string { return "unsubscribe SUB_ID ..." }
func (u *unsubscribe) SetFlags(f *flag.FlagSet) {
	u.subIDs = f.Args()
}

func (u *unsubscribe) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	err := u.app.Init()
	log := u.app.GetLogger()

	if log == nil {
		log = u.loggerFactory()
	}

	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize")
		return subcommands.ExitFailure
	}

	svc := u.app.GetTwitchOnlineSubscriptionService()

	if svc == nil {
		log.Error().Msg("Event subscription service was not initialized")

		return subcommands.ExitFailure
	}

	for _, subID := range u.subIDs {
		err := svc.Unsubscribe(subID)
		subLog := log.With().Str("subID", subID).Logger()

		if err == nil {
			subLog.Error().Err(err).Msg("Failed to unsubscribe")
		} else {
			subLog.Info().Msg("Succesfully unsubscribed")
		}
	}

	return subcommands.ExitSuccess
}
