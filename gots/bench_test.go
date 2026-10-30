package gots

import (
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// TestBenchmarksOnDeno runs testdata/bench's Go benchmarks under deno
// bench, as go test -bench runs them.
func TestBenchmarksOnDeno(t *testing.T) {
	t.Parallel()
	deno, err := exec.LookPath("deno")
	if err != nil {
		t.Skip("deno is not installed")
	}
	out := t.TempDir()
	if _, err := Transpile(Options{Dir: ".", Out: out, Patterns: []string{"./testdata/bench"}, Tests: true}); err != nil {
		t.Fatalf("transpile: %v", err)
	}
	mods, _ := filepath.Glob(filepath.Join(out, "*", "testdata", "bench", "mod_test.ts"))
	if len(mods) != 1 {
		mods, _ = filepath.Glob(filepath.Join(out, "testdata", "bench", "mod_test.ts"))
	}
	if len(mods) != 1 {
		t.Fatalf("no mod_test.ts under %s", out)
	}
	got, err := exec.Command(deno, "bench", "--no-check", "--no-config", "--no-lock", mods[0]).CombinedOutput()
	if err != nil {
		t.Fatalf("deno bench: %v\n%s", err, got)
	}
	// Each prints its line as go test -bench does, b.N grown until the
	// loop takes long enough to time.
	for _, name := range []string{"BenchmarkSumLoop", "BenchmarkSumN", "BenchmarkChanPingPong", "BenchmarkChanBuffered", "BenchmarkWaitGroup"} {
		m := regexp.MustCompile(`(?m)^` + name + `\s+(\d+)\s+[\d.]+ ns/op$`).FindSubmatch(got)
		if m == nil {
			t.Errorf("deno bench did not report %s in Go's format:\n%s", name, got)
			continue
		}
		if n, _ := strconv.Atoi(string(m[1])); n < 1000 {
			t.Errorf("%s ran %d times; b.N was not grown", name, n)
		}
	}
	// A benchmark is not a test: deno test must still pass, and run only TestSum.
	got, err = exec.Command(deno, "test", "--no-check", "--no-config", "--no-lock", mods[0]).CombinedOutput()
	if err != nil || !strings.Contains(string(got), "1 passed") {
		t.Errorf("deno test: %v\n%s", err, got)
	}
}
