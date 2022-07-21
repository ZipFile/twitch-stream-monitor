package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/google/subcommands"
	"github.com/rs/zerolog"

	tsm "github.com/ZipFile/twitch-stream-monitor/internal"
)

type testApp struct {
	init error
	run  error
	toss tsm.TwitchOnlineSubscriptionService
	log  *zerolog.Logger
}

func (a *testApp) Init() error {
	return a.init
}

func (a *testApp) Run() error {
	return a.run
}

func (a *testApp) GetTwitchOnlineSubscriptionService() tsm.TwitchOnlineSubscriptionService {
	return a.toss
}

func (a *testApp) GetLogger() *zerolog.Logger {
	return a.log
}

func noopLoggerFactory() *zerolog.Logger {
	log := zerolog.Nop()

	return &log
}

func panicLoggerFactory() *zerolog.Logger {
	panic("should not be called")

	return nil
}

func TestInitFailure(t *testing.T) {
	initError := errors.New("test")
	m := &monitor{
		app: &testApp{
			init: initError,
		},
		loggerFactory: noopLoggerFactory,
	}

	exitCode := m.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitFailure {
		t.Errorf("exitCode: %v; expected: subcommands.ExitFailure", exitCode)
	}
}

func TestInitFailureWithLogger(t *testing.T) {
	initError := errors.New("test")
	log := zerolog.Nop()
	m := &monitor{
		app: &testApp{
			init: initError,
			log:  &log,
		},
		loggerFactory: panicLoggerFactory,
	}

	exitCode := m.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitFailure {
		t.Errorf("exitCode: %v; expected: subcommands.ExitFailure", exitCode)
	}
}

func TestRunFailure(t *testing.T) {
	runError := errors.New("test")
	m := &monitor{
		app: &testApp{
			run: runError,
		},
		loggerFactory: noopLoggerFactory,
	}

	exitCode := m.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitFailure {
		t.Errorf("exitCode: %v; expected: subcommands.ExitFailure", exitCode)
	}
}

func TestRunSucces(t *testing.T) {
	m := &monitor{
		app:           &testApp{},
		loggerFactory: noopLoggerFactory,
	}

	exitCode := m.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitSuccess {
		t.Errorf("exitCode: %v; expected: subcommands.ExitSuccess", exitCode)
	}
}
