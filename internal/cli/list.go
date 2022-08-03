package cli

import (
	"context"
	"flag"

	"github.com/google/subcommands"

	tsm_app "github.com/ZipFile/twitch-stream-monitor/internal/app"
)

type list struct {
	app            tsm_app.App
	appInitializer AppInitializer
}

func NewList(app tsm_app.App) subcommands.Command {
	return &list{
		app:            app,
		appInitializer: DefaultAppInitializer,
	}
}

func (*list) Name() string           { return "list" }
func (*list) Synopsis() string       { return "List stream.online subscriptions." }
func (*list) Usage() string          { return "list" }
func (*list) SetFlags(*flag.FlagSet) {}

func (l *list) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	log, err := l.appInitializer.Init(l.app)

	if err != nil {
		return subcommands.ExitFailure
	}

	svc := l.app.GetTwitchOnlineSubscriptionService()

	if svc == nil {
		log.Error().Msg("Event subscription service was not initialized")

		return subcommands.ExitFailure
	}

	subs, err := svc.List()

	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve list of stream.online subscriptions")
		return subcommands.ExitFailure
	}

	for _, sub := range subs {
		log.Info().Str("id", sub.ID).
			Str("status", sub.Status).
			Str("user_id", sub.UserID).
			Str("callback_url", sub.CallbackURL).
			Send()
	}

	return subcommands.ExitSuccess
}
