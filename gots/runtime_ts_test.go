package gots

import (
	"os/exec"
	"testing"
)

// TestRuntimeOnDeno runs the runtime's own TypeScript tests, in
// testdata/rt, which exercise what is only reachable from JavaScript.
func TestRuntimeOnDeno(t *testing.T) {
	deno, err := exec.LookPath("deno")
	if err != nil {
		t.Skip("deno is not installed")
	}
	cmd := exec.Command(deno, "test", "-q", "--no-config", "--no-lock", "--allow-net=127.0.0.1", "testdata/rt")
	if got, err := cmd.CombinedOutput(); err != nil {
		t.Errorf("deno test: %v\n%s", err, got)
	}
}
