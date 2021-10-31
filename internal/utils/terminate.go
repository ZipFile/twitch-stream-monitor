//go:build !windows

package utils

import (
	"os"
	"syscall"
)

// Send SIGTERM signal to the proccess.
func Terminate(process *os.Process) error {
	return process.Signal(syscall.SIGTERM)
}
