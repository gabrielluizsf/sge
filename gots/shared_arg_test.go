package gots

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestReadOnlyArgumentsAreNotCopied guards the speed of struct arguments:
// a parameter the callee only reads gets the caller's variable as it is,
// where nothing else reaches it (testdata/conformance/argcopy checks the
// cases that must still copy).
func TestReadOnlyArgumentsAreNotCopied(t *testing.T) {
	out := t.TempDir()
	if _, err := Transpile(Options{Dir: ".", Out: out, Patterns: []string{"./testdata/conformance/argcopy"}}); err != nil {
		t.Fatalf("transpile: %v", err)
	}
	src, err := os.ReadFile(filepath.Join(out, "gots", "testdata", "conformance", "argcopy", "mod.ts"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"sum(p)", "p.Add(p)", "first(a)"} {
		if !strings.Contains(string(src), want) {
			t.Errorf("%s is copied:\n%s", want, src)
		}
	}
	for _, want := range []string{"bump(p.$clone())", "later(p.$clone())", "sumThen(w.$clone()"} {
		if !strings.Contains(string(src), want) {
			t.Errorf("%s is not copied", want)
		}
	}
}
