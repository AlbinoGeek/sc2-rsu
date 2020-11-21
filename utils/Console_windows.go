// +build windows

package utils

import (
	"fmt"
	"os"
	"syscall"

	"github.com/kataras/golog"
	"golang.org/x/sys/windows"
)

// win32api constants
const (
	ATTACH_PARENT_PROCESS = ^uintptr(0)
	ERROR_INVALID_HANDLE  = 6
)

// to handle redirection
var (
	prevStderr *os.File
	prevStdin  *os.File
	prevStdout *os.File
)

// AttachConsole gives us access to the parent console on Windows, when built
// as a GUI application, where the console would otherwise not be available.
func AttachConsole() error {
	proc := syscall.MustLoadDLL("kernel32.dll").MustFindProc("AttachConsole")
	r1, _, err := proc.Call(ATTACH_PARENT_PROCESS)
	if r1 == 0 {
		if errno, is := err.(syscall.Errno); !is || errno != ERROR_INVALID_HANDLE {
			return err
		}
	}

	return nil
}

// FixRedirection restores the parent console's ability to redirect this
// program's output through a pipe to a file or device, where without
// doing so, the output would be forced to the console even if redirected.
func FixRedirection() error {
	prevStderr, prevStdin, prevStdout = os.Stderr, os.Stdin, os.Stdout
	stderr, stdin, stdout, err := sysGetHandles()

	// if any handles are invalid we need to AttachConsole first
	var invalid syscall.Handle
	if err != nil || stderr == invalid || stdin == invalid || stdout == invalid {
		if err = AttachConsole(); err != nil {
			return fmt.Errorf("failed to AttachConsole: %v", err)
		}

		if stderr == invalid {
			stderr, _ = syscall.GetStdHandle(syscall.STD_ERROR_HANDLE)
		}
		if stdin == invalid {
			stdin, _ = syscall.GetStdHandle(syscall.STD_INPUT_HANDLE)
		}
		if stdout == invalid {
			stdout, _ = syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
		}
	}

	// golang expects line ending conversion handled by the writer itself
	for _, c := range []syscall.Handle{stderr, stdout} {
		if c != invalid {
			if err = imposeLineEndings(windows.Handle(c)); err != nil {
				return err
			}
		}
	}

	// assign newly acquired file handles
	if stderr != invalid {
		os.Stderr = os.NewFile(uintptr(stderr), "stderr")
	}
	if stdin != invalid {
		os.Stdin = os.NewFile(uintptr(stdin), "stdin")
	}
	if stdout != invalid {
		os.Stdout = os.NewFile(uintptr(stdout), "stdout")
	}

	return nil
}

// FreeConsole would relinquish our control of the parent console on Windows
// ? but is it actually necessary ? requires more testing
// TODO: NOT YET IMPLEMENTED
func FreeConsole() error {
	golog.Infof("FreeConsole: not yet implemented")
	return fmt.Errorf("not yet implemented")
}

func imposeLineEndings(h windows.Handle) error {
	var st uint32

	if err := windows.GetConsoleMode(h, &st); err != nil {
		return fmt.Errorf("GetConsoleMode: %v", err)
	}
	if err := windows.SetConsoleMode(h, st&^windows.DISABLE_NEWLINE_AUTO_RETURN); err != nil {
		return fmt.Errorf("SetConsoleMode: %v", err)
	}

	return nil
}

func sysGetHandles() (stderr, stdin, stdout syscall.Handle, err error) {
	stderr, err = syscall.GetStdHandle(syscall.STD_ERROR_HANDLE)
	if err != nil {
		stdin, err = syscall.GetStdHandle(syscall.STD_INPUT_HANDLE)
		if err != nil {
			stdout, err = syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
		}
	}

	return
}
