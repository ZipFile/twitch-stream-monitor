package cli

import (
	"context"
	"flag"
	"testing"

	"github.com/google/subcommands"
	"github.com/rs/zerolog"

	tsm_testing "github.com/ZipFile/twitch-stream-monitor/internal/testing"
)

func TestResolveUsernameName(t *testing.T) {
	var ru resolveUsername
	name := ru.Name()
	expected := "resolve-username"

	if name != expected {
		t.Errorf("name: %s; expected: %s", name, expected)
	}
}

func TestResolveUsernameSynopsis(t *testing.T) {
	var ru resolveUsername
	synopsis := ru.Synopsis()
	expected := "Find twitch user id by its username."

	if synopsis != expected {
		t.Errorf("synopsis: %s; expected: %s", synopsis, expected)
	}
}

func TestResolveUsernameUsage(t *testing.T) {
	var ru resolveUsername
	usage := ru.Usage()
	expected := "resolve-username USERNAME"

	if usage != expected {
		t.Errorf("usage: %s; expected: %s", usage, expected)
	}
}

func TestResolveUsernameSetFlags(t *testing.T) {
	var ru resolveUsername

	ru.SetFlags(nil)
}

func TestResolveUsernameExecuteInitFailure(t *testing.T) {
	log := zerolog.Nop()
	ru := &resolveUsername{
		app: &tsm_testing.App{
			InitError: tsm_testing.Error,
			Log:       &log,
		},
		appInitializer: &appInitializer{
			loggerFactory: tsm_testing.NoopLoggerFactory,
		},
	}

	exitCode := ru.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitFailure {
		t.Errorf("exitCode: %v; expected: subcommands.ExitFailure", exitCode)
	}
}

func TestResolveUsernameExecutMissingUsername(t *testing.T) {
	f := flag.NewFlagSet("test", flag.ContinueOnError)
	log := zerolog.Nop()
	ru := &resolveUsername{
		app: &tsm_testing.App{
			Log: &log,
		},
		appInitializer: &appInitializer{
			loggerFactory: tsm_testing.NoopLoggerFactory,
		},
	}

	exitCode := ru.Execute(context.Background(), f)

	if exitCode != subcommands.ExitUsageError {
		t.Errorf("exitCode: %v; expected: subcommands.ExitUsageError", exitCode)
	}
}

func TestResolveUsernameExecutGetUsernameResolverFailure(t *testing.T) {
	f := flag.NewFlagSet("test", flag.ContinueOnError)
	log := zerolog.Nop()
	ru := &resolveUsername{
		app: &tsm_testing.App{
			Log: &log,
		},
		appInitializer: &appInitializer{
			loggerFactory: tsm_testing.NoopLoggerFactory,
		},
	}

	f.Parse([]string{"test"})

	exitCode := ru.Execute(context.Background(), f)

	if exitCode != subcommands.ExitFailure {
		t.Errorf("exitCode: %v; expected: subcommands.ExitFailure", exitCode)
	}
}

func TestResolveUsernameExecuteResolveFailure(t *testing.T) {
	f := flag.NewFlagSet("test", flag.ContinueOnError)
	log := zerolog.Nop()
	ru := &resolveUsername{
		app: &tsm_testing.App{
			Log: &log,
			UR: &tsm_testing.UsernameResolver{
				Error: tsm_testing.Error,
			},
		},
		appInitializer: &appInitializer{
			loggerFactory: tsm_testing.NoopLoggerFactory,
		},
	}

	f.Parse([]string{"test"})

	exitCode := ru.Execute(context.Background(), f)

	if exitCode != subcommands.ExitFailure {
		t.Errorf("exitCode: %v; expected: subcommands.ExitFailure", exitCode)
	}
}

func TestResolveUsernameExecuteNoPrint(t *testing.T) {
	f := flag.NewFlagSet("test", flag.ContinueOnError)
	log := zerolog.Nop()
	ru := &resolveUsername{
		app: &tsm_testing.App{
			Log: &log,
			UR:  &tsm_testing.UsernameResolver{},
		},
		appInitializer: &appInitializer{
			loggerFactory: tsm_testing.NoopLoggerFactory,
		},
		print: func(s string) error { panic("should not be called") },
	}

	f.Parse([]string{"test"})

	exitCode := ru.Execute(context.Background(), f)

	if exitCode != subcommands.ExitSuccess {
		t.Errorf("exitCode: %v; expected: subcommands.ExitSuccess", exitCode)
	}
}

func TestResolveUsernameExecutePrintFailure(t *testing.T) {
	f := flag.NewFlagSet("test", flag.ContinueOnError)
	log := zerolog.Nop()
	ru := &resolveUsername{
		app: &tsm_testing.App{
			Log: &log,
			UR: &tsm_testing.UsernameResolver{
				Usernames: map[string]string{
					"test": "123",
				},
			},
		},
		appInitializer: &appInitializer{
			loggerFactory: tsm_testing.NoopLoggerFactory,
		},
		print: func(s string) error { return tsm_testing.Error },
	}

	f.Parse([]string{"test"})

	exitCode := ru.Execute(context.Background(), f)

	if exitCode != subcommands.ExitFailure {
		t.Errorf("exitCode: %v; expected: subcommands.ExitFailure", exitCode)
	}
}

func TestResolveUsernameExecuteOK(t *testing.T) {
	var printedValue string
	f := flag.NewFlagSet("test", flag.ContinueOnError)
	log := zerolog.Nop()
	broadcasterID := "123"
	ru := &resolveUsername{
		app: &tsm_testing.App{
			Log: &log,
			UR: &tsm_testing.UsernameResolver{
				Usernames: map[string]string{
					"test": broadcasterID,
				},
			},
		},
		appInitializer: &appInitializer{
			loggerFactory: tsm_testing.NoopLoggerFactory,
		},
		print: func(s string) error {
			printedValue = s
			return nil
		},
	}

	f.Parse([]string{"test"})

	exitCode := ru.Execute(context.Background(), f)

	if exitCode != subcommands.ExitSuccess {
		t.Errorf("exitCode: %v; expected: subcommands.ExitSuccess", exitCode)
	}

	if printedValue != broadcasterID {
		t.Errorf("printedValue: %v; expected: %v", printedValue, broadcasterID)
	}
}
