package installer

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
)

// copyFile copies e from files into dir in root, made new (never through
// what is there), executable when it is the program, and fails when it
// holds more bytes than its size said, or than *left.
func copyFile(files fs.FS, to placed, e entry) error {
	src, err := files.Open(e.name)
	if err != nil {
		return err
	}
	defer src.Close()
	name := path.Join(to.dir, e.name)
	if err := to.root.MkdirAll(path.Dir(name), 0o755); err != nil {
		return err
	}
	perm := os.FileMode(0o644)
	if e.name == to.exe {
		perm = 0o755
	}
	dst, err := to.root.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, perm)
	if err != nil {
		return err
	}
	n, err := io.Copy(dst, io.LimitReader(src, min(e.size, *to.left)+1))
	if cerr := dst.Close(); err == nil {
		err = cerr
	}
	*to.left -= n
	if err == nil && (n > e.size || *to.left < 0) {
		err = errors.New("installer: " + e.name + " holds more than it said")
	}
	return err
}

// placed is where files are copied to: dir in root, the program exe
// among them, and how many bytes may still be written, left.
type placed struct {
	root *os.Root
	dir  string
	exe  string
	left *int64
}
