package gots

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// headlessChrome is a headless Chromium to run pages in: $GOTS_CHROME, or
// Playwright's, or "" for none.
func headlessChrome() string {
	if p := os.Getenv("GOTS_CHROME"); p != "" {
		return p
	}
	home, _ := os.UserHomeDir()
	found, _ := filepath.Glob(filepath.Join(home, ".cache/ms-playwright/chromium_headless_shell-*/chrome-headless-shell-linux64/chrome-headless-shell"))
	if len(found) == 0 {
		return ""
	}
	return found[len(found)-1]
}

// testPage drives testdata/browser/window's canvas as a player would, and
// reads back what it drew.
const testPage = `<!doctype html><html><head><meta charset="utf-8"></head><body>
<script type="module" src="app.js"></script>
<script>
setTimeout(() => {
  const c = document.querySelector("canvas[data-gots]");
  const send = (type, init) => c.dispatchEvent(new (type.startsWith("key") ? KeyboardEvent : MouseEvent)(type, { bubbles: true, ...init }));
  send("keydown", { code: "ArrowRight", key: "ArrowRight" });
  send("keyup", { code: "ArrowRight", key: "ArrowRight" });
  send("keydown", { code: "KeyA", key: "a" });
  send("keyup", { code: "KeyA", key: "a" });
  const r = c.getBoundingClientRect();
  send("mousedown", { button: 0, clientX: r.left + 10, clientY: r.top + 10 });
  send("mouseup", { button: 0, clientX: r.left + 10, clientY: r.top + 10 });
  setTimeout(() => console.log("pixel " + c.getContext("2d").getImageData(32, 24, 1, 1).data.join(",")), 300);
  console.log("favicon " + (document.querySelector("link[rel~='icon']")?.href ?? "").slice(0, 22));
  console.log("cursor " + c.style.cursor.replace(/url\("?data:[^)]*\)/, "url(png)"));
}, 1000);
</script></body></html>`

// TestWindowInBrowser runs a program that opens a window in a page of a
// headless Chromium: the window is a canvas, the page's events its input.
func TestWindowInBrowser(t *testing.T) {
	chrome := headlessChrome()
	if chrome == "" {
		t.Skip("no headless Chromium (set GOTS_CHROME)")
	}
	if _, err := exec.LookPath("deno"); err != nil {
		t.Skip("deno is not installed")
	}
	site := t.TempDir()
	// The test's own page, which writeWeb keeps, drives the program's.
	if err := os.WriteFile(filepath.Join(site, "index.html"), []byte(testPage), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Transpile(Options{Dir: ".", Out: t.TempDir(), Patterns: []string{"./testdata/browser/window"}, Web: site}); err != nil {
		t.Fatalf("transpile: %v", err)
	}
	srv := httptest.NewServer(http.FileServer(http.Dir(site)))
	defer srv.Close()
	b, _ := exec.Command(chrome, "--no-sandbox", "--enable-logging=stderr", "--v=0", "--virtual-time-budget=8000",
		"--dump-dom", srv.URL+"/index.html").CombinedOutput()
	var lines []string
	for _, m := range regexp.MustCompile(`CONSOLE[^\]]*\] "([^"]*)"`).FindAllSubmatch(b, -1) {
		lines = append(lines, string(m[1]))
	}
	got := strings.Join(lines, "\n")
	for _, want := range []string{"icon true cursor true", "favicon data:image/png;base64", "cursor url(png) 1 2, auto", "key right", "text a", "click left", "done", "pixel 229,72,77,255"} {
		if !strings.Contains(got, want) {
			t.Errorf("the page did not log %q; it logged:\n%s", want, got)
		}
	}
}
