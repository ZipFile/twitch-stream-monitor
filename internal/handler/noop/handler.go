package noop

import (
	"context"
	"strings"

	tsm "github.com/ZipFile/twitch-stream-monitor/internal"
)

// A handler that does nothing. For use when testing.
type Handler struct {
	CheckError  error
	HandleError error
}

func New() tsm.TwitchStreamOnlineEventHandler {
	return &Handler{}
}

// Always returns "noop".
func (*Handler) Name() string {
	return "noop"
}

// Return value of the CheckError field.
func (h *Handler) Check(ctx context.Context) error {
	return h.CheckError
}

// If event.UserLogin contains "error" substring - returns value of the HandleError field.
// If event.UserLogin contains "wait" substring - waits for context to be done, then returns ErrCancelled.
// In all other cases - returns nil.
func (h *Handler) Handle(ctx context.Context, event tsm.TwitchStreamOnlineEvent) error {
	if strings.Contains(event.UserLogin, "error") {
		return h.HandleError
	}

	if strings.Contains(event.UserLogin, "wait") {
		<-ctx.Done()

		return tsm.ErrCancelled
	}

	return nil
}
