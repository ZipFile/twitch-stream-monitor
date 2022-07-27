package cli

import (
	"context"
	"flag"
	"testing"

	"github.com/google/subcommands"
	"github.com/rs/zerolog"

	tsm_testing "github.com/ZipFile/twitch-stream-monitor/internal/testing"
)

func TestSubscribeName(t *testing.T) {
	var s subscribe
	name := s.Name()
	expected := "subscribe"

	if name != expected {
		t.Errorf("name: %s; expected: %s", name, expected)
	}
}

func TestSubscribeSynopsis(t *testing.T) {
	var s subscribe
	synopsis := s.Synopsis()
	expected := "Subscribe to stream.online events for given twitch users."

	if synopsis != expected {
		t.Errorf("synopsis: %s; expected: %s", synopsis, expected)
	}
}

func TestSubscribeUsage(t *testing.T) {
	var s subscribe
	usage := s.Usage()
	expected := "subscribe BROADCASTER_ID ..."

	if usage != expected {
		t.Errorf("usage: %s; expected: %s", usage, expected)
	}
}

func TestSubscribeExecuteInitFailure(t *testing.T) {
	log := zerolog.Nop()
	s := &subscribe{
		app: &tsm_testing.App{
			InitError: tsm_testing.Error,
			Log:       &log,
		},
		appInitializer: &appInitializer{
			loggerFactory: tsm_testing.NoopLoggerFactory,
		},
	}

	exitCode := s.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitFailure {
		t.Errorf("exitCode: %v; expected: subcommands.ExitFailure", exitCode)
	}
}

func TestSubscribeExecuteSubsSvcNotInitialized(t *testing.T) {
	log := zerolog.Nop()
	s := &subscribe{
		app: &tsm_testing.App{Log: &log},
		appInitializer: &appInitializer{
			loggerFactory: tsm_testing.NoopLoggerFactory,
		},
	}

	exitCode := s.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitFailure {
		t.Errorf("exitCode: %v; expected: subcommands.ExitFailure", exitCode)
	}
}

func TestSubscribeExecuteNoSubs(t *testing.T) {
	f := flag.NewFlagSet("test", flag.ContinueOnError)
	log := zerolog.Nop()
	s := &subscribe{
		app: &tsm_testing.App{
			Log:  &log,
			TOSS: tsm_testing.NewFakeTwitchOnlineSubscriptionService(),
		},
		appInitializer: &appInitializer{
			loggerFactory: tsm_testing.NoopLoggerFactory,
		},
	}

	exitCode := s.Execute(context.Background(), f)

	if exitCode != subcommands.ExitSuccess {
		t.Errorf("exitCode: %v; expected: subcommands.ExitSuccess", exitCode)
	}
}

func TestSubscribeExecuteOK(t *testing.T) {
	f := flag.NewFlagSet("test", flag.ContinueOnError)
	log := zerolog.Nop()
	s := &subscribe{
		app: &tsm_testing.App{
			Log:  &log,
			TOSS: tsm_testing.NewFakeTwitchOnlineSubscriptionService("123", "345 subError"),
		},
		appInitializer: &appInitializer{
			loggerFactory: tsm_testing.NoopLoggerFactory,
		},
	}

	f.Parse([]string{"123", "456", "789"})

	exitCode := s.Execute(context.Background(), f)

	if exitCode != subcommands.ExitSuccess {
		t.Errorf("exitCode: %v; expected: subcommands.ExitSuccess", exitCode)
	}
}
