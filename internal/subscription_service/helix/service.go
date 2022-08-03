package helix

import (
	"context"
	"errors"

	"github.com/nicklaw5/helix"
	"github.com/rs/zerolog"

	tsm "github.com/ZipFile/twitch-stream-monitor/internal"
)

// https://dev.twitch.tv/docs/eventsub
// https://github.com/nicklaw5/helix/blob/main/eventsub.go
// https://github.com/nicklaw5/helix/blob/main/docs/eventsub_docs.md

type service struct {
	Client        *helix.Client
	CallbackURL   string
	WebhookSecret string
	Bind          string
	Log           zerolog.Logger
}

func New(
	client *helix.Client,
	callbackURL, webhookSecret, bind string,
	log zerolog.Logger,
) (tsm.TwitchOnlineSubscriptionService, error) {
	if bind == "" {
		bind = "localhost:8000"
	}
	webhookSecretLength := len(webhookSecret)

	if webhookSecretLength < 10 || webhookSecretLength > 100 {
		return nil, errors.New("Webhook secret length must be between 10 and 100 chars.")
	}

	return &service{
		Client:        client,
		CallbackURL:   callbackURL,
		WebhookSecret: webhookSecret,
		Bind:          bind,
		Log:           log.With().Str("component", "helix_subscription_service").Logger(),
	}, nil
}

func (s *service) Subscribe(broadcaster_id string) (string, error) {
	response, err := s.Client.CreateEventSubSubscription(&helix.EventSubSubscription{
		Type:    helix.EventSubTypeStreamOnline,
		Version: "1",
		Condition: helix.EventSubCondition{
			BroadcasterUserID: broadcaster_id,
		},
		Transport: helix.EventSubTransport{
			Method:   "webhook",
			Callback: s.CallbackURL,
			Secret:   s.WebhookSecret,
		},
	})

	if err != nil {
		return "", err
	}

	if response.StatusCode == 409 {
		return "", tsm.ErrAlreadySubscribed
	}

	if response.StatusCode >= 400 {
		return "", errorFromResponse(&response.ResponseCommon)
	}

	for _, sub := range response.Data.EventSubSubscriptions {
		return sub.ID, nil
	}

	return "", errors.New("No subscriptions were registered")
}

func (s *service) Unsubscribe(subscription_id string) error {
	if subscription_id == "" {
		return errors.New("Missing subscription id")
	}

	response, err := s.Client.RemoveEventSubSubscription(subscription_id)

	if err != nil {
		return err
	}

	if response.StatusCode == 404 {
		return tsm.ErrNotFound
	}

	if response.StatusCode >= 400 {
		return errorFromResponse(&response.ResponseCommon)
	}

	return nil
}

func (svc *service) Listen(ctx context.Context) (<-chan tsm.TwitchStreamOnlineEvent, error) {
	out := make(chan tsm.TwitchStreamOnlineEvent)
	srv := server{
		Svc: svc,
		Out: out,
		Log: svc.Log.With().Str("component", "helix_subscription_server").Logger(),
	}
	err := srv.start(ctx)

	if err == nil {
		return out, nil
	}

	return nil, err
}

func (s *service) list(after string) ([]helix.EventSubSubscription, string, error) {
	response, err := s.Client.GetEventSubSubscriptions(&helix.EventSubSubscriptionsParams{
		Type:  helix.EventSubTypeStreamOnline,
		After: after,
	})

	if err != nil {
		return nil, "", err
	}

	return response.Data.EventSubSubscriptions, response.Data.Pagination.Cursor, nil
}

func (s *service) List() ([]tsm.TwitchStreamOnlineEventSubscription, error) {
	var after string
	var subs []helix.EventSubSubscription
	var err error
	var out []tsm.TwitchStreamOnlineEventSubscription

	for {
		s.Log.Trace().Str("after", after).Msg("Fetching subs")

		subs, after, err = s.list(after)

		if err != nil {
			return nil, err
		}

		s.Log.Trace().Int("subs_count", len(subs)).Msg("Fetched subs")

		if out == nil {
			out = make([]tsm.TwitchStreamOnlineEventSubscription, 0, len(subs))
		}

		for _, sub := range subs {
			out = append(out, tsm.TwitchStreamOnlineEventSubscription{
				ID:          sub.ID,
				Status:      sub.Status,
				UserID:      sub.Condition.BroadcasterUserID,
				CallbackURL: sub.Transport.Callback,
			})
		}

		if len(subs) == 0 || after == "" {
			break
		}
	}

	return out, nil
}
