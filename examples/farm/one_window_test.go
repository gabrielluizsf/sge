package farm

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Each game opens one window, once, and shows its scenes in it: the
// loading screen, the menu, the character's making and the play follow
// one another without the window closing and opening again between them.
func TestEachGameOpensOneWindow(t *testing.T) {
	for dir, want := range map[string]int{".": 1, "../realfarm": 1} {
		opens, runs := 0, 0
		files, _ := filepath.Glob(filepath.Join(dir, "*.go"))
		for _, f := range files {
			if strings.HasSuffix(f, "_test.go") {
				continue
			}
			b, err := os.ReadFile(f)
			if err != nil {
				t.Fatal(err)
			}
			src := string(b)
			opens += strings.Count(src, "antui.Open(")
			runs += strings.Count(src, ".Game.Run(") + strings.Count(src, "g.Run(")
		}
		if opens != want {
			t.Errorf("%s opens a window in %d places, not %d", dir, opens, want)
		}
		if dir == "../realfarm" && runs != 0 {
			t.Errorf("%s runs a game in a window of its own %d times", dir, runs)
		}
	}
}
