package installer

import (
	"errors"
	"os"
)

// finish writes the system's own files into the new folder and puts it
// in the old one's place: the old is moved aside first and taken away
// only once the new is in.
func (w *work) finish() error {
	l := w.l
	for _, m := range l.Inside[1:] {
		if err := writeIn(w.root, m); err != nil {
			return err
		}
	}
	if l.Uninstaller != "" && w.s.Uninstaller != nil {
		if err := writeIn(w.root, made{Path: l.Uninstaller, Data: w.s.Uninstaller, Mode: 0o755}); err != nil {
			return err
		}
	}
	err := w.root.Close()
	w.root = nil
	if err != nil {
		return err
	}
	_, there, err := ours(l.Dir, l.markerDir(), w.s.App)
	if err != nil {
		return err
	}
	if there {
		if err := os.Rename(l.Dir, l.old()); err != nil {
			return errors.Join(errors.New("installer: is "+w.s.App.Name+" still open?"), err)
		}
	}
	if err := os.Rename(l.stage(), l.Dir); err != nil {
		return err
	}
	return clear(l.old(), l, w.s.App)
}
