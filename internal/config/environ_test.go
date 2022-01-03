package config

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
)

type testEnviron map[string]string

func (environ testEnviron) LookupEnv(name string) string {
	value, _ := environ[name]
	return value
}

func TestNewEnvironLoaderDefault(t *testing.T) {
	cl := NewEnvironLoader(nil).(*environLoader)
	name := fmt.Sprintf("TEST_ENV_N%d", rand.Int())

	os.Setenv(name, "test")
	defer os.Unsetenv(name)

	value := cl.LookupEnv(name)

	if value != "test" {
		t.Errorf("cl.LookupEnv(%v): %v; expected: \"test\"", name, value)
	}
}

func TestNewEnvironLoaderCustom(t *testing.T) {
	testEnviron := make(testEnviron)
	testEnviron["test"] = "test"
	cl := NewEnvironLoader(testEnviron.LookupEnv).(*environLoader)
	value := cl.LookupEnv("test")

	if value != "test" {
		t.Errorf("cl.LookupEnv(\"test\"): %v; expected: \"test\"", value)
	}
}

func TestEnvironLoaderOK(t *testing.T) {
	testEnviron := testEnviron{
		"TWITCH_BROADCASTERS":                      "123,456",
		"TWITCH_MONITOR_KEEP_EXISTING_SUBS":        "true",
		"TWITCH_MONITOR_KEEP_NEW_SUBS":             "true",
		"TWITCH_MONITOR_IGNORE_START_ERRORS":       "true",
		"TWITCH_MONITOR_IGNORE_SUB_ERRORS":         "true",
		"TWITCH_MONITOR_HANDLER":                   "noop",
		"TWITCH_CLIENT_ID":                         "test_client_id",
		"TWITCH_CLIENT_SECRET":                     "test_client_secret",
		"TWITCH_EVENTSUB_CALLBACK_URL":             "https://example.com/tsm",
		"TWITCH_WEBHOOK_SECRET":                    "hackme",
		"TWITCH_MONITOR_HOST":                      "localhost",
		"TWITCH_MONITOR_PORT":                      "9000",
		"TWITCH_MONITOR_STREAMLINK_PATH":           "/tmp/streamlink",
		"TWITCH_MONITOR_STREAMLINK_FILE_DIR":       "/tmp/streams",
		"TWITCH_MONITOR_STREAMLINK_LOG_DIR":        "/tmp/logs",
		"TWITCH_MONITOR_STREAMLINK_CONFIG":         "/tmp/streamlink.config",
		"TWITCH_MONITOR_STREAMLINK_KILL_TIMEOUT":   "99s",
		"TWITCH_MONITOR_LOG_LEVEL":                 "debug",
		"TWITCH_MONITOR_HTTP_NOTIFICATOR_URL":      "http://localhost:8888/online",
		"TWITCH_MONITOR_HTTP_NOTIFICATOR_USERNAME": "root",
		"TWITCH_MONITOR_HTTP_NOTIFICATOR_PASSWORD": "hackme",
		"TWITCH_MONITOR_NGROK_TUNNELS_URL":         "http://localhost:4040/api/tunnels",
		"TWITCH_MONITOR_NGROK_TUNNEL_NAME":         "test",
	}

	c, err := NewEnvironLoader(testEnviron.LookupEnv).Load()

	if err != nil {
		t.Error(err)
		return
	}

	if strings.Join(c.BroadcasterUserIDs, ",") != "123,456" {
		t.Errorf("c.BroadcasterUserIDs: %v; expected: []string{\"123\",\"456\"}", c.BroadcasterUserIDs)
	}

	if c.Subscription.Port != 9000 {
		t.Errorf("c.Subscription.Port: %d; expected: 9000", c.Subscription.Port)
	}

	if c.Streamlink.KillTimeout != 99*time.Second {
		t.Errorf("c.Streamlink.KillTimeout: %d; expected: %d", c.Streamlink.KillTimeout, 10*time.Second)
	}
}

func TestEnvironLoaderPortError(t *testing.T) {
	testEnviron := make(testEnviron)
	testEnviron["TWITCH_MONITOR_PORT"] = "test"
	c, err := NewEnvironLoader(testEnviron.LookupEnv).Load()

	if err == nil {
		t.Errorf("err: nil; expected: error")
		return
	}

	if c != nil {
		t.Errorf("c: %v; expected: nil", c)
	}
}

func TestEnvironLoaderKillTimeoutError(t *testing.T) {
	testEnviron := make(testEnviron)
	testEnviron["TWITCH_MONITOR_STREAMLINK_KILL_TIMEOUT"] = "test"
	c, err := NewEnvironLoader(testEnviron.LookupEnv).Load()

	if err == nil {
		t.Errorf("err: nil; expected: error")
		return
	}

	if c != nil {
		t.Errorf("c: %v; expected: nil", c)
	}
}
