package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type environLoader struct {
	LookupEnv func(string) string
}

func defaultLookupEnv(name string) string {
	value, _ := os.LookupEnv(name)
	return value
}

func NewEnvironLoader(lookupEnv func(string) string) Loader {
	if lookupEnv == nil {
		lookupEnv = defaultLookupEnv
	}

	return &environLoader{
		LookupEnv: lookupEnv,
	}
}

func (l *environLoader) lookupStringList(name, sep string) []string {
	if s := l.LookupEnv(name); s != "" {
		return strings.Split(s, sep)
	}

	return nil
}

func (l *environLoader) lookupUint64(name string) (uint64, error) {
	if s := l.LookupEnv(name); s != "" {
		return strconv.ParseUint(s, 10, 64)
	}

	return 0, nil
}

func (l *environLoader) lookupDuration(name string) (time.Duration, error) {
	if s := l.LookupEnv(name); s != "" {
		return time.ParseDuration(s)
	}

	return 0, nil
}

func (l *environLoader) lookupBool(name string) (bool, error) {
	if s := l.LookupEnv(name); s != "" {
		return strconv.ParseBool(s)
	}

	return false, nil
}

func (l *environLoader) Load() (*Config, error) {
	var err error

	c := Config{}

	c.BroadcasterUserIDs = l.lookupStringList("TWITCH_BROADCASTERS", ",")
	c.KeepExistingSubs, err = l.lookupBool("TWITCH_MONITOR_KEEP_EXISTING_SUBS")

	if err != nil {
		return nil, err
	}

	c.KeepNewSubs, err = l.lookupBool("TWITCH_MONITOR_KEEP_NEW_SUBS")

	if err != nil {
		return nil, err
	}

	c.IgnoreStartupErrors, err = l.lookupBool("TWITCH_MONITOR_IGNORE_START_ERRORS")

	if err != nil {
		return nil, err
	}

	c.IgnoreSubscriptionErrors, err = l.lookupBool("TWITCH_MONITOR_IGNORE_SUB_ERRORS")

	if err != nil {
		return nil, err
	}

	c.EventHandlerType = l.LookupEnv("TWITCH_MONITOR_HANDLER")
	c.CheckTimeout, err = l.lookupDuration("TWITCH_MONITOR_CHECK_TIMEOUT")

	if err != nil {
		return nil, err
	}

	c.Twitch.ClientID = l.LookupEnv("TWITCH_CLIENT_ID")
	c.Twitch.ClientSecret = l.LookupEnv("TWITCH_CLIENT_SECRET")
	c.Twitch.AppAccessToken = l.LookupEnv("TWITCH_APP_ACCESS_TOKEN")
	c.Twitch.AppAccessTokenLocation = l.LookupEnv("TWITCH_APP_ACCESS_TOKEN_LOCATION")

	c.Subscription.CallbackURL = l.LookupEnv("TWITCH_EVENTSUB_CALLBACK_URL")
	c.Subscription.WebhookSecret = l.LookupEnv("TWITCH_WEBHOOK_SECRET")
	c.Subscription.Host = l.LookupEnv("TWITCH_MONITOR_HOST")
	c.Subscription.Port, err = l.lookupUint64("TWITCH_MONITOR_PORT")

	if err != nil {
		return nil, err
	}

	c.Streamlink.Path = l.LookupEnv("TWITCH_MONITOR_STREAMLINK_PATH")
	c.Streamlink.FileDir = l.LookupEnv("TWITCH_MONITOR_STREAMLINK_FILE_DIR")
	c.Streamlink.LogDir = l.LookupEnv("TWITCH_MONITOR_STREAMLINK_LOG_DIR")
	c.Streamlink.ConfigPath = l.LookupEnv("TWITCH_MONITOR_STREAMLINK_CONFIG")
	c.Streamlink.KillTimeout, err = l.lookupDuration("TWITCH_MONITOR_STREAMLINK_KILL_TIMEOUT")

	if err != nil {
		return nil, err
	}

	c.HTTPNotificator.URL = l.LookupEnv("TWITCH_MONITOR_HTTP_NOTIFICATOR_URL")
	c.HTTPNotificator.UserName = l.LookupEnv("TWITCH_MONITOR_HTTP_NOTIFICATOR_USERNAME")
	c.HTTPNotificator.Password = l.LookupEnv("TWITCH_MONITOR_HTTP_NOTIFICATOR_PASSWORD")

	c.Ngrok.TunnelsURL = l.LookupEnv("TWITCH_MONITOR_NGROK_TUNNELS_URL")
	c.Ngrok.TunnelName = l.LookupEnv("TWITCH_MONITOR_NGROK_TUNNEL_NAME")

	c.Log.Level = l.LookupEnv("TWITCH_MONITOR_LOG_LEVEL")
	c.Log.Pretty, err = l.lookupBool("TWITCH_MONITOR_LOG_PRETTY")

	if err != nil {
		return nil, err
	}

	c.Log.Stdout, err = l.lookupBool("TWITCH_MONITOR_LOG_STDOUT")

	if err != nil {
		return nil, err
	}

	return &c, nil
}
