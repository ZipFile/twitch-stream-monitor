package testing

import (
	"context"
	"errors"
	"fmt"
	"strings"

	tsm "github.com/ZipFile/twitch-stream-monitor/internal"
)

type FakeTwitchOnlineSubscriptionService struct {
	I           int
	IDs         map[string]bool
	Subs        map[string]string
	Sbus        map[string]string
	Events      chan tsm.TwitchStreamOnlineEvent
	Started     chan interface{}
	ListenError error
}

var Error = errors.New("test error")
var FakeRecordingError = errors.New("recording failed")

func NewFakeTwitchOnlineSubscriptionService(ids ...string) *FakeTwitchOnlineSubscriptionService {
	knownIds := make(map[string]bool)

	for _, id := range ids {
		knownIds[id] = true
	}

	return &FakeTwitchOnlineSubscriptionService{
		IDs:    knownIds,
		Subs:   make(map[string]string),
		Sbus:   make(map[string]string),
		Events: make(chan tsm.TwitchStreamOnlineEvent),
	}
}

func (svc *FakeTwitchOnlineSubscriptionService) Subscribe(id string) (string, error) {
	if strings.Contains(id, " subError") {
		return "", Error
	}

	_, ok := svc.IDs[id]

	if !ok {
		return "", tsm.ErrNotFound
	}

	subId, ok := svc.Sbus[id]

	if ok {
		return subId, tsm.ErrAlreadySubscribed
	}

	svc.I++

	subId = fmt.Sprintf("sub%d", svc.I)

	svc.Subs[subId] = id
	svc.Sbus[id] = subId

	return subId, nil
}

func (svc *FakeTwitchOnlineSubscriptionService) Unsubscribe(subId string) error {
	id, ok := svc.Subs[subId]

	if !ok {
		return tsm.ErrNotFound
	}

	if strings.Contains(id, " unsubError") {
		return Error
	}

	delete(svc.Subs, subId)
	delete(svc.Sbus, id)

	return nil
}

func (svc *FakeTwitchOnlineSubscriptionService) Listen(ctx context.Context) (<-chan tsm.TwitchStreamOnlineEvent, error) {
	if svc.ListenError != nil {
		return nil, svc.ListenError
	}

	svc.Started = make(chan interface{})
	out := make(chan tsm.TwitchStreamOnlineEvent)

	go func() {
		close(svc.Started)
		for {
			select {
			case <-ctx.Done():
				close(out)
				return
			case event := <-svc.Events:
				out <- event
			}
		}
	}()

	return out, nil
}
