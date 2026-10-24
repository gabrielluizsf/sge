package gots

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestMathIntrinsicsAreJavaScripts guards the speed of math: it is
// transpiled from Go's source, whose Sqrt is a loop in software and whose
// Sin reduces in BigInts, so calls of these are JavaScript's Math.
func TestMathIntrinsicsAreJavaScripts(t *testing.T) {
	out := t.TempDir()
	if _, err := Transpile(Options{Dir: ".", Out: out, Patterns: []string{"./testdata/conformance/mathfuncs"}}); err != nil {
		t.Fatalf("transpile: %v", err)
	}
	src, err := os.ReadFile(filepath.Join(out, "gots", "testdata", "conformance", "mathfuncs", "mod.ts"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"$Math.sqrt(2)", "$Math.floor(", "$Math.trunc(", "$Math.atan2("} {
		if !strings.Contains(string(src), want) {
			t.Errorf("no %s in\n%s", want, src)
		}
	}
}
