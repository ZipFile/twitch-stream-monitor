package noop

import (
	"context"
	"errors"
	"testing"

	tsm "github.com/ZipFile/twitch-stream-monitor/internal"
)

func TestNew(t *testing.T) {
	h := New()

	if _, ok := h.(*Handler); !ok {
		t.Errorf("Non-noop instance")
	}
}

func TestHandlerName(t *testing.T) {
	var h *Handler
	name := h.Name()

	if name != "noop" {
		t.Errorf("name: %v; expected: \"noop\"", name)
	}
}

func TestHandlerCheck(t *testing.T) {
	var h Handler

	for _, checkError := range []error{nil, errors.New("test")} {
		h.CheckError = checkError
		err := h.Check(context.Background())

		if err != checkError {
			t.Errorf("err: %v; expected: %v", err, checkError)
		}
	}
}

func TestHandlerHandleOK(t *testing.T) {
	var h Handler
	var event tsm.TwitchStreamOnlineEvent

	err := h.Handle(context.TODO(), event)

	if err != nil {
		t.Errorf("err: %v; expected: nil", err)
	}
}

func TestHandlerHandleError(t *testing.T) {
	var h Handler
	event := tsm.TwitchStreamOnlineEvent{
		UserLogin: "test_error_123",
	}

	for _, handleError := range []error{nil, errors.New("test")} {
		h.HandleError = handleError
		err := h.Handle(context.TODO(), event)

		if err != handleError {
			t.Errorf("err: %v; expected: %v", err, handleError)
		}
	}
}

func TestHandlerHandleWait(t *testing.T) {
	var h Handler
	event := tsm.TwitchStreamOnlineEvent{
		UserLogin: "test_wait_123",
	}
	ctx, cancel := context.WithCancel(context.Background())

	cancel()

	err := h.Handle(ctx, event)

	if err != tsm.ErrCancelled {
		t.Errorf("err: %v; expected: ErrCancelled", err)
	}
}
