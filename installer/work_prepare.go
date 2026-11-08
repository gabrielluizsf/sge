package installer

import (
	"errors"
	"os"
	"path"
	"path/filepath"
)

// prepare checks that the folder installed into is the app's or not
// there, and makes the new one, marked as the app's from the first.
func (w *work) prepare() error {
	l, app := w.l, w.s.App
	if mine, there, err := ours(l.Dir, l.markerDir(), app); err != nil || there && !mine {
		return errors.Join(err, errors.New("installer: "+l.Dir+" is there and is not "+app.Name+"'s"))
	}
	if err := clear(l.stage(), l, app); err != nil {
		return err
	}
	if err := clear(l.old(), l, app); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(l.Dir), 0o755); err != nil {
		return err
	}
	if err := os.Mkdir(l.stage(), 0o755); err != nil {
		return err
	}
	root, err := os.OpenRoot(l.stage())
	if err != nil {
		return err
	}
	w.root = root
	return writeIn(root, l.Inside[0])
}

// writeIn writes m new into root, its folders made.
func writeIn(root *os.Root, m made) error {
	if err := root.MkdirAll(path.Dir(m.Path), 0o755); err != nil {
		return err
	}
	f, err := root.OpenFile(m.Path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, os.FileMode(m.Mode))
	if err != nil {
		return err
	}
	_, err = f.Write(m.Data)
	return errors.Join(err, f.Close())
}
