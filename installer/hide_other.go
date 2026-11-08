//go:build !windows

package installer

import "os/exec"

// hideWindow does nothing: only Windows opens a window for a console.
func hideWindow(*exec.Cmd) {}
