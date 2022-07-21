package testing

import (
	"github.com/rs/zerolog"

	tsm "github.com/ZipFile/twitch-stream-monitor/internal"
)

type App struct {
	InitError error
	RunError  error
	TOSS      tsm.TwitchOnlineSubscriptionService
	Log       *zerolog.Logger
}

func (a *App) Init() error {
	return a.InitError
}

func (a *App) Run() error {
	return a.RunError
}

func (a *App) GetTwitchOnlineSubscriptionService() tsm.TwitchOnlineSubscriptionService {
	return a.TOSS
}

func (a *App) GetLogger() *zerolog.Logger {
	return a.Log
}
