package gots

import (
	"embed"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// runtimeFS holds the TypeScript runtime and standard library every
// transpiled program is written against.
//
//go:embed runtime
var runtimeFS embed.FS

// hasStd reports whether there is a TypeScript version of the standard
// library package path.
func hasStd(path string) bool {
	_, err := fs.Stat(runtimeFS, "runtime/std/"+path+"/mod.ts")
	return err == nil
}

// writeRuntime copies the runtime into out/_gots.
func writeRuntime(out string) error {
	return fs.WalkDir(runtimeFS, "runtime", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := runtimeFS.ReadFile(p)
		if err != nil {
			return err
		}
		dst := filepath.Join(out, "_gots", filepath.FromSlash(p[len("runtime/"):]))
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		return os.WriteFile(dst, data, 0o644)
	})
}

var tsImport = regexp.MustCompile(`from "((?:\.\./)+)([a-z0-9/]+)/mod\.ts"`)

// stdDeps are the standard packages a hand-written package's TypeScript
// imports, which must be there too, transpiled if they are not hand-written.
func stdDeps(pkgPath string) []string {
	src, err := runtimeFS.ReadFile("runtime/std/" + pkgPath + "/mod.ts")
	if err != nil {
		return nil
	}
	var out []string
	for _, m := range tsImport.FindAllStringSubmatch(string(src), -1) {
		up := strings.Count(m[1], "../")
		dir := strings.Split(pkgPath, "/")
		if up > len(dir) {
			continue
		}
		p := path.Join(append(dir[:len(dir)-up], m[2])...)
		if !strings.HasPrefix(p, "rt") || up < len(dir)+1 {
			out = append(out, p)
		}
	}
	return out
}
