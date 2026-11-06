package farm

import (
	"os"
	"testing"
)

// TestMain gives the tests a config folder of their own: a menu may keep
// a choice or a save, never on the player's computer.
func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "farm-config-")
	if err != nil {
		panic(err)
	}
	os.Setenv("XDG_CONFIG_HOME", dir)
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}
