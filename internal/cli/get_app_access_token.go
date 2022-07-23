package cli

import (
	"context"
	"flag"
	"time"

	"github.com/google/subcommands"

	tsm_app "github.com/ZipFile/twitch-stream-monitor/internal/app"
)

type getAppAccessToken struct {
	app            tsm_app.App
	appInitializer AppInitializer
	now            func() time.Time
	print          func(string) error
}

func NewGetAppAccessToken(app tsm_app.App) subcommands.Command {
	return &getAppAccessToken{
		app:            app,
		appInitializer: DefaultAppInitializer,
		now:            time.Now,
		print:          simplePrintln,
	}
}

func (*getAppAccessToken) Name() string { return "get-app-access-token" }
func (*getAppAccessToken) Synopsis() string {
	return "Retrieve fresh app access token to use in API queries."
}
func (*getAppAccessToken) Usage() string                 { return "get-app-access-token" }
func (gaat *getAppAccessToken) SetFlags(f *flag.FlagSet) {}

func (gaat *getAppAccessToken) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	log, err := gaat.appInitializer.Init(gaat.app)

	if err != nil {
		return subcommands.ExitFailure
	}

	ts := gaat.app.GetTokenStore()

	if ts == nil {
		log.Error().Msg("Token store failed to initialized")

		return subcommands.ExitFailure
	}

	token, expires, err := ts.Load()

	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve app access token")

		return subcommands.ExitFailure
	}

	if expires.Before(gaat.now()) {
		log.Error().Err(err).Msg("App access token is expired")

		return subcommands.ExitFailure
	}

	err = gaat.print(token)

	if err != nil {
		log.Error().Err(err).Msg("Failed to print app access token")

		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
