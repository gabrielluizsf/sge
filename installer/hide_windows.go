//go:build windows

package installer

import (
	"os/exec"
	"syscall"
)

// hideWindow keeps a console program from opening a window of its own.
func hideWindow(cmd *exec.Cmd) {
	const createNoWindow = 0x08000000
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: createNoWindow}
}
