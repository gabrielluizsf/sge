package gots

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestAudioInBrowser plays testdata/audio/stream with Web Audio in a
// headless Chromium let play sound unasked. The page's time runs ahead of
// the sound's there, so what is checked is that the stream opens and,
// though it is asked for more than is played, never queues more than a
// little; TestAudioOnDeno checks it keeps up.
func TestAudioInBrowser(t *testing.T) {
	chrome := headlessChrome()
	if chrome == "" {
		t.Skip("no headless Chromium (set GOTS_CHROME)")
	}
	site := t.TempDir()
	if _, err := Transpile(Options{Dir: ".", Out: t.TempDir(), Patterns: []string{"./testdata/audio/stream"}, Web: site}); err != nil {
		t.Fatalf("transpile: %v", err)
	}
	if _, err := os.Stat(filepath.Join(site, "index.html")); err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(http.FileServer(http.Dir(site)))
	defer srv.Close()
	b, _ := exec.Command(chrome, "--no-sandbox", "--enable-logging=stderr", "--v=0", "--autoplay-policy=no-user-gesture-required",
		"--virtual-time-budget=5000", "--dump-dom", srv.URL+"/index.html").CombinedOutput()
	if want := "queue short true"; !strings.Contains(string(b), "rate 44100, kept up") || !strings.Contains(string(b), want) {
		t.Errorf("the page did not log %q; it logged:\n%s", want, b)
	}
}
