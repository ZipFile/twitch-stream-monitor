package helix

import (
	"time"

	"github.com/nicklaw5/helix"
	"github.com/rs/zerolog"

	"github.com/ZipFile/twitch-stream-monitor/internal/app_access_token"
)

type manager struct {
	Client *helix.Client
	Log    zerolog.Logger
	Now    func() time.Time
}

func NewManager(client *helix.Client, log zerolog.Logger) app_access_token.Manager {
	return &manager{
		Client: client,
		Log:    log.With().Str("component", "helix_token_manager").Logger(),
		Now:    time.Now,
	}
}

func (m *manager) Validate(token string) (bool, error) {
	ok, _, err := m.Client.ValidateToken(token)

	return ok, err
}

func (m *manager) Request() (string, time.Time, error) {
	var expires time.Time
	response, err := m.Client.RequestAppAccessToken(nil)

	if err != nil {
		return "", expires, err
	}

	expires = m.Now().Add(time.Duration(response.Data.ExpiresIn) * time.Second)

	return response.Data.AccessToken, expires, nil
}
