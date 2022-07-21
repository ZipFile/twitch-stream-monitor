package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/google/subcommands"
	"github.com/rs/zerolog"

	tsm "github.com/ZipFile/twitch-stream-monitor/internal"
	tsm_testing "github.com/ZipFile/twitch-stream-monitor/internal/testing"
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

func TestInitFailure(t *testing.T) {
	initError := errors.New("test")
	m := &monitor{
		app: &testApp{
			init: initError,
		},
		loggerFactory: tsm_testing.NoopLoggerFactory,
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
		loggerFactory: tsm_testing.PanicLoggerFactory,
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
		loggerFactory: tsm_testing.NoopLoggerFactory,
	}

	exitCode := m.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitFailure {
		t.Errorf("exitCode: %v; expected: subcommands.ExitFailure", exitCode)
	}
}

func TestRunSucces(t *testing.T) {
	m := &monitor{
		app:           &testApp{},
		loggerFactory: tsm_testing.NoopLoggerFactory,
	}

	exitCode := m.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitSuccess {
		t.Errorf("exitCode: %v; expected: subcommands.ExitSuccess", exitCode)
	}
}
