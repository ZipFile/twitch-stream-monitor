package cli

import (
	"context"
	"testing"
	"time"

	"github.com/google/subcommands"
	"github.com/rs/zerolog"

	tsm_aat "github.com/ZipFile/twitch-stream-monitor/internal/app_access_token"
	tsm_testing "github.com/ZipFile/twitch-stream-monitor/internal/testing"
)

func TestGetAppAccessTokenExecuteInitFailure(t *testing.T) {
	log := zerolog.Nop()
	m := &getAppAccessToken{
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

func TestGetAppAccessTokenExecuteGetTokenStoreFailure(t *testing.T) {
	log := zerolog.Nop()
	m := &getAppAccessToken{
		app: &tsm_testing.App{
			Log: &log,
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

func TestGetAppAccessTokenExecuteLoadTokenFailure(t *testing.T) {
	log := zerolog.Nop()
	m := &getAppAccessToken{
		app: &tsm_testing.App{
			Log:  &log,
			AATS: &tsm_aat.NoopStore{LoadError: tsm_testing.Error},
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

func TestGetAppAccessTokenExecuteExpiredToken(t *testing.T) {
	now := time.Date(2022, time.July, 23, 12, 0, 0, 0, time.UTC)
	log := zerolog.Nop()
	m := &getAppAccessToken{
		app: &tsm_testing.App{
			Log: &log,
			AATS: &tsm_aat.NoopStore{
				Expires: time.Date(2022, time.July, 23, 6, 0, 0, 0, time.UTC),
			},
		},
		appInitializer: &appInitializer{
			loggerFactory: tsm_testing.NoopLoggerFactory,
		},
		now: func() time.Time { return now },
	}

	exitCode := m.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitFailure {
		t.Errorf("exitCode: %v; expected: subcommands.ExitFailure", exitCode)
	}
}

func TestGetAppAccessTokenExecutePrintFailure(t *testing.T) {
	now := time.Date(2022, time.July, 23, 12, 0, 0, 0, time.UTC)
	log := zerolog.Nop()
	m := &getAppAccessToken{
		app: &tsm_testing.App{
			Log: &log,
			AATS: &tsm_aat.NoopStore{
				Expires: time.Date(2022, time.July, 24, 12, 0, 0, 0, time.UTC),
			},
		},
		appInitializer: &appInitializer{
			loggerFactory: tsm_testing.NoopLoggerFactory,
		},
		now:   func() time.Time { return now },
		print: func(s string) error { return tsm_testing.Error },
	}

	exitCode := m.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitFailure {
		t.Errorf("exitCode: %v; expected: subcommands.ExitFailure", exitCode)
	}
}

func TestGetAppAccessTokenExecuteOK(t *testing.T) {
	var printedValue string
	token := "test123"
	now := time.Date(2022, time.July, 23, 12, 0, 0, 0, time.UTC)
	log := zerolog.Nop()
	m := &getAppAccessToken{
		app: &tsm_testing.App{
			Log: &log,
			AATS: &tsm_aat.NoopStore{
				Token:   token,
				Expires: time.Date(2022, time.July, 24, 12, 0, 0, 0, time.UTC),
			},
		},
		appInitializer: &appInitializer{
			loggerFactory: tsm_testing.NoopLoggerFactory,
		},
		now: func() time.Time { return now },
		print: func(s string) error {
			printedValue = s
			return nil
		},
	}

	exitCode := m.Execute(context.Background(), nil)

	if exitCode != subcommands.ExitSuccess {
		t.Errorf("exitCode: %v; expected: subcommands.ExitSuccess", exitCode)
	}

	if printedValue != token {
		t.Errorf("printedValue: %v; expected: %v", printedValue, token)
	}
}
