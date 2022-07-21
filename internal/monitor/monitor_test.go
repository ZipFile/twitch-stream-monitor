package tsm

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"

	tsm "github.com/ZipFile/twitch-stream-monitor/internal"
	noop "github.com/ZipFile/twitch-stream-monitor/internal/handler/noop"
	tsm_testing "github.com/ZipFile/twitch-stream-monitor/internal/testing"
)

func TestMonitorCheckFail(t *testing.T) {
	errors := []error{
		errors.New("check failed"),
		context.Canceled,
		context.DeadlineExceeded,
	}

	for _, expected := range errors {
		err := Monitor(
			context.TODO(),
			nil,
			&noop.Handler{CheckError: expected},
			Settings{},
			zerolog.Nop(),
		)

		if err != expected {
			t.Errorf("err: %v; expected: %v", err, expected)
		}
	}
}

func TestMonitorSubNotFound(t *testing.T) {
	handler := &noop.Handler{CheckError: tsm.ErrUncheckable}
	svc := tsm_testing.NewFakeTwitchOnlineSubscriptionService("123")
	err := Monitor(
		context.TODO(),
		svc,
		handler,
		Settings{UserIDs: []string{"123", "456"}},
		zerolog.Nop(),
	)

	if err != tsm.ErrNotFound {
		t.Errorf("err: %v; expected: ErrNotFound", err)
	}
}

func TestMonitorSubFail(t *testing.T) {
	handler := &noop.Handler{CheckError: tsm.ErrUncheckable}
	svc := tsm_testing.NewFakeTwitchOnlineSubscriptionService("123 subError")
	err := Monitor(
		context.TODO(),
		svc,
		handler,
		Settings{UserIDs: []string{"123 subError"}},
		zerolog.Nop(),
	)

	if err != tsm_testing.Error {
		t.Errorf("err: %v; expected: tsm_testing.Error", err)
	}
}

func TestMonitorKeepExistingSubs(t *testing.T) {
	handler := &noop.Handler{CheckError: tsm.ErrUncheckable}
	svc := tsm_testing.NewFakeTwitchOnlineSubscriptionService("123", "456")
	ctx, cancel := context.WithCancel(context.Background())

	svc.Subscribe("123")

	go func() {
		defer cancel()
		<-svc.Started
	}()

	err := Monitor(
		ctx,
		svc,
		handler,
		Settings{UserIDs: []string{"123", "456"}, KeepExistingSubs: true},
		zerolog.Nop(),
	)

	if err != nil {
		t.Errorf("err: %v; expected: nil", err)
	}

	if _, ok := svc.Sbus["123"]; !ok {
		t.Errorf("should keep existing subs")
	}

	if _, ok := svc.Sbus["456"]; ok {
		t.Errorf("should delete new subs")
	}
}

func TestMonitorKeepNewSubs(t *testing.T) {
	handler := &noop.Handler{CheckError: tsm.ErrUncheckable}
	svc := tsm_testing.NewFakeTwitchOnlineSubscriptionService("123", "456")
	ctx, cancel := context.WithCancel(context.Background())

	svc.Subscribe("123")

	go func() {
		defer cancel()
		<-svc.Started
	}()

	err := Monitor(
		ctx,
		svc,
		handler,
		Settings{UserIDs: []string{"123", "456"}, KeepNewSubs: true},
		zerolog.Nop(),
	)

	if err != nil {
		t.Errorf("err: %v; expected: nil", err)
	}

	if _, ok := svc.Sbus["123"]; ok {
		t.Errorf("should delete existing subs")
	}

	if _, ok := svc.Sbus["456"]; !ok {
		t.Errorf("should keep new subs")
	}
}

func TestMonitorUnsubFail(t *testing.T) {
	handler := &noop.Handler{}
	svc := tsm_testing.NewFakeTwitchOnlineSubscriptionService("123 unsubError")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := Monitor(
		ctx,
		svc,
		handler,
		Settings{UserIDs: []string{"123 unsubError"}},
		zerolog.Nop(),
	)

	if err != nil {
		t.Errorf("err: %v; expected: nil", err)
	}

	// TODO: inspect logs to do proper assertion
}

func TestMonitorListenFail(t *testing.T) {
	handler := &noop.Handler{}
	svc := tsm_testing.NewFakeTwitchOnlineSubscriptionService()
	svc.ListenError = tsm_testing.Error
	err := Monitor(
		context.TODO(),
		svc,
		handler,
		Settings{},
		zerolog.Nop(),
	)

	if err != tsm_testing.Error {
		t.Errorf("err: %v; expected: tsm_testing.Error", err)
	}
}

func TestMonitorEventsOK(t *testing.T) {
	handler := &noop.Handler{HandleError: tsm_testing.Error}
	svc := tsm_testing.NewFakeTwitchOnlineSubscriptionService()
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		svc.Events <- tsm.TwitchStreamOnlineEvent{
			UserID:    "1234",
			UserLogin: "test",
			UserName:  "Test",
			StartedAt: time.Date(2021, time.September, 17, 23, 45, 0, 0, time.UTC),
		}
		svc.Events <- tsm.TwitchStreamOnlineEvent{
			UserID:    "5678",
			UserLogin: "wait",
			UserName:  "Wait",
			StartedAt: time.Date(2021, time.September, 17, 23, 50, 0, 0, time.UTC),
		}
		svc.Events <- tsm.TwitchStreamOnlineEvent{
			UserID:    "9000",
			UserLogin: "error",
			UserName:  "Error",
			StartedAt: time.Date(2021, time.September, 17, 23, 55, 0, 0, time.UTC),
		}
		cancel()
	}()

	err := Monitor(ctx, svc, handler, Settings{}, zerolog.Nop())

	if err != nil {
		t.Errorf("err: %v; expected: nil", err)
	}

	// TODO: inspect logs to do proper assertion
}

func TestMonitorIgnoreErrors(t *testing.T) {
	handler := &noop.Handler{CheckError: tsm_testing.Error}
	svc := tsm_testing.NewFakeTwitchOnlineSubscriptionService("123 subError")
	ctx, cancel := context.WithCancel(context.Background())
	settings := Settings{
		UserIDs:                  []string{"123 subError"},
		IgnoreStartupErrors:      true,
		IgnoreSubscriptionErrors: true,
	}

	cancel()

	err := Monitor(ctx, svc, handler, settings, zerolog.Nop())

	if err != nil {
		t.Errorf("err: %v; expected: nil", err)
	}
}
