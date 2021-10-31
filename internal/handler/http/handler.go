package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	tsm "github.com/ZipFile/twitch-stream-monitor/internal"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Handler struct {
	URL      string
	UserName string
	Password string
	Client   HTTPClient
}

const UserAgent = "TwitchStreamMonitor/1.0"

var ErrInvalidAuthCredentials = errors.New("Invalid auth credentials")

func New(url, userName, password string) tsm.TwitchStreamOnlineEventHandler {
	return &Handler{
		URL:      url,
		UserName: userName,
		Password: password,
		Client:   &http.Client{},
	}
}

func (h *Handler) Name() string {
	return "http"
}

func (h *Handler) Check(ctx context.Context) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodHead, h.URL, nil)

	if err != nil {
		return err
	}

	h.setRequiredHeaders(request)

	response, err := h.Client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK || response.StatusCode == http.StatusMethodNotAllowed {
		return nil
	}

	if response.StatusCode == http.StatusUnauthorized {
		return ErrInvalidAuthCredentials
	}

	return tsm.ErrUncheckable
}

func (h *Handler) setRequiredHeaders(request *http.Request) {
	if h.UserName != "" && h.Password != "" {
		request.SetBasicAuth(h.UserName, h.Password)
	}

	request.Header.Set("user-agent", UserAgent)
}

func (h *Handler) Handle(ctx context.Context, event tsm.TwitchStreamOnlineEvent) error {
	body := new(bytes.Buffer)

	json.NewEncoder(body).Encode(event)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, h.URL, body)

	if err != nil {
		return err
	}

	h.setRequiredHeaders(request)
	request.Header.Set("content-type", "application/json")

	response, err := h.Client.Do(request)

	if err != nil {
		return err
	}

	response.Body.Close()

	return nil
}
