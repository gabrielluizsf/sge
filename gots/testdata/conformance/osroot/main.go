// Command osroot opens files through an os.Root: names inside it open,
// and those that climb out with "..", are absolute, or follow a link out
// of it are refused, with Go's error, on any platform the same.
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dir, err := os.MkdirTemp("", "osroot")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)
	inside := filepath.Join(dir, "in")
	os.MkdirAll(filepath.Join(inside, "sub"), 0o755)
	os.WriteFile(filepath.Join(inside, "a.txt"), []byte("inside"), 0o644)
	os.WriteFile(filepath.Join(inside, "sub", "b.txt"), []byte("deeper"), 0o644)
	os.WriteFile(filepath.Join(dir, "secret.txt"), []byte("secret"), 0o644)
	os.Symlink(filepath.Join(dir, "secret.txt"), filepath.Join(inside, "out.txt"))
	os.Symlink(filepath.Join(inside, "a.txt"), filepath.Join(inside, "alias.txt"))
	os.Symlink(dir, filepath.Join(inside, "up"))

	root, err := os.OpenRoot(inside)
	if err != nil {
		panic(err)
	}
	for _, name := range []string{"a.txt", "sub/b.txt", "sub/../a.txt", "./a.txt", "alias.txt",
		"../secret.txt", "sub/../../secret.txt", filepath.Join(dir, "secret.txt"), "out.txt", "up/secret.txt", "missing.txt"} {
		f, err := root.Open(name)
		if err != nil {
			fmt.Println(show(name, dir), "refused:", reason(err))
			continue
		}
		b, _ := io.ReadAll(f)
		f.Close()
		fmt.Println(show(name, dir), "read:", string(b))
	}
	sub, err := root.OpenRoot("sub")
	if err != nil {
		panic(err)
	}
	_, err = sub.Open("../a.txt")
	fmt.Println("sub ../a.txt refused:", err != nil)
	if b, err := root.ReadFile("sub/b.txt"); err == nil {
		fmt.Println("ReadFile:", string(b))
	}
	f, err := root.Create("new.txt")
	fmt.Println("create inside:", err == nil)
	if f != nil {
		f.Close()
	}
	_, err = root.Create("up/new.txt")
	fmt.Println("create through a link out refused:", err != nil)
	fmt.Println("mkdirall inside:", root.MkdirAll("made/deeper", 0o755) == nil)
	fmt.Println("write inside:", root.WriteFile("made/deeper/c.txt", []byte("written"), 0o644) == nil)
	fmt.Println("rename inside:", root.Rename("made/deeper/c.txt", "made/c.txt") == nil)
	if b, err := root.ReadFile("made/c.txt"); err == nil {
		fmt.Println("renamed reads:", string(b))
	}
	fmt.Println("mkdirall out refused:", root.MkdirAll("../outside/x", 0o755) != nil)
	fmt.Println("write through a link out refused:", root.WriteFile("up/planted.txt", []byte("x"), 0o644) != nil)
	fmt.Println("rename out refused:", root.Rename("made/c.txt", "../stolen.txt") != nil)
	_, err = os.Stat(filepath.Join(dir, "planted.txt"))
	fmt.Println("nothing planted outside:", err != nil)
	fmt.Println("remove inside:", root.Remove("made/c.txt") == nil)
	fmt.Println("remove a link:", root.Remove("alias.txt") == nil)
	_, err = os.Stat(filepath.Join(inside, "a.txt"))
	fmt.Println("its target kept:", err == nil)
	fmt.Println("remove the root refused:", root.Remove(".") != nil)
	os.WriteFile(filepath.Join(inside, "abs.txt"), []byte("x"), 0o644)
	fmt.Println("remove an absolute name refused:", root.Remove("/abs.txt") != nil)
	_, err = os.Stat(filepath.Join(inside, "abs.txt"))
	fmt.Println("its namesake inside kept:", err == nil)
	root.Close()
	_, err = root.Open("a.txt")
	fmt.Println("closed refused:", err != nil)
}

// show is name with the temporary directory's path hidden.
func show(name, dir string) string { return strings.Replace(name, dir, "<dir>", 1) }

// reason is why an open failed, in words that do not name the paths.
func reason(err error) string {
	switch {
	case errors.Is(err, os.ErrNotExist):
		return "does not exist"
	case strings.Contains(err.Error(), "escapes"):
		return "escapes"
	}
	return "other: " + err.Error()
}
