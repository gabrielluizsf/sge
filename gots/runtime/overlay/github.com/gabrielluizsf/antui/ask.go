//go:build gots_overlay

// Overlay for ask.go: the program runs synchronously, as there are no
// goroutines to time it out from.

package antui

import (
	"os/exec"
	"strings"
)

func ask(name string, args ...string) string {
	path, err := exec.LookPath(name)
	if err != nil {
		return ""
	}
	out, err := exec.Command(path, args...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
