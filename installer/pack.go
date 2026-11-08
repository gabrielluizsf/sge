package installer

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// payloadTag starts the comment of the zip an installer carries, which
// tells how many bytes of program come before it.
const payloadTag = "sge-installer program="

// Pack writes out, an installer: the program at stub (a build of the
// game's installer, for the system it installs on), and after it the
// files in dir, as a zip. Links and other non-regular files are refused.
func Pack(stub, dir, out string) error {
	program, err := os.ReadFile(stub)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	err = packInto(f, program, os.DirFS(dir))
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		os.Remove(out)
	}
	return err
}

func packInto(w io.Writer, program []byte, files fs.FS) error {
	if _, err := w.Write(program); err != nil {
		return err
	}
	z := zip.NewWriter(w)
	z.SetOffset(int64(len(program)))
	err := fs.WalkDir(files, ".", func(name string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if !d.Type().IsRegular() || !validFile(name) {
			return fmt.Errorf("installer: cannot pack %q", filepath.FromSlash(name))
		}
		return packFile(z, files, name)
	})
	if err != nil {
		return err
	}
	return errors.Join(z.SetComment(fmt.Sprint(payloadTag, len(program))), z.Close())
}

func packFile(z *zip.Writer, files fs.FS, name string) error {
	src, err := files.Open(name)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := z.CreateHeader(&zip.FileHeader{Name: name, Method: zip.Deflate})
	if err != nil {
		return err
	}
	_, err = io.Copy(dst, src)
	return err
}
