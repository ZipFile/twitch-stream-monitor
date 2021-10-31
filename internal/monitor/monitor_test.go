package tsm

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"

	tsm "github.com/ZipFile/twitch-stream-monitor/internal"
	noop "github.com/ZipFile/twitch-stream-monitor/internal/handler/noop"
)

type fakeTwitchOnlineSubscriptionService struct {
	i           int
	ids         map[string]bool
	subs        map[string]string
	sbus        map[string]string
	events      chan tsm.TwitchStreamOnlineEvent
	started     chan interface{}
	listenError error
}

var testError = errors.New("test error")
var fakeRecordingError = errors.New("recording failed")

func newFakeTwitchOnlineSubscriptionService(ids ...string) *fakeTwitchOnlineSubscriptionService {
	knownIds := make(map[string]bool)

	for _, id := range ids {
		knownIds[id] = true
	}

	return &fakeTwitchOnlineSubscriptionService{
		ids:    knownIds,
		subs:   make(map[string]string),
		sbus:   make(map[string]string),
		events: make(chan tsm.TwitchStreamOnlineEvent),
	}
}

func (svc *fakeTwitchOnlineSubscriptionService) Subscribe(id string) (string, error) {
	if strings.Contains(id, " subError") {
		return "", testError
	}

	_, ok := svc.ids[id]

	if !ok {
		return "", tsm.ErrNotFound
	}

	subId, ok := svc.sbus[id]

	if ok {
		return subId, tsm.ErrAlreadySubscribed
	}

	svc.i++

	subId = fmt.Sprintf("sub%d", svc.i)

	svc.subs[subId] = id
	svc.sbus[id] = subId

	return subId, nil
}

func (svc *fakeTwitchOnlineSubscriptionService) Unsubscribe(subId string) error {
	id, ok := svc.subs[subId]

	if !ok {
		return tsm.ErrNotSubscribed
	}

	if strings.Contains(id, " unsubError") {
		return testError
	}

	delete(svc.subs, subId)
	delete(svc.sbus, id)

	return nil
}

func (svc *fakeTwitchOnlineSubscriptionService) Listen(ctx context.Context) (<-chan tsm.TwitchStreamOnlineEvent, error) {
	if svc.listenError != nil {
		return nil, svc.listenError
	}

	svc.started = make(chan interface{})
	out := make(chan tsm.TwitchStreamOnlineEvent)

	go func() {
		close(svc.started)
		for {
			select {
			case <-ctx.Done():
				close(out)
				return
			case event := <-svc.events:
				out <- event
			}
		}
	}()

	return out, nil
}

func TestMonitorCheckFail(t *testing.T) {
	checkError := errors.New("check failed")
	handler := &noop.Handler{CheckError: checkError}
	err := Monitor(context.TODO(), nil, handler, nil, zerolog.Nop())

	if err != checkError {
		t.Errorf("err: %v; expected: checkError", err)
	}
}

func TestMonitorSubNotFound(t *testing.T) {
	handler := &noop.Handler{CheckError: tsm.ErrUncheckable}
	svc := newFakeTwitchOnlineSubscriptionService("123")
	err := Monitor(
		context.TODO(),
		svc,
		handler,
		[]string{"123", "456"},
		zerolog.Nop(),
	)

	if err != tsm.ErrNotFound {
		t.Errorf("err: %v; expected: ErrNotFound", err)
	}
}

func TestMonitorSubFail(t *testing.T) {
	handler := &noop.Handler{CheckError: tsm.ErrUncheckable}
	svc := newFakeTwitchOnlineSubscriptionService("123 subError")
	err := Monitor(
		context.TODO(),
		svc,
		handler,
		[]string{"123 subError"},
		zerolog.Nop(),
	)

	if err != testError {
		t.Errorf("err: %v; expected: testError", err)
	}
}

func TestMonitorDirty(t *testing.T) {
	handler := &noop.Handler{CheckError: tsm.ErrUncheckable}
	svc := newFakeTwitchOnlineSubscriptionService("123", "456")
	ctx, cancel := context.WithCancel(context.Background())

	svc.Subscribe("123")

	go func() {
		defer cancel()
		<-svc.started
		svc.Unsubscribe(svc.sbus["456"])
	}()

	err := Monitor(
		ctx,
		svc,
		handler,
		[]string{"123", "456"},
		zerolog.Nop(),
	)

	if err != nil {
		t.Errorf("err: %v; expected: nil", err)
	}

	if _, ok := svc.sbus["123"]; !ok {
		t.Errorf("should not unsubscribe already subscribed channels")
	}
}

func TestMonitorUnsubFail(t *testing.T) {
	handler := &noop.Handler{}
	svc := newFakeTwitchOnlineSubscriptionService("123 unsubError")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := Monitor(
		ctx,
		svc,
		handler,
		[]string{"123 unsubError"},
		zerolog.Nop(),
	)

	if err != nil {
		t.Errorf("err: %v; expected: nil", err)
	}

	// TODO: inspect logs to do proper assertion
}

func TestMonitorListenFail(t *testing.T) {
	handler := &noop.Handler{}
	svc := newFakeTwitchOnlineSubscriptionService()
	svc.listenError = testError
	err := Monitor(
		context.TODO(),
		svc,
		handler,
		nil,
		zerolog.Nop(),
	)

	if err != testError {
		t.Errorf("err: %v; expected: testError", err)
	}
}

func TestMonitorEventsOK(t *testing.T) {
	handler := &noop.Handler{HandleError: testError}
	svc := newFakeTwitchOnlineSubscriptionService()
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		svc.events <- tsm.TwitchStreamOnlineEvent{
			UserID:    "1234",
			UserLogin: "test",
			UserName:  "Test",
			StartedAt: time.Date(2021, time.September, 17, 23, 45, 0, 0, time.UTC),
		}
		svc.events <- tsm.TwitchStreamOnlineEvent{
			UserID:    "5678",
			UserLogin: "wait",
			UserName:  "Wait",
			StartedAt: time.Date(2021, time.September, 17, 23, 50, 0, 0, time.UTC),
		}
		svc.events <- tsm.TwitchStreamOnlineEvent{
			UserID:    "9000",
			UserLogin: "error",
			UserName:  "Error",
			StartedAt: time.Date(2021, time.September, 17, 23, 55, 0, 0, time.UTC),
		}
		cancel()
	}()

	err := Monitor(ctx, svc, handler, nil, zerolog.Nop())

	if err != nil {
		t.Errorf("err: %v; expected: nil", err)
	}

	// TODO: inspect logs to do proper assertion
}
