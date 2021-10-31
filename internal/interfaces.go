package tsm

import (
	"context"
	"errors"
	"time"
)

type TwitchStreamOnlineEvent struct {
	UserID    string    `json:"user_id"`
	UserLogin string    `json:"user_login"`
	UserName  string    `json:"user_name"`
	StartedAt time.Time `json:"started_at"`
}

var ErrCancelled = errors.New("Recording cancelled")
var ErrUncheckable = errors.New("Unable to check if handler works")
var ErrAlreadySubscribed = errors.New("Already subscribed")
var ErrNotSubscribed = errors.New("Not subscribed")
var ErrNotFound = errors.New("Not found")

type TwitchStreamOnlineEventHandler interface {
	Name() string
	Check(context.Context) error
	Handle(context.Context, TwitchStreamOnlineEvent) error
}

type TwitchOnlineSubscriptionService interface {
	Subscribe(string) (string, error)
	Unsubscribe(string) error
	Listen(context.Context) (<-chan TwitchStreamOnlineEvent, error)
}

type CallbackURLGetter interface {
	GetCallbackURL(context.Context) (string, error)
}
