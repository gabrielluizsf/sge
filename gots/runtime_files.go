package gots

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
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
