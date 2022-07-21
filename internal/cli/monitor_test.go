package cli

import (
	"context"
	"testing"

	"github.com/google/subcommands"
	"github.com/rs/zerolog"

	tsm_testing "github.com/ZipFile/twitch-stream-monitor/internal/testing"
)

func TestMonitorExecuteInitFailure(t *testing.T) {
	m := &monitor{
		app: &tsm_testing.App{
			InitError: tsm_testing.Error,
		},
		loggerFactory: tsm_testing.NoopLoggerFactory,
	}

	exitCode := m.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitFailure {
		t.Errorf("exitCode: %v; expected: subcommands.ExitFailure", exitCode)
	}
}

func TestMonitorExecuteInitFailureWithLogger(t *testing.T) {
	log := zerolog.Nop()
	m := &monitor{
		app: &tsm_testing.App{
			InitError: tsm_testing.Error,
			Log:       &log,
		},
		loggerFactory: tsm_testing.PanicLoggerFactory,
	}

	exitCode := m.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitFailure {
		t.Errorf("exitCode: %v; expected: subcommands.ExitFailure", exitCode)
	}
}

func TestMonitorExecuteRunFailure(t *testing.T) {
	m := &monitor{
		app: &tsm_testing.App{
			RunError: tsm_testing.Error,
		},
		loggerFactory: tsm_testing.NoopLoggerFactory,
	}

	exitCode := m.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitFailure {
		t.Errorf("exitCode: %v; expected: subcommands.ExitFailure", exitCode)
	}
}

func TestMonitorExecuteOK(t *testing.T) {
	m := &monitor{
		app:           &tsm_testing.App{},
		loggerFactory: tsm_testing.NoopLoggerFactory,
	}

	exitCode := m.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitSuccess {
		t.Errorf("exitCode: %v; expected: subcommands.ExitSuccess", exitCode)
	}
}
