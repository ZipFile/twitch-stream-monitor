package cli

import (
	"context"
	"testing"

	"github.com/google/subcommands"
	"github.com/rs/zerolog"

	tsm_testing "github.com/ZipFile/twitch-stream-monitor/internal/testing"
)

func TestListName(t *testing.T) {
	var l list
	name := l.Name()
	expected := "list"

	if name != expected {
		t.Errorf("name: %s; expected: %s", name, expected)
	}
}

func TestListSynopsis(t *testing.T) {
	var l list
	synopsis := l.Synopsis()
	expected := "List stream.online subscriptions."

	if synopsis != expected {
		t.Errorf("synopsis: %s; expected: %s", synopsis, expected)
	}
}

func TestListUsage(t *testing.T) {
	var l list
	usage := l.Usage()
	expected := "list"

	if usage != expected {
		t.Errorf("usage: %s; expected: %s", usage, expected)
	}
}

func TestListSetFlags(t *testing.T) {
	var l list

	l.SetFlags(nil)
}

func TestListExecuteInitFailure(t *testing.T) {
	log := zerolog.Nop()
	l := &list{
		app: &tsm_testing.App{
			InitError: tsm_testing.Error,
			Log:       &log,
		},
		appInitializer: &appInitializer{
			loggerFactory: tsm_testing.NoopLoggerFactory,
		},
	}

	exitCode := l.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitFailure {
		t.Errorf("exitCode: %v; expected: subcommands.ExitFailure", exitCode)
	}
}

func TestListExecuteSubsSvcNotInitialized(t *testing.T) {
	log := zerolog.Nop()
	l := &list{
		app: &tsm_testing.App{Log: &log},
		appInitializer: &appInitializer{
			loggerFactory: tsm_testing.NoopLoggerFactory,
		},
	}

	exitCode := l.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitFailure {
		t.Errorf("exitCode: %v; expected: subcommands.ExitFailure", exitCode)
	}
}

func TestListExecuteListError(t *testing.T) {
	log := zerolog.Nop()
	toss := tsm_testing.NewFakeTwitchOnlineSubscriptionService()
	l := &list{
		app: &tsm_testing.App{
			Log:  &log,
			TOSS: toss,
		},
		appInitializer: &appInitializer{
			loggerFactory: tsm_testing.NoopLoggerFactory,
		},
	}

	toss.ListError = tsm_testing.Error

	exitCode := l.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitFailure {
		t.Errorf("exitCode: %v; expected: subcommands.ExitFailure", exitCode)
	}
}

func TestListExecuteListOK(t *testing.T) {
	log := zerolog.Nop()
	toss := tsm_testing.NewFakeTwitchOnlineSubscriptionService("123", "456")
	l := &list{
		app: &tsm_testing.App{
			Log:  &log,
			TOSS: toss,
		},
		appInitializer: &appInitializer{
			loggerFactory: tsm_testing.NoopLoggerFactory,
		},
	}

	toss.Subscribe("123")
	toss.Subscribe("456")

	exitCode := l.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitSuccess {
		t.Errorf("exitCode: %v; expected: subcommands.ExitSuccess", exitCode)
	}

	// TODO: assert on logger output
}
