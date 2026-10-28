package gots

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestPointerMethodsAreCalledDirectly guards the speed of method calls:
// Class.prototype.M.call(recv, ...) costs V8 three to five times a plain
// call, so Go code calls a struct's pointer methods through their
// standalone functions, nil receivers included (testdata/conformance/nilrecv
// checks they still behave).
func TestPointerMethodsAreCalledDirectly(t *testing.T) {
	out := t.TempDir()
	if _, err := Transpile(Options{Dir: ".", Out: out, Patterns: []string{"./testdata/conformance/nilrecv", "github.com/gabrielluizsf/antui/canvas"}}); err != nil {
		t.Fatalf("transpile: %v", err)
	}
	filepath.WalkDir(out, func(p string, d os.DirEntry, err error) error {
		if err != nil || !strings.HasSuffix(p, ".ts") || strings.Contains(p, "_gots"+string(filepath.Separator)+"std") || strings.Contains(p, "_gots"+string(filepath.Separator)+"rt") {
			return err
		}
		src, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		if i := strings.Index(string(src), ".prototype."); i >= 0 && strings.Contains(string(src[i:min(len(src), i+80)]), ".call(") {
			t.Errorf("%s calls a method through its prototype: %s", p, src[i:min(len(src), i+80)])
		}
		return nil
	})
}
