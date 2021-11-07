package utils

import (
	"context"
	"os"
	"os/exec"
)

// Check if path is a directory and is writable.
//
// Attempts to create a file at the given path and then writes "test" text.
// Returns nil on success. Created file is removed afterwards.
func CheckDirIsWritable(path string) error {
	file, err := os.CreateTemp(path, "DeleteMe.*.txt")

	if err != nil {
		return err
	}

	defer os.Remove(file.Name())

	_, err = file.Write([]byte("test"))

	if err != nil {
		file.Close()

		return err
	}

	return file.Close()
}

// Check if path is a file and is readablr.
//
// Returns nil on success.
func CheckFileIsReadable(path string) error {
	f, err := os.Open(path)

	if err != nil {
		return err
	}

	buff := make([]byte, 16)
	_, err = f.Read(buff)

	return err
}

// Check if executable at path exists and is executable.
//
// Calls path with provided args. Returns nil on success.
// If proccess takes more than 5 seconds to execute, it will be killed and error
// returned.
func CheckCLI(ctx context.Context, path string, args ...string) error {
	cmd := exec.CommandContext(ctx, path, args...)

	return cmd.Run()
}
