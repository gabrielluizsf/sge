package gots

import (
	"errors"
	"fmt"
	"html"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

// writeWeb writes each requested main package to opts.Web as a page a
// browser runs: its TypeScript bundled into app.js by deno bundle, beside
// an index.html showing the canvas its window draws to. A page already
// there is the program's own and is kept. One program goes at the top of
// opts.Web, several in a directory each, named after their packages.
func writeWeb(opts Options, prog *Program, lay layout) error {
	deno, err := exec.LookPath("deno")
	if err != nil {
		return errors.New("gots: writing a page needs deno, to bundle the program")
	}
	var mains []*Package
	for _, p := range prog.Roots {
		if p.Name == "main" {
			mains = append(mains, p)
		}
	}
	if len(mains) == 0 {
		return errors.New("gots: no main package to write a page for")
	}
	for _, p := range mains {
		dir := opts.Web
		if len(mains) > 1 {
			dir = filepath.Join(dir, path.Base(p.Path))
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		mod := filepath.Join(opts.Out, filepath.FromSlash(lay.file(p)))
		out, err := exec.Command(deno, "bundle", "-q", "--platform", "browser", "-o", filepath.Join(dir, "app.js"), mod).CombinedOutput()
		if err != nil {
			return fmt.Errorf("gots: bundling %s: %v\n%s", p.Path, err, out)
		}
		page := filepath.Join(dir, "index.html")
		if _, err := os.Stat(page); err == nil {
			continue
		}
		if err := os.WriteFile(page, []byte(fmt.Sprintf(webPage, html.EscapeString(path.Base(p.Path)))), 0o644); err != nil {
			return err
		}
	}
	return nil
}

// webPage centres the canvas on a dark page, scaled down to fit it and
// kept sharp; the program's window sets its title.
const webPage = `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>%s</title>
<style>
  html, body { margin: 0; height: 100%%; background: #111; }
  body { display: flex; align-items: center; justify-content: center; }
  canvas { max-width: 100vw; max-height: 100vh; image-rendering: pixelated; outline: none; }
</style>
</head>
<body>
<canvas data-gots></canvas>
<script type="module" src="app.js"></script>
</body>
</html>
`
