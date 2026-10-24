package gots

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestTuplesAreNotDestructured guards the speed of multiple results and
// parallel assignment: V8 takes a destructuring through the iterator
// protocol, about a tenth slower than reading the elements in a hot loop
// such as antui's BlitMatrix, and [a, b] = [b, a] also makes an array.
func TestTuplesAreNotDestructured(t *testing.T) {
	out := t.TempDir()
	if _, err := Transpile(Options{Dir: ".", Out: out, Patterns: []string{"./testdata/conformance/tuples"}}); err != nil {
		t.Fatalf("transpile: %v", err)
	}
	src, err := os.ReadFile(filepath.Join(out, "gots", "testdata", "conformance", "tuples", "mod.ts"))
	if err != nil {
		t.Fatal(err)
	}
	for _, bad := range []string{"let [", "] = ["} {
		if strings.Contains(string(src), bad) {
			t.Errorf("found %q:\n%s", bad, src)
		}
	}
}
