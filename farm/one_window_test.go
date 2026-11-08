package farm

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The game opens one window, once, and shows its scenes in it: the
// loading screen, the menu, the character's making and the play follow
// one another without the window closing and opening again between them.
func TestTheGameOpensOneWindow(t *testing.T) {
	opens := 0
	files, _ := filepath.Glob("*.go")
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		opens += strings.Count(string(b), "antui.Open(")
	}
	if opens != 1 {
		t.Errorf("package farm opens a window in %d places, not 1", opens)
	}
}
