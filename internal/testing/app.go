package testing

import (
	"github.com/rs/zerolog"

	tsm "github.com/ZipFile/twitch-stream-monitor/internal"
	tsm_aat "github.com/ZipFile/twitch-stream-monitor/internal/app_access_token"
	tsm_ur "github.com/ZipFile/twitch-stream-monitor/internal/username_resolver"
)

type App struct {
	InitError error
	RunError  error
	UR        tsm_ur.UsernameResolver
	AATS      tsm_aat.Store
	AATM      tsm_aat.Manager
	TOSS      tsm.TwitchOnlineSubscriptionService
	Log       *zerolog.Logger
}

func (a *App) Init() error {
	return a.InitError
}

func (a *App) Run() error {
	return a.RunError
}

func (a *App) GetUsernameResolver() tsm_ur.UsernameResolver {
	return a.UR
}

func (a *App) GetTokenManager() tsm_aat.Manager {
	return a.AATM
}

func (a *App) GetTokenStore() tsm_aat.Store {
	return a.AATS
}

func (a *App) GetTwitchOnlineSubscriptionService() tsm.TwitchOnlineSubscriptionService {
	return a.TOSS
}

func (a *App) GetLogger() *zerolog.Logger {
	return a.Log
}
