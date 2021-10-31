package tsm

import (
	"context"
	"errors"
	"sync"

	"github.com/rs/zerolog"

	tsm "github.com/ZipFile/twitch-stream-monitor/internal"
)

// Main logic for monitor sub-command.
//
// Performs check if handler is functional, subscribes to the each user and then
// starts listening to online events. Unsubscribes on shutdown.
//
// Args:
// * ctx: Context used for graceful shutdown. Passed directly to svc.Listen and into the handler.Handle.
// * svc: Instantiate of the online subscription instance.
// * handler: Online event handler.
// * userIDs: List of Twitch user (broadcaster) ids. Not usernames.
// * log: Logger instance.
func Monitor(
	ctx context.Context,
	svc tsm.TwitchOnlineSubscriptionService,
	handler tsm.TwitchStreamOnlineEventHandler,
	userIDs []string,
	log zerolog.Logger,
) error {
	log = log.With().Str("component", "monitor").Logger()

	err := handler.Check(ctx) // TODO: introduce proper timeout

	if err != nil {
		if errors.Is(err, tsm.ErrUncheckable) {
			log.Warn().Msg("Cannot verify that handler is functional")
		} else {
			log.Error().Err(err).Msg("Handler is not functional")

			return err
		}
	}

	if len(userIDs) == 0 {
		log.Info().Msg("Not subscribing to anyone")
	}

	for _, userID := range userIDs {
		subId, err := svc.Subscribe(userID)
		log := log.With().Str("user_id", userID).Logger()

		if err == nil {
			log := log.With().Str("sub_id", subId).Logger()

			log.Info().Msg("Subscribed")

			defer func(userID, subId string, log *zerolog.Logger) {
				err := svc.Unsubscribe(subId)

				if err == nil {
					log.Info().Msg("Unsubscribed")
				} else if errors.Is(err, tsm.ErrNotSubscribed) {
					log.Warn().Msg("Already unsubscribed")
				} else {
					log.Error().Err(err).Msg("Failed to unsubscribe")
				}
			}(userID, subId, &log)
		} else if errors.Is(err, tsm.ErrAlreadySubscribed) {
			log.Warn().Msg("Already subscribed")
		} else {
			log.Error().Err(err).Msg("Failed to subscribe")

			return err
		}
	}

	var wg sync.WaitGroup

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
