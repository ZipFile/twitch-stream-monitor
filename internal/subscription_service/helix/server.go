package helix

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/nicklaw5/helix"
	"github.com/rs/zerolog"

	tsm "github.com/ZipFile/twitch-stream-monitor/internal"
)

type server struct {
	Svc *service
	Log zerolog.Logger
	Out chan tsm.TwitchStreamOnlineEvent
	Err error
}

type eventSubNotification struct {
	Subscription helix.EventSubSubscription `json:"subscription"`
	Challenge    string                     `json:"challenge"`
	Event        json.RawMessage            `json:"event"`
}

func (srv *server) start(ctx context.Context) error {
	if ctx == nil {
		srv.Log.Panic().Msg("Context is required")
	}

	httpServer := http.Server{
		Addr:    srv.Svc.Bind,
		Handler: srv,
	}
	serverStartErrContext, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)

	go func() {
		defer cancel()

		srv.Log.Trace().Str("addr", httpServer.Addr).Msg("Starting HTTP server")

		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			srv.Err = err

			srv.Log.Error().Err(err).Msg("Failed to start HTTP server")
		}
	}()

	<-serverStartErrContext.Done()

	go func() {
		<-ctx.Done()

		srv.Log.Debug().Msg("Shutting down")

		err := httpServer.Shutdown(context.Background())

		close(srv.Out)

		if err != nil {
			srv.Log.Error().Err(err).Msg("Failed to shutdown")
		}
	}()

	if srv.Err == nil {
		srv.Log.Info().Str("addr", httpServer.Addr).Msg("HTTP server started")
	}

	return srv.Err
}

func (srv *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		srv.Log.Error().Err(err).Msg("Failed to read request body")
		return
	}

	status, text := 500, "unknown error"

	defer r.Body.Close()
	defer func() {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(status)
		w.Write([]byte(text))
	}()

	if helix.VerifyEventSubNotification(srv.Svc.WebhookSecret, r.Header, string(body)) {
		srv.Log.Trace().Msg("Signature check OK")
	} else {
		srv.Log.Trace().Msg("Signature check failed")
		status, text = 403, "invalid signature"
		return
	}

	var payload eventSubNotification
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&payload)

	if err != nil {
		srv.Log.Error().Err(err).Msg("Failed to decode request body")
		status, text = 400, "invalid payload"
		return
	}

	if payload.Challenge != "" {
		srv.Log.Debug().Str("challenge", payload.Challenge).Msg("Received challenge")
		status, text = 200, payload.Challenge
		return
	}

	var onlineEvent helix.EventSubStreamOnlineEvent

	err = json.NewDecoder(bytes.NewReader(payload.Event)).Decode(&onlineEvent)

	if err != nil {
		srv.Log.Error().Err(err).Msg("Failed to decode event payload")
		status, text = 400, "invalid payload"
		return
	}

	srv.Log.Info().Str("user_login", onlineEvent.BroadcasterUserLogin).Str("event_type", onlineEvent.Type).Msg("New event")
	status, text = 200, "ok"

	if onlineEvent.Type == "live" {
		srv.Out <- tsm.TwitchStreamOnlineEvent{
			UserID:    onlineEvent.BroadcasterUserID,
			UserLogin: onlineEvent.BroadcasterUserLogin,
			UserName:  onlineEvent.BroadcasterUserName,
			StartedAt: onlineEvent.StartedAt.Time,
		}
	}
}
