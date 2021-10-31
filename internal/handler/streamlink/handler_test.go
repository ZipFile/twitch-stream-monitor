package streamlink

import (
	"context"
	"io/ioutil"
	"os/exec"
	"path"
	"testing"
	"time"

	"github.com/rs/zerolog"

	tsm "github.com/ZipFile/twitch-stream-monitor/internal"
)

func TestHandlerName(t *testing.T) {
	var h *Handler

	name := h.Name()

	if name != "streamlink" {
		t.Errorf("name: %v; expected: \"streamlink\"", name)
	}
}

func TestHandlerNegativeKillTimeout(t *testing.T) {
	_, err := New(
		"./fake_streamlink.py",
		"",
		"",
		-1*time.Second,
		zerolog.Nop(),
	)

	if err != ErrBadKillTimeout {
		t.Errorf("err: nil; expected: ErrBadKillTimeout")
	}
}

func TestHandlerCheck(t *testing.T) {
	h, err := New(
		"./fake_streamlink.py",
		t.TempDir(),
		t.TempDir(),
		1*time.Second,
		zerolog.Nop(),
	)

	if err != nil {
		t.Error(err)
		return
	}

	err = h.Check(context.Background())

	if err != nil {
		t.Error(err)
	}
}

func TestHandlerCheckSameDir(t *testing.T) {
	dir := t.TempDir()
	h, err := New(
		"./fake_streamlink.py",
		dir,
		dir,
		1*time.Second,
		zerolog.Nop(),
	)

	if err != nil {
		t.Error(err)
		return
	}

	err = h.Check(context.Background())

	if err != nil {
		t.Error(err)
	}
}

func TestHandlerCheckFailFileDir(t *testing.T) {
	dir := path.Join(t.TempDir(), "does not exists")
	h, err := New(
		"./fake_streamlink.py",
		dir,
		dir,
		1*time.Second,
		zerolog.Nop(),
	)

	if err != nil {
		t.Error(err)
		return
	}

	err = h.Check(context.Background())

	if err == nil {
		t.Error("err: nil; expect: error")
	}
}

func TestHandlerCheckFailLogDir(t *testing.T) {
	h, err := New(
		"./fake_streamlink.py",
		t.TempDir(),
		path.Join(t.TempDir(), "does not exists"),
		1*time.Second,
		zerolog.Nop(),
	)

	if err != nil {
		t.Error(err)
		return
	}

	err = h.Check(context.Background())

	if err == nil {
		t.Error("err: nil; expect: error")
	}
}

func TestHandlerHandleOK(t *testing.T) {
	fileDir := t.TempDir()
	logDir := t.TempDir()
	event := tsm.TwitchStreamOnlineEvent{
		UserID:    "123",
		UserLogin: "test",
		UserName:  "Test",
		StartedAt: time.Date(2021, time.October, 28, 21, 56, 0, 0, time.UTC),
	}
	h, err := New(
		"./fake_streamlink.py",
		fileDir,
		logDir,
		1*time.Second,
		zerolog.Nop(),
	)

	if err != nil {
		t.Error(err)
	}

	err = h.Handle(context.Background(), event)

	if err != nil {
		t.Error(err)
		return
	}

	file, err := ioutil.ReadFile(path.Join(fileDir, "test 20211028215600.mp4"))

	if err == nil {
		if string(file) != "test stream" {
			t.Errorf("string(file): %v; expected: \"test stream\"", string(file))
		}
	} else {
		t.Error(err)
	}

	log, err := ioutil.ReadFile(path.Join(logDir, "test 20211028215600.log"))

	if err == nil {
		if string(log) != "test log" {
			t.Errorf("string(log): %v; expected: \"test log\"", string(log))
		}
	} else {
		t.Error(err)
	}
}

func TestHandlerHandleGracefulShutdown(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	fileDir := t.TempDir()
	logDir := t.TempDir()
	event := tsm.TwitchStreamOnlineEvent{
		UserID:    "123",
		UserLogin: "wait",
		UserName:  "Wait",
		StartedAt: time.Date(2021, time.October, 29, 18, 18, 0, 0, time.UTC),
	}
	h, err := New(
		"./fake_streamlink.py",
		fileDir,
		logDir,
		1*time.Second,
		zerolog.Nop(),
	)

	if err != nil {
		t.Error(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()

	err = h.Handle(ctx, event)

	if err != nil {
		t.Error(err)
	}
}

func TestHandlerHandleKill(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	fileDir := t.TempDir()
	logDir := t.TempDir()
	event := tsm.TwitchStreamOnlineEvent{
		UserID:    "123",
		UserLogin: "wait_forever",
		UserName:  "Wait_Forever",
		StartedAt: time.Date(2021, time.October, 29, 18, 18, 0, 0, time.UTC),
	}
	h, err := New(
		"./fake_streamlink.py",
		fileDir,
		logDir,
		1*time.Second,
		zerolog.Nop(),
	)

	if err != nil {
		t.Error(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()

	err = h.Handle(ctx, event)

	if err == nil {
		t.Errorf("err: nil; expected: err")
		return
	}

	exiterr, ok := err.(*exec.ExitError)

	if !ok {
		t.Errorf("exiterr: nil; expected: ExitError")
		return
	}

	if exiterr.ExitCode() != -1 {
		t.Errorf("exiterr.ExitCode(): %d; expected: -1", exiterr.ExitCode())
	}
}
