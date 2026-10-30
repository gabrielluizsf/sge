package gots

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestTLSConfigNeverFetches runs testdata/nosockets/tlsconfig on Deno made
// to look like a browser, with no Deno.connect and a fetch that tells it
// was called: a TLSClientConfig must not be left aside by fetch.
func TestTLSConfigNeverFetches(t *testing.T) {
	t.Parallel()
	deno, err := exec.LookPath("deno")
	if err != nil {
		t.Skip("deno is not installed")
	}
	out := programs(t)
	mod := filepath.Join(out, "gots", "testdata", "nosockets", "tlsconfig", "mod.ts")
	shim := filepath.Join(out, "shim.ts")
	src := `globalThis.fetch = (u) => Promise.reject(new Error("fetched " + u));
Object.defineProperty(Deno, "connect", { value: undefined });
await import(` + "`" + `file://${"` + filepath.ToSlash(mod) + `"}` + "`" + `);
`
	if err := os.WriteFile(shim, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := exec.Command(deno, "run", "-q", "--no-config", "--no-lock", "--allow-read", "--allow-env", shim).CombinedOutput()
	if err != nil {
		t.Fatalf("deno: %v\n%s", err, got)
	}
	if !strings.Contains(string(got), "pinned fetched: false") || !strings.Contains(string(got), "plain: Get \"https://example.test/\": net/http: fetch() failed: Error: fetched https://example.test/") {
		t.Errorf("got\n%s", got)
	}
}
