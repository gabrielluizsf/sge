//go:build windows

package installer

import (
	"os"
	"path/filepath"
	"testing"
)

// testPlace is a new home on Windows: its own Local and Roaming folders,
// the system's own programs where they are.
func testPlace(t *testing.T) Place {
	home := t.TempDir()
	local := filepath.Join(home, "AppData", "Local")
	return Place{System: Windows, Home: home, Data: local, Roaming: filepath.Join(home, "AppData", "Roaming"), SystemRoot: os.Getenv("SystemRoot"), Caches: local}
}
