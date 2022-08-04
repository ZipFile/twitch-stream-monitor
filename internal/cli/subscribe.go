package cli

import (
	"context"
	"flag"
	"regexp"

	"github.com/google/subcommands"

	tsm_app "github.com/ZipFile/twitch-stream-monitor/internal/app"
)

type subscribe struct {
	app            tsm_app.App
	appInitializer AppInitializer
}

func NewSubscribe(app tsm_app.App) subcommands.Command {
	return &subscribe{
		app:            app,
		appInitializer: DefaultAppInitializer,
	}
}

func (*subscribe) Name() string { return "subscribe" }
func (*subscribe) Synopsis() string {
	return "Subscribe to stream.online events for given twitch users."
}
func (*subscribe) Usage() string          { return "subscribe BROADCASTER_ID ..." }
func (*subscribe) SetFlags(*flag.FlagSet) {}

var numericIdRegexp = regexp.MustCompile(`^\d+$`)

func (s *subscribe) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	var err error
	log, err := s.appInitializer.Init(s.app)

	if err != nil {
		return subcommands.ExitFailure
	}

	svc := s.app.GetTwitchOnlineSubscriptionService()

	if svc == nil {
		log.Error().Msg("Event subscription service was not initialized")

		return subcommands.ExitFailure
	}

	ur := s.app.GetUsernameResolver()

	if ur == nil {
		log.Error().Msg("Username resolver failed to initialized")

		return subcommands.ExitFailure
	}

	for _, broadcasterID := range f.Args() {
		if !numericIdRegexp.MatchString(broadcasterID) {
			broadcasterID, err = ur.Resolve(broadcasterID)
		}

		if err != nil {
			log.Error().Err(err).Msg("Failed to resolve username")

			return subcommands.ExitFailure
		}

		subID, err := svc.Subscribe(broadcasterID)
		subLog := log.With().Str("subID", subID).Str("broadcasterID", broadcasterID).Logger()

		if err == nil {
			subLog.Info().Msg("Succesfully subscribed")
		} else {
			subLog.Error().Err(err).Msg("Failed to subscribe")
		}
	}

	return subcommands.ExitSuccess
}
