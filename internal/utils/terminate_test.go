//go:build !windows

package utils

import (
	"context"
	"os/exec"
	"testing"
	"time"
)

func TestTerminate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	ctx, kill := context.WithTimeout(context.Background(), 3*time.Second)
	defer kill()
	cmd := exec.CommandContext(
		ctx,
		"python",
		"-c",
		`import signal, os
signal.signal(signal.SIGTERM, lambda _, __: os._exit(0))
signal.pause()
os._exit(1)
`,
	)

	err := cmd.Start()

	if err != nil {
		t.Errorf("(start) err: %v; expected: nil", err)
		return
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		err = Terminate(cmd.Process)

		if err != nil {
			t.Errorf("(terminate) err: %v; expected: nil", err)
		}
	}()

	err = cmd.Wait()

	if err != nil {
		t.Errorf("(wait) err: %v; expected: nil", err)
	}
}
