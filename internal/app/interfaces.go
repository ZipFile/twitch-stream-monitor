package app

import (
	"github.com/rs/zerolog"

	tsm "github.com/ZipFile/twitch-stream-monitor/internal"
	tsm_aat "github.com/ZipFile/twitch-stream-monitor/internal/app_access_token"
	tsm_ur "github.com/ZipFile/twitch-stream-monitor/internal/username_resolver"
)

type App interface {
	Init() error
	Run() error
	GetUsernameResolver() tsm_ur.UsernameResolver
	GetTwitchOnlineSubscriptionService() tsm.TwitchOnlineSubscriptionService
	GetTokenStore() tsm_aat.Store
	GetTokenManager() tsm_aat.Manager
	GetLogger() *zerolog.Logger
}
