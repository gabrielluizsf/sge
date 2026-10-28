package gots

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestAudioOnDeno plays a stream through SDL2's dummy driver, which takes
// sound as fast as speakers would and plays it nowhere: the stream keeps
// it fed, and never queues more than a little.
func TestAudioOnDeno(t *testing.T) {
	deno, err := exec.LookPath("deno")
	if err != nil {
		t.Skip("deno is not installed")
	}
	out := t.TempDir()
	if _, err := Transpile(Options{Dir: ".", Out: out, Patterns: []string{"./testdata/audio/stream"}}); err != nil {
		t.Fatalf("transpile: %v", err)
	}
	cmd := exec.Command(deno, "run", "-q", "--no-config", "--no-lock", "--allow-read", "--allow-env", "--allow-ffi",
		filepath.Join(out, "gots", "testdata", "audio", "stream", "mod.ts"))
	cmd.Env = append(os.Environ(), "SDL_AUDIODRIVER=dummy")
	got, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("deno run: %v\n%s", err, got)
	}
	if strings.HasPrefix(string(got), "no audio") && strings.Contains(string(got), "SDL2 could not be loaded") {
		t.Skip("SDL2 is not installed")
	}
	if want := "rate 44100, kept up true, queue short true\n"; string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
