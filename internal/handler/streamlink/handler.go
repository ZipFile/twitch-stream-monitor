package streamlink

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path"
	"time"

	"github.com/rs/zerolog"

	tsm "github.com/ZipFile/twitch-stream-monitor/internal"
	utils "github.com/ZipFile/twitch-stream-monitor/internal/utils"
)

type Handler struct {
	Path        string
	FileDir     string
	LogDir      string
	Log         zerolog.Logger
	KillTimeout time.Duration
}

var ErrBadKillTimeout = errors.New("0 or negative kill timeout is not supported")

// Instantiate event handler utilizing streamlink to record streams.
//
// Args:
// * path: Location of the streamlink executable. Keep empty to use default.
// * fileDir: Location where to store streams. Keep empty to use current working dir.
// * logDir: Location where to store stream recording logs. Keep empty to use current working dir.
// * killTimeout: Time to wait for process to gracefully shutdown before kiling.
// * log: Logger instance.
func New(path, fileDir, logDir string, killTimeout time.Duration, log zerolog.Logger) (tsm.TwitchStreamOnlineEventHandler, error) {
	if killTimeout <= 0 {
		return nil, ErrBadKillTimeout
	}

	return &Handler{
		Path:        utils.OrStr(path, "streamlink"),
		FileDir:     utils.OrStr(fileDir, "."),
		LogDir:      utils.OrStr(logDir, "."),
		Log:         log.With().Str("component", "streamlink_handler").Logger(),
		KillTimeout: killTimeout,
	}, nil
}

// Always returns "streamlink".
func (*Handler) Name() string {
	return "streamlink"
}

// Checks that FileDir and LogDir are writable, and streamlink is executable.
func (h *Handler) Check(ctx context.Context) error {
	err := utils.CheckDirIsWritable(h.FileDir)

	if err != nil {
		return err
	}

	if h.LogDir != h.FileDir {
		err = utils.CheckDirIsWritable(h.LogDir)

		if err != nil {
			return err
		}
	}

	return utils.CheckCLI(ctx, h.Path, "--help")
}

// Start stream recording.
//
// Logs are written to LogDir and the output file to FileDir.
func (h *Handler) Handle(ctx context.Context, event tsm.TwitchStreamOnlineEvent) error {
	killCtx, kill := context.WithCancel(context.Background())
	baseName := fmt.Sprintf(
		"%s %s",
		event.UserLogin,
		event.StartedAt.Format("20060102150405"),
	)
	streamFilename := baseName + ".mp4"
	logFilename := baseName + ".log"
	cmd := exec.CommandContext(
		killCtx,
		h.Path,
		// "--retry-streams", "5",
		// "--retry-max", "10",
		// "--retry-open", "3",
		"--logfile", path.Join(h.LogDir, logFilename),
		"--twitch-disable-ads",
		fmt.Sprintf("twitch.tv/%s", event.UserLogin),
		"best",
		"-o", path.Join(h.FileDir, streamFilename),
	)

	log := h.Log.With().Stringer("cmd", cmd).Logger()
	exit := make(chan interface{})

	defer close(exit)
	log.Debug().Msg("Starting stream recording")

	go func() {
		defer kill()

		for {
			select {
			case <-ctx.Done():
				log.Debug().Msg("Terminating")
				utils.Terminate(cmd.Process)
				// TODO: do not oversleep when successfully terminated
				time.Sleep(h.KillTimeout)
				return
			case <-exit:
				return
			}
		}
	}()

	err := cmd.Run()

	if err == nil {
		log.Debug().Msg("Done!")

		return nil
	}

	if exiterr, ok := err.(*exec.ExitError); ok {
		h.Log.Warn().Int("exit_status", exiterr.ExitCode()).Msg("Streamlink exited")
	} else {
		h.Log.Error().Err(err).Msg("Failed to run streamlink")
	}

	return err
}
