package cli

import (
	"context"
	"flag"
	"reflect"
	"testing"

	"github.com/google/subcommands"
	"github.com/rs/zerolog"

	tsm_testing "github.com/ZipFile/twitch-stream-monitor/internal/testing"
)

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

	exitCode := s.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitSuccess {
		t.Errorf("exitCode: %v; expected: subcommands.ExitSuccess", exitCode)
	}
}

func TestSubscribeExecuteOK(t *testing.T) {
	log := zerolog.Nop()
	s := &subscribe{
		app: &tsm_testing.App{
			Log:  &log,
			TOSS: tsm_testing.NewFakeTwitchOnlineSubscriptionService("123", "345 subError"),
		},
		appInitializer: &appInitializer{
			loggerFactory: tsm_testing.NoopLoggerFactory,
		},
		broadcasterIDs: []string{"123", "456", "789"},
	}

	exitCode := s.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitSuccess {
		t.Errorf("exitCode: %v; expected: subcommands.ExitSuccess", exitCode)
	}
}

func TestSubscribeSetFlags(t *testing.T) {
	f := flag.NewFlagSet("test", flag.PanicOnError)

	f.Parse([]string{"123", "456"})

	s := &subscribe{}

	s.SetFlags(f)

	if !reflect.DeepEqual(s.broadcasterIDs, []string{"123", "456"}) {
		t.Errorf("s.broadcasterIDs: %v: expected: [\"123\", \"456\"]", s.broadcasterIDs)
	}
}
