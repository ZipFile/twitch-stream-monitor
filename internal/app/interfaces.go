package app

import (
	"github.com/rs/zerolog"

	tsm "github.com/ZipFile/twitch-stream-monitor/internal"
)

type App interface {
	Init() error
	Run() error
	GetTwitchOnlineSubscriptionService() tsm.TwitchOnlineSubscriptionService
	GetLogger() *zerolog.Logger
}
