package gots

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestInterfaceValuesAreNotCopiedNeedlessly guards the speed of values in
// interfaces: boxes are built with their class named, which V8 can inline,
// and a type switch case that only reads its variable is handed the boxed
// value itself (testdata/conformance/ifacecopy checks the cases that copy).
func TestInterfaceValuesAreNotCopiedNeedlessly(t *testing.T) {
	out := t.TempDir()
	if _, err := Transpile(Options{Dir: ".", Out: out, Patterns: []string{"./testdata/conformance/ifacecopy"}}); err != nil {
		t.Fatalf("transpile: %v", err)
	}
	src, err := os.ReadFile(filepath.Join(out, "gots", "testdata", "conformance", "ifacecopy", "mod.ts"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(src)
	for _, want := range []string{"new V$Val(V.$type, new V(8, 9))", "new $.Val(plain.$type, new plain(3))", "let v: V = $x.v;"} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %s:\n%s", want, s)
		}
	}
	if !strings.Contains(s, "let v$1: plain = $x.v.$clone();") {
		t.Errorf("a case that changes its variable is not given a copy:\n%s", s)
	}
}
