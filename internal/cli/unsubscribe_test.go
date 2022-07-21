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

func TestUnsubscribeExecuteInitFailure(t *testing.T) {
	u := &unsubscribe{
		app: &tsm_testing.App{
			InitError: tsm_testing.Error,
		},
		loggerFactory: tsm_testing.NoopLoggerFactory,
	}

	exitCode := u.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitFailure {
		t.Errorf("exitCode: %v; expected: subcommands.ExitFailure", exitCode)
	}
}

func TestUnsubscribeExecuteInitFailureWithLogger(t *testing.T) {
	log := zerolog.Nop()
	u := &unsubscribe{
		app: &tsm_testing.App{
			InitError: tsm_testing.Error,
			Log:       &log,
		},
		loggerFactory: tsm_testing.PanicLoggerFactory,
	}

	exitCode := u.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitFailure {
		t.Errorf("exitCode: %v; expected: subcommands.ExitFailure", exitCode)
	}
}

func TestUnsubscribeExecuteSubsSvcNotInitialized(t *testing.T) {
	u := &unsubscribe{
		app:           &tsm_testing.App{},
		loggerFactory: tsm_testing.NoopLoggerFactory,
	}

	exitCode := u.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitFailure {
		t.Errorf("exitCode: %v; expected: subcommands.ExitFailure", exitCode)
	}
}

func TestUnsubscribeExecuteNoSubs(t *testing.T) {
	u := &unsubscribe{
		app: &tsm_testing.App{
			TOSS: tsm_testing.NewFakeTwitchOnlineSubscriptionService(),
		},
		loggerFactory: tsm_testing.NoopLoggerFactory,
	}

	exitCode := u.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitSuccess {
		t.Errorf("exitCode: %v; expected: subcommands.ExitSuccess", exitCode)
	}
}

func TestUnsubscribeExecuteOK(t *testing.T) {
	toss := tsm_testing.NewFakeTwitchOnlineSubscriptionService("123", "345 unsubError")
	u := &unsubscribe{
		app:           &tsm_testing.App{TOSS: toss},
		loggerFactory: tsm_testing.NoopLoggerFactory,
		subIDs:        []string{toss.Sbus["123"], toss.Sbus["345"]},
	}

	exitCode := u.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitSuccess {
		t.Errorf("exitCode: %v; expected: subcommands.ExitSuccess", exitCode)
	}
}

func TestUnsubscribeSetFlags(t *testing.T) {
	f := flag.NewFlagSet("test", flag.PanicOnError)

	f.Parse([]string{"123", "456"})

	u := &unsubscribe{}

	u.SetFlags(f)

	if !reflect.DeepEqual(u.subIDs, []string{"123", "456"}) {
		t.Errorf("u.subsIDs: %v: expected: [\"123\", \"456\"]", u.subIDs)
	}
}
