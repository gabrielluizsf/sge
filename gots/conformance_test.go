package gots

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// jsRuntime is the command that runs a TypeScript file: deno when it is
// installed, node otherwise.
func jsRuntime(t *testing.T) []string {
	if p, err := exec.LookPath("deno"); err == nil {
		// The network the conformance programs reach: loopback, and the one
		// Unix socket unixsock listens on, which Deno grants by its path.
		sock := filepath.Join(os.TempDir(), "gots-conformance.sock")
		return []string{p, "run", "-q", "--no-config", "--no-lock", "--allow-read", "--allow-write", "--allow-env", "--allow-run", "--unstable-net", "--allow-net=127.0.0.1,unix:" + sock}
	}
	if p, err := exec.LookPath("node"); err == nil {
		return []string{p, "--no-warnings"}
	}
	t.Skip("neither deno nor node is installed")
	return nil
}

// TestConformance runs every program in testdata/conformance with go run and
// as transpiled TypeScript, and wants the same output from both.
func TestConformance(t *testing.T) {
	run := jsRuntime(t)
	dirs, err := filepath.Glob("testdata/conformance/*")
	if err != nil {
		t.Fatal(err)
	}
	for _, dir := range dirs {
		t.Run(filepath.Base(dir), func(t *testing.T) {
			t.Parallel()
			want, err := exec.Command("go", "run", "./"+dir).CombinedOutput()
			if err != nil {
				t.Fatalf("go run: %v\n%s", err, want)
			}
			out := t.TempDir()
			_, err = Transpile(Options{Dir: ".", Out: out, Patterns: []string{"./" + dir}})
			if err != nil {
				if strings.Contains(err.Error(), "no TypeScript version of standard package") {
					t.Skip(err)
				}
				t.Fatalf("transpile: %v", err)
			}
			mod := filepath.Join(out, "_gots", "mod", "github.com", "gabrielluizsf", "sge", "gots", filepath.ToSlash(dir), "mod.ts")
			if _, err := os.Stat(mod); err != nil {
				mod = filepath.Join(out, filepath.ToSlash("gots/"+dir), "mod.ts")
			}
			cmd := exec.Command(run[0], append(run[1:], mod)...)
			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			got, err := cmd.Output()
			if err != nil {
				t.Fatalf("running %s: %v\n%s\n%s", mod, err, got, stderr.String())
			}
			if !bytes.Equal(got, want) {
				t.Errorf("output differs\n%s", diffLines(string(want), string(got)))
			}
		})
	}
}

// diffLines shows the lines that differ, Go's first.
func diffLines(want, got string) string {
	w, g := strings.Split(want, "\n"), strings.Split(got, "\n")
	var b strings.Builder
	for i := 0; i < len(w) || i < len(g); i++ {
		var wl, gl string
		if i < len(w) {
			wl = w[i]
		}
		if i < len(g) {
			gl = g[i]
		}
		if wl != gl {
			b.WriteString("go: " + wl + "\nts: " + gl + "\n")
		}
	}
	return b.String()
}
