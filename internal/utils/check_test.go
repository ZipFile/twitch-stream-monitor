package utils

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"testing"
	"time"
)

func TestCheckDirIsWritableOK(t *testing.T) {
	dir, err := ioutil.TempDir("", "test")

	if err != nil {
		t.Error(err)
		return
	}

	defer os.RemoveAll(dir)

	err = CheckDirIsWritable(dir)

	if err != nil {
		t.Errorf("err: %v; expected: nil", err)
	}
}

func TestCheckDirIsWritableError(t *testing.T) {
	dir, err := ioutil.TempDir("", "test")

	if err != nil {
		t.Error(err)
		return
	}

	os.Chmod(dir, 0555)

	defer func() {
		os.Chmod(dir, 0777)
		os.RemoveAll(dir)
	}()

	err = CheckDirIsWritable(dir)

	if err == nil {
		t.Errorf("err: nil; expected: error")
	}
}

func TestCheckFileIsReadableOK(t *testing.T) {
	f := path.Join(t.TempDir(), "test")
	err := ioutil.WriteFile(f, []byte("test"), 0644)

	if err != nil {
		t.Error(err)
		return
	}

	err = CheckFileIsReadable(f)

	if err != nil {
		t.Errorf("err: %v; expected: nil", err)
	}
}

func TestCheckFileIsReadableError(t *testing.T) {
	err := CheckFileIsReadable(path.Join(t.TempDir(), "test"))

	if err == nil {
		t.Errorf("err: nil; expected: error")
	}
}

func TestCheckCLIOK(t *testing.T) {
	err := CheckCLI(context.Background(), "echo", "test")

	if err != nil {
		t.Errorf("err: %v; expected: nil", err)
	}
}

func TestCheckCLIExit(t *testing.T) {
	err := CheckCLI(context.Background(), "exit", "1")

	if err == nil {
		t.Errorf("err: nil; expected: error")
	}
}

func TestCheckCLINotFound(t *testing.T) {
	err := CheckCLI(context.Background(), fmt.Sprintf("test-%d", rand.Int()))

	if err == nil {
		t.Errorf("err: nil; expected: error")
	}
}

func TestCheckCLITimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)

	defer cancel()

	err := CheckCLI(ctx, "sleep", "10")

	if err == nil {
		t.Errorf("err: nil; expected: error")
	}
}
