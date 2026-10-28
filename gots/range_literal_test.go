package gots

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRangeOverLiteralAllocatesNoTypedArray guards the speed of ranging
// over an array literal: a typed array costs V8 about a microsecond to
// make, ten times a plain one, and the literal is only read by the loop.
func TestRangeOverLiteralAllocatesNoTypedArray(t *testing.T) {
	out := t.TempDir()
	if _, err := Transpile(Options{Dir: ".", Out: out, Patterns: []string{"./testdata/conformance/rangelit"}}); err != nil {
		t.Fatalf("transpile: %v", err)
	}
	var mods []string
	filepath.WalkDir(out, func(p string, d os.DirEntry, err error) error {
		if err == nil && strings.HasSuffix(p, filepath.Join("rangelit", "mod.ts")) {
			mods = append(mods, p)
		}
		return nil
	})
	if len(mods) != 1 {
		t.Fatalf("found %v", mods)
	}
	src, err := os.ReadFile(mods[0])
	if err != nil {
		t.Fatal(err)
	}
	for _, b := range []string{"Float64Array.of", "Float32Array.of", "Int8Array.of", "Uint8Array.of", "BigUint64Array.of"} {
		if strings.Contains(string(src), b) {
			t.Errorf("a ranged literal is a %s:\n%s", b, src)
		}
	}
}
