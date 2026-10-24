package gots

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestFarmDemoOnDeno plays the farm example's demo in Go and, transpiled,
// on Deno, and wants the same farm at the end and the same last frame,
// pixel for pixel: tile maps, endless terrain, SVG sprites and all.
func TestFarmDemoOnDeno(t *testing.T) {
	deno, err := exec.LookPath("deno")
	if err != nil {
		t.Skip("deno is not installed")
	}
	if testing.Short() {
		t.Skip("plays a minute of the game twice")
	}
	out := t.TempDir()
	if _, err := Transpile(Options{Dir: ".", Out: out, Patterns: []string{"../examples/farm"}}); err != nil {
		t.Fatalf("transpile: %v", err)
	}
	env := func(shot string) []string {
		return append(os.Environ(), "FARM_DEMO=3600", "FARM_VIEW=1280x720", "FARM_SHOT="+shot)
	}
	goShot, denoShot := filepath.Join(out, "go.png"), filepath.Join(out, "deno.png")
	goCmd := exec.Command("go", "run", ".")
	goCmd.Dir, goCmd.Env = "../examples/farm", env(goShot)
	want, err := goCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run: %v\n%s", err, want)
	}
	denoCmd := exec.Command(deno, "run", "-q", "--no-config", "--no-lock", "--allow-read", "--allow-env", "--allow-write="+out,
		filepath.Join(out, "examples", "farm", "mod.ts"))
	denoCmd.Env = env(denoShot)
	got, err := denoCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("deno run: %v\n%s", err, got)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("Deno's farm:\n%s\nGo's:\n%s", got, want)
	}
	a, _ := os.ReadFile(goShot)
	b, _ := os.ReadFile(denoShot)
	if len(a) == 0 || !bytes.Equal(a, b) {
		t.Error("the last frames differ")
	}
}
