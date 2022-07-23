package helix_username_resolver

import (
	"github.com/nicklaw5/helix"
	"github.com/rs/zerolog"

	"github.com/ZipFile/twitch-stream-monitor/internal/username_resolver"
)

type usernameResolver struct {
	Client *helix.Client
	Log    zerolog.Logger
}

func New(client *helix.Client, log zerolog.Logger) username_resolver.UsernameResolver {
	return &usernameResolver{
		Client: client,
		Log:    log.With().Str("component", "helix_username_resolver").Logger(),
	}
}

func (ur *usernameResolver) Resolve(username string) (string, error) {
	response, err := ur.Client.GetUsers(&helix.UsersParams{Logins: []string{username}})

	if err != nil {
		return "", err
	}

	for _, user := range response.Data.Users {
		return user.ID, nil
	}

	return "", err
}
