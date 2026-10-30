package gots

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// enginePackages are sge's own packages, whose Go tests must pass on Deno
// too once transpiled.
var enginePackages = []string{"vector", "transform", "collider", "clock", "sprite", "animation", "game", "tilemap", "terrain", "iso", "random", "inventory", "request", "pathfind", "town", "scene", "daytime", "label"}

// TestEngineOnDeno transpiles sge's packages with their tests and runs
// them under deno test, from each package's directory as go test would.
func TestEngineOnDeno(t *testing.T) {
	deno, err := exec.LookPath("deno")
	if err != nil {
		t.Skip("deno is not installed")
	}
	var patterns []string
	for _, p := range enginePackages {
		patterns = append(patterns, "../"+p)
	}
	out := t.TempDir()
	if _, err := Transpile(Options{Dir: ".", Out: out, Patterns: patterns, Tests: true}); err != nil {
		t.Fatalf("transpile: %v", err)
	}
	for _, p := range enginePackages {
		t.Run(p, func(t *testing.T) {
			cmd := exec.Command(deno, "test", "-q", "--no-config", "--no-lock", "--allow-read", "--allow-env", filepath.Join(out, p, "mod_test.ts"))
			cmd.Dir = filepath.Join("..", p)
			got, err := cmd.CombinedOutput()
			if err != nil || !strings.Contains(string(got), " passed") {
				t.Errorf("deno test: %v\n%s", err, got)
			}
		})
	}
}
