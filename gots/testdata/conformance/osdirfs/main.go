// Command osdirfs reads a tree of files through os.DirFS: names inside it
// open, a walk sees every file, and a name fs.ValidPath refuses, or one
// not there, gives Go's error, named by its name in the tree.
package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	dir, err := os.MkdirTemp("", "osdirfs")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)
	os.MkdirAll(filepath.Join(dir, "sub", "deep"), 0o755)
	os.WriteFile(filepath.Join(dir, "a.txt"), []byte("top"), 0o644)
	os.WriteFile(filepath.Join(dir, "sub", "deep", "b.txt"), []byte("deeper"), 0o644)
	fsys := os.DirFS(dir + "/")

	fs.WalkDir(fsys, ".", func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}
		fmt.Println("walk:", name, d.IsDir())
		return nil
	})
	for _, name := range []string{"a.txt", "sub/deep/b.txt", "../a.txt", "/a.txt", "sub/../a.txt", "missing.txt"} {
		b, err := fs.ReadFile(fsys, name)
		var pe *fs.PathError
		switch {
		case errors.As(err, &pe):
			fmt.Println(name, "refused:", pe.Op, pe.Path, errors.Is(err, fs.ErrInvalid), errors.Is(err, fs.ErrNotExist))
		case err != nil:
			fmt.Println(name, "error:", err)
		default:
			fmt.Println(name, "read:", string(b))
		}
	}
	info, err := fs.Stat(fsys, "sub")
	fmt.Println("stat sub:", info.Name(), info.IsDir(), err)
}
