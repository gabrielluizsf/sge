package installer

import (
	"errors"
	"os"
	"path/filepath"
)

// register makes the app known to the system: its entry, its link and,
// on Windows, its shortcut and its place among the installed apps.
func (w *work) register() error {
	l := w.l
	for _, m := range l.Outside {
		if err := writeAtomic(m); err != nil {
			return err
		}
	}
	if l.Link != "" {
		if err := link(l); err != nil {
			return err
		}
	}
	if l.system == Windows {
		return runScript(w.s, windowsInstall(w.s.App, l, size(w)))
	}
	return nil
}

// size is about how much the installed files hold, in KiB.
func size(w *work) int64 {
	return (w.s.Limits.orDefault().Bytes - w.left) / 1024
}

// writeAtomic writes m's file whole or not at all: into a new file
// beside it, then moved over it (a link there is replaced, not followed).
func writeAtomic(m made) error {
	dir := filepath.Dir(m.Path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	f, err := os.CreateTemp(dir, ".sge-installer-*")
	if err != nil {
		return err
	}
	_, err = f.Write(m.Data)
	err = errors.Join(err, f.Chmod(os.FileMode(m.Mode)), f.Close())
	if err == nil {
		err = os.Rename(f.Name(), m.Path)
	}
	if err != nil {
		os.Remove(f.Name())
	}
	return err
}

// link points l.Link at the program, unless something not the app's
// holds that name: a file, or a link to elsewhere, is left as it is.
func link(l layout) error {
	if err := os.MkdirAll(filepath.Dir(l.Link), 0o755); err != nil {
		return err
	}
	if to, err := os.Readlink(l.Link); err == nil {
		if to != l.Exe {
			return nil
		}
		if err := os.Remove(l.Link); err != nil {
			return err
		}
	} else if _, err := os.Lstat(l.Link); err == nil {
		return nil
	}
	return os.Symlink(l.Exe, l.Link)
}
