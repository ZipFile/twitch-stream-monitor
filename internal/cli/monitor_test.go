package cli

import (
	"context"
	"testing"

	"github.com/google/subcommands"
	"github.com/rs/zerolog"

	tsm_testing "github.com/ZipFile/twitch-stream-monitor/internal/testing"
)

func TestMonitorName(t *testing.T) {
	var m monitor
	name := m.Name()
	expected := "monitor"

	if name != expected {
		t.Errorf("name: %s; expected: %s", name, expected)
	}
}

func TestMonitorSynopsis(t *testing.T) {
	var m monitor
	synopsis := m.Synopsis()
	expected := "Launch Twitch Stream Monitor."

	if synopsis != expected {
		t.Errorf("synopsis: %s; expected: %s", synopsis, expected)
	}
}

func TestMonitorUsage(t *testing.T) {
	var m monitor
	usage := m.Usage()
	expected := "monitor"

	if usage != expected {
		t.Errorf("usage: %s; expected: %s", usage, expected)
	}
}

func TestMonitorSetFlags(t *testing.T) {
	var m monitor

	m.SetFlags(nil)
}

func TestMonitorExecuteInitFailure(t *testing.T) {
	log := zerolog.Nop()
	m := &monitor{
		app: &tsm_testing.App{
			InitError: tsm_testing.Error,
			Log:       &log,
		},
		appInitializer: &appInitializer{
			loggerFactory: tsm_testing.NoopLoggerFactory,
		},
	}

	exitCode := m.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitFailure {
		t.Errorf("exitCode: %v; expected: subcommands.ExitFailure", exitCode)
	}
}

func TestMonitorExecuteRunFailure(t *testing.T) {
	log := zerolog.Nop()
	m := &monitor{
		app: &tsm_testing.App{
			RunError: tsm_testing.Error,
			Log:      &log,
		},
		appInitializer: &appInitializer{
			loggerFactory: tsm_testing.NoopLoggerFactory,
		},
	}

	exitCode := m.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitFailure {
		t.Errorf("exitCode: %v; expected: subcommands.ExitFailure", exitCode)
	}
}

func TestMonitorExecuteOK(t *testing.T) {
	log := zerolog.Nop()
	m := &monitor{
		app: &tsm_testing.App{Log: &log},
		appInitializer: &appInitializer{
			loggerFactory: tsm_testing.NoopLoggerFactory,
		},
	}

	exitCode := m.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitSuccess {
		t.Errorf("exitCode: %v; expected: subcommands.ExitSuccess", exitCode)
	}
}
