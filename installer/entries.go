package installer

import (
	"fmt"
	"io/fs"
)

// entry is one file to copy, and its size as the files tell it.
type entry struct {
	name string
	size int64
}

// entries lists the files to copy, checked before any is written: only
// folders and regular files, under names valid on every system, none the
// installer's own, at most lim of them, in all no more bytes than lim.
func entries(files fs.FS, lim Limits, own []string) ([]entry, error) {
	var list []entry
	var total int64
	err := fs.WalkDir(files, ".", func(name string, d fs.DirEntry, err error) error {
		if err != nil || name == "." {
			return err
		}
		if !validFile(name) {
			return fmt.Errorf("installer: refused the file name %q", name)
		}
		if d.IsDir() {
			return nil
		}
		if !d.Type().IsRegular() {
			return fmt.Errorf("installer: %q is not a regular file", name)
		}
		for _, o := range own {
			if equalFold(name, o) {
				return fmt.Errorf("installer: %q is the installer's own file", name)
			}
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		total += max(info.Size(), 0)
		list = append(list, entry{name: name, size: info.Size()})
		if len(list) > lim.Files || total > lim.Bytes {
			return fmt.Errorf("installer: the files are more than %d files or %d bytes", lim.Files, lim.Bytes)
		}
		return nil
	})
	return list, err
}
