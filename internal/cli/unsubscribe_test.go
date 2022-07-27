package cli

import (
	"context"
	"flag"
	"testing"

	"github.com/google/subcommands"
	"github.com/rs/zerolog"

	tsm_testing "github.com/ZipFile/twitch-stream-monitor/internal/testing"
)

func TestUnsubscribeName(t *testing.T) {
	var u unsubscribe
	name := u.Name()
	expected := "unsubscribe"

	if name != expected {
		t.Errorf("name: %s; expected: %s", name, expected)
	}
}

func TestUnsubscribeSynopsis(t *testing.T) {
	var u unsubscribe
	synopsis := u.Synopsis()
	expected := "Unsubscribe from stream.online by subscription id."

	if synopsis != expected {
		t.Errorf("synopsis: %s; expected: %s", synopsis, expected)
	}
}

func TestUnsubscribeUsage(t *testing.T) {
	var u unsubscribe
	usage := u.Usage()
	expected := "unsubscribe SUB_ID ..."

	if usage != expected {
		t.Errorf("usage: %s; expected: %s", usage, expected)
	}
}

func TestUnsubscribeExecuteInitFailure(t *testing.T) {
	log := zerolog.Nop()
	i := &unsubscribe{
		app: &tsm_testing.App{
			InitError: tsm_testing.Error,
			Log:       &log,
		},
		appInitializer: &appInitializer{
			loggerFactory: tsm_testing.NoopLoggerFactory,
		},
	}

	exitCode := i.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitFailure {
		t.Errorf("exitCode: %v; expected: subcommands.ExitFailure", exitCode)
	}
}

func TestUnsubscribeExecuteSubsSvcNotInitialized(t *testing.T) {
	log := zerolog.Nop()
	u := &unsubscribe{
		app: &tsm_testing.App{Log: &log},
		appInitializer: &appInitializer{
			loggerFactory: tsm_testing.NoopLoggerFactory,
		},
	}

	exitCode := u.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitFailure {
		t.Errorf("exitCode: %v; expected: subcommands.ExitFailure", exitCode)
	}
}

func TestUnsubscribeExecuteNoSubs(t *testing.T) {
	f := flag.NewFlagSet("test", flag.ContinueOnError)
	log := zerolog.Nop()
	u := &unsubscribe{
		app: &tsm_testing.App{
			Log:  &log,
			TOSS: tsm_testing.NewFakeTwitchOnlineSubscriptionService(),
		},
		appInitializer: &appInitializer{
			loggerFactory: tsm_testing.NoopLoggerFactory,
		},
	}

	exitCode := u.Execute(context.Background(), f)

	if exitCode != subcommands.ExitSuccess {
		t.Errorf("exitCode: %v; expected: subcommands.ExitSuccess", exitCode)
	}
}

func TestUnsubscribeExecuteOK(t *testing.T) {
	f := flag.NewFlagSet("test", flag.ContinueOnError)
	log := zerolog.Nop()
	toss := tsm_testing.NewFakeTwitchOnlineSubscriptionService("123", "456 unsubError")

	toss.Subscribe("123")
	toss.Subscribe("456")

	u := &unsubscribe{
		app: &tsm_testing.App{Log: &log, TOSS: toss},
		appInitializer: &appInitializer{
			loggerFactory: tsm_testing.NoopLoggerFactory,
		},
	}

	f.Parse([]string{toss.Sbus["123"], toss.Sbus["456"]})

	exitCode := u.Execute(context.Background(), f)

	if exitCode != subcommands.ExitSuccess {
		t.Errorf("exitCode: %v; expected: subcommands.ExitSuccess", exitCode)
	}
}
