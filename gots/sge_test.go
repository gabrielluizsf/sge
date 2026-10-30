package gots

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// enginePackages are sge's own packages, whose Go tests must pass on Deno
// too once transpiled.
var enginePackages = []string{"vector", "transform", "collider", "clock", "sprite", "animation", "game", "tilemap", "terrain", "iso", "random", "inventory", "request", "pathfind", "town", "scene", "daytime", "label", "gltf", "bake", "loader", "fetch", "flora", "human", "life", "bond", "lineage", "pose", "render3d", "render3d/soft"}

// TestEngineOnDeno transpiles sge's packages with their tests and runs
// them under deno test, from each package's directory as go test would.
func TestEngineOnDeno(t *testing.T) {
	t.Parallel()
	deno, err := exec.LookPath("deno")
	if err != nil {
		t.Skip("deno is not installed")
	}
	out := engine(t)
	for _, p := range enginePackages {
		t.Run(p, func(t *testing.T) {
			t.Parallel()
			// The Go type checker has checked the module, which says so
			// with @ts-nocheck: TypeScript's would only take its time.
			cmd := exec.Command(deno, "test", "-q", "--no-check", "--no-config", "--no-lock", "--allow-read", "--allow-env", "--allow-write="+os.TempDir(), "--allow-net=127.0.0.1,[::1]", filepath.Join(out, p, "mod_test.ts"))
			cmd.Dir = filepath.Join("..", p)
			got, err := cmd.CombinedOutput()
			if err != nil || !strings.Contains(string(got), " passed") {
				t.Errorf("deno test: %v\n%s", err, got)
			}
		})
	}
}
