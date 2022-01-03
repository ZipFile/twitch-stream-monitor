package config

import (
	"context"
	"time"
)

type TwitchConfig struct {
	ClientID               string
	ClientSecret           string
	AppAccessToken         string
	AppAccessTokenLocation string
}

type StreamlinkConfig struct {
	Path        string
	FileDir     string
	LogDir      string
	ConfigPath  string
	KillTimeout time.Duration
}

type HTTPNotificatorConfig struct {
	URL      string
	UserName string
	Password string
}

type LoggingConfig struct {
	Level  string
	Pretty bool
	Stdout bool
}

type SubscriptionConfig struct {
	CallbackURL   string
	WebhookSecret string
	Host          string
	Port          uint64
}

type NgrokConfig struct {
	TunnelsURL string
	TunnelName string
}

type Config struct {
	BroadcasterUserIDs       []string
	KeepExistingSubs         bool
	KeepNewSubs              bool
	IgnoreStartupErrors      bool
	IgnoreSubscriptionErrors bool

	EventHandlerType string
	CheckTimeout     time.Duration
	Twitch           TwitchConfig
	Subscription     SubscriptionConfig
	Streamlink       StreamlinkConfig
	HTTPNotificator  HTTPNotificatorConfig
	Ngrok            NgrokConfig
	Log              LoggingConfig
}

func (c *Config) GetCallbackURL(context.Context) (string, error) {
	if c == nil {
		return "", nil
	}

	return c.Subscription.CallbackURL, nil
}

func mergeStr(a, b *string) {
	if *a == "" {
		*a = *b
	}
}

func mergeStrList(a, b *[]string) {
	if len(*a) == 0 {
		*a = *b
	}
}

func mergeBool(a, b *bool) {
	if *b {
		*a = true
	}
}

func mergeUint64(a, b *uint64) {
	if *a == 0 {
		*a = *b
	}
}

func mergeDuration(a, b *time.Duration) {
	if *a == 0 {
		*a = *b
	}
}

func (a *Config) Merge(b *Config) {
	if a == nil || b == nil {
		return
	}

	mergeStrList(&a.BroadcasterUserIDs, &b.BroadcasterUserIDs)
	mergeBool(&a.KeepExistingSubs, &b.KeepExistingSubs)
	mergeBool(&a.KeepNewSubs, &b.KeepNewSubs)
	mergeBool(&a.IgnoreStartupErrors, &b.IgnoreStartupErrors)
	mergeBool(&a.IgnoreSubscriptionErrors, &b.IgnoreSubscriptionErrors)
	mergeStr(&a.EventHandlerType, &b.EventHandlerType)
	mergeDuration(&a.CheckTimeout, &b.CheckTimeout)
	mergeStr(&a.Twitch.ClientID, &b.Twitch.ClientID)
	mergeStr(&a.Twitch.ClientSecret, &b.Twitch.ClientSecret)
	mergeStr(&a.Subscription.CallbackURL, &b.Subscription.CallbackURL)
	mergeStr(&a.Subscription.WebhookSecret, &b.Subscription.WebhookSecret)
	mergeStr(&a.Subscription.Host, &b.Subscription.Host)
	mergeUint64(&a.Subscription.Port, &b.Subscription.Port)
	mergeStr(&a.Streamlink.Path, &b.Streamlink.Path)
	mergeStr(&a.Streamlink.FileDir, &b.Streamlink.FileDir)
	mergeStr(&a.Streamlink.LogDir, &b.Streamlink.LogDir)
	mergeDuration(&a.Streamlink.KillTimeout, &b.Streamlink.KillTimeout)
	mergeStr(&a.HTTPNotificator.URL, &b.HTTPNotificator.URL)
	mergeStr(&a.HTTPNotificator.UserName, &b.HTTPNotificator.UserName)
	mergeStr(&a.HTTPNotificator.Password, &b.HTTPNotificator.Password)
	mergeStr(&a.Ngrok.TunnelsURL, &b.Ngrok.TunnelsURL)
	mergeStr(&a.Ngrok.TunnelName, &b.Ngrok.TunnelName)
	mergeStr(&a.Log.Level, &b.Log.Level)
	mergeBool(&a.Log.Pretty, &b.Log.Pretty)
	mergeBool(&a.Log.Stdout, &b.Log.Stdout)
}

func (c *Config) SetDefaults() {
	c.Merge(&DefaultConfig)
}

var DefaultConfig = Config{
	EventHandlerType: "streamlink",
	CheckTimeout:     5 * time.Second,
	Subscription: SubscriptionConfig{
		Port: 29177,
	},
	Streamlink: StreamlinkConfig{
		Path:        "streamlink",
		FileDir:     ".",
		LogDir:      ".",
		KillTimeout: 60 * time.Second,
	},
	Log: LoggingConfig{
		Level: "info",
	},
}
