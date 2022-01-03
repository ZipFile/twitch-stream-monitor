package tsm

import (
	"context"
	"errors"
	"sync"

	"github.com/rs/zerolog"

	tsm "github.com/ZipFile/twitch-stream-monitor/internal"
)

type Settings struct {
	// List of Twitch user (broadcaster) ids. Not usernames.
	UserIDs []string
	// Do not delete existing subscriptions at exit.
	KeepExistingSubs bool
	// Do not delete created subscriptions at exit.
	KeepNewSubs bool
	// Keep going if online event handler is nonoperational.
	IgnoreStartupErrors bool
	// Keep going if subscription creation fails.
	IgnoreSubscriptionErrors bool
}

// Main logic for monitor sub-command.
//
// Performs check if handler is functional, subscribes to the each user and then
// starts listening to online events. Unsubscribes on shutdown.
//
// Args:
// * ctx: Context used for graceful shutdown. Passed directly to svc.Listen and into the handler.Handle.
// * svc: Instantiate of the online subscription instance.
// * handler: Online event handler.
// * settings: Monitoring parameters.
// * log: Logger instance.
func Monitor(
	ctx context.Context,
	svc tsm.TwitchOnlineSubscriptionService,
	handler tsm.TwitchStreamOnlineEventHandler,
	settings Settings,
	log zerolog.Logger,
) error {
	log = log.With().Str("component", "monitor").Logger()

	log.Trace().Msg("Checking handler")

	err := handler.Check(ctx) // TODO: introduce proper timeout

	switch {
	case err == nil:
		log.Info().Msg("Handler seems to be functional")
	case errors.Is(err, tsm.ErrUncheckable):
		log.Warn().Msg("Cannot verify that handler is functional")
	case errors.Is(err, context.Canceled):
		log.Trace().Msg("Handler check cancelled")
		return err
	case errors.Is(err, context.DeadlineExceeded):
		// XXX: should not occur atm since no timeout
		log.Trace().Msg("Handler check timed out")
		return err
	case settings.IgnoreStartupErrors:
		log.Error().Err(err).Msg("Handler is not functional")
	default:
		log.Warn().Err(err).Msg("Handler is not functional")
		return err
	}

	if len(settings.UserIDs) == 0 {
		log.Info().Msg("Not subscribing to anyone")
	}

	for _, userID := range settings.UserIDs {
		keep := false
		log := log.With().Str("user_id", userID).Logger()

		log.Trace().Msg("Subscribing")

		// TODO: introduce proper timeout and cancellation
		subID, err := svc.Subscribe(userID)

		if subID != "" {
			log = log.With().Str("sub_id", subID).Logger()
		}

		switch {
		case err == nil:
			keep = settings.KeepNewSubs
			log.Info().Msg("Subscribed")
		case errors.Is(err, tsm.ErrAlreadySubscribed):
			keep = settings.KeepExistingSubs
			log.Info().Msg("Already subscribed")
		case errors.Is(err, context.Canceled):
			// XXX: should not occur atm since no cancellation
			log.Trace().Msg("Subscription cancelled")
			return err
		case errors.Is(err, context.DeadlineExceeded):
			// XXX: should not occur atm since no timeout
			if settings.IgnoreSubscriptionErrors {
				log.Warn().Msg("Subscription timed out")
			} else {
				log.Error().Msg("Subscription timed out")
				return err
			}
		default:
			if settings.IgnoreSubscriptionErrors {
				log.Warn().Err(err).Msg("Failed to subscribe")
			} else {
				log.Error().Err(err).Msg("Failed to subscribe")
				return err
			}
		}

		if keep {
			log.Trace().Msg("Not unsubscribing at exit")
			continue
		} else {
			log.Trace().Msg("Will be unsubscribed at exit")
		}

		defer func(subID string, log zerolog.Logger) {
			log.Trace().Msg("Unsubscribing")

			// TODO: introduce proper timeout
			err := svc.Unsubscribe(subID)

			switch {
			case err == nil:
				log.Info().Msg("Unsubscribed")
			case errors.Is(err, tsm.ErrNotFound):
				log.Warn().Msg("Already unsubscribed")
			default:
				log.Error().Err(err).Msg("Failed to unsubscribe")
			}
		}(subID, log)
	}

	var wg sync.WaitGroup

	log.Trace().Msg("Starting event listener")

	events, err := svc.Listen(ctx)

	if err == nil {
		log.Info().Msg("Listening to events")
	} else {
		log.Error().Err(err).Msg("Failed to start listening")

		return err
	}

	for event := range events {
		wg.Add(1)

		go func(event tsm.TwitchStreamOnlineEvent) {
			defer wg.Done()

			log := log.With().Str("user_id", event.UserID).Str("user_name", event.UserName).Logger()

			log.Info().Msg("Started")

			err := handler.Handle(ctx, event)

			if err == nil {
				log.Info().Msg("Completed")
			} else if errors.Is(err, tsm.ErrCancelled) {
				log.Warn().Msg("Canceled")
			} else {
				log.Error().Err(err).Msg("Failed")
			}
		}(event)
	}

	log.Info().Msg("Waiting for running jobs to finish")
	wg.Wait()
	log.Info().Msg("Done")

	return nil
}
