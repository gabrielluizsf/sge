package gots

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestTypeScriptBoundary runs testdata/tsapi against the engine transpiled
// with an index, as a TypeScript program uses it.
func TestTypeScriptBoundary(t *testing.T) {
	t.Parallel()
	deno, err := exec.LookPath("deno")
	if err != nil {
		t.Skip("deno is not installed")
	}
	out := engine(t)
	src, err := os.ReadFile("testdata/tsapi/boundary_test.ts")
	if err != nil {
		t.Fatal(err)
	}
	test := filepath.Join(out, "boundary_test.ts") // beside the index, which it imports
	if err := os.WriteFile(test, src, 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(deno, "test", "-q", "--no-config", "--no-lock", test)
	if got, err := cmd.CombinedOutput(); err != nil {
		t.Errorf("deno test: %v\n%s", err, got)
	}
}
