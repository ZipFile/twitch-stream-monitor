package config

import (
	"context"
	"testing"
)

func TestConfigGetCallbackURLNil(t *testing.T) {
	var c *Config

	url, err := c.GetCallbackURL(context.TODO())

	if err != nil {
		t.Errorf("err: %v; expected: nil", err)
	}

	if url != "" {
		t.Errorf("url: %v; expected: \"\"", url)
	}
}

func TestConfigGetCallbackURLOK(t *testing.T) {
	c := Config{Subscription: SubscriptionConfig{CallbackURL: "test"}}
	url, err := c.GetCallbackURL(context.TODO())

	if err != nil {
		t.Errorf("err: %v; expected: nil", err)
	}

	if url != "test" {
		t.Errorf("url: %v; expected: \"test\"", url)
	}
}

func TestConfigSetDefaultsNil(t *testing.T) {
	var c *Config

	c.SetDefaults()
}

func TestConfigSetDefaultsOK(t *testing.T) {
	var c Config

	c.SetDefaults()

	if c.EventHandlerType == "" {
		t.Errorf("c.EventHandlerType: %v; expected: \"streamlink\"", c.EventHandlerType)
	}

	if c.Subscription.Port == 0 {
		t.Errorf("c.Subscription.Port: %v; expected: 7578", c.Subscription.Port)
	}

	if c.Streamlink.Path == "" {
		t.Errorf("c.Streamlink.Path: %v; expected: \"streamlink\"", c.Streamlink.Path)
	}

	if c.Streamlink.FileDir == "" {
		t.Errorf("c.Streamlink.FileDir: %v; expected: \".\"", c.Streamlink.FileDir)
	}

	if c.Streamlink.LogDir == "" {
		t.Errorf("c.Streamlink.LogDir: %v; expected: \".\"", c.Streamlink.LogDir)
	}

	if c.Streamlink.KillTimeout == 0 {
		t.Errorf("c.Streamlink.KillTimeout: %v; expected: \".\"", c.Streamlink.KillTimeout)
	}

	if c.Log.Level == "" {
		t.Errorf("c.Log.Level: %v; expected: \"info\"", c.Log.Level)
	}
}
