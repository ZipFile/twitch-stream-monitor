//go:build windows

package utils

import (
	"os"
	"syscall"
)

// Send a ctrl-break to the process.
func Terminate(process *os.Process) error {
	// https://github.com/golang/go/blob/bdefb77309fdc6e47102a8d6272fd2293aefa1d9/src/os/signal/signal_windows_test.go#L18-L31
	dll, err := syscall.LoadDLL("kernel32.dll")

	if err != nil {
		return err
	}

	proc, err := dll.FindProc("GenerateConsoleCtrlEvent")

	if err != nil {
		return err
	}

	result, _, err := proc.Call(syscall.CTRL_BREAK_EVENT, uintptr(process.Pid))

	if result == 0 {
		return err
	}

	return nil
}
