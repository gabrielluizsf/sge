package installer

import (
	"errors"
	"os"
)

// Uninstall takes away what Install made for s's app: its folder, its
// entries and links, its marked cache; never the player's saves. Every
// path is worked out again from the app, none read from disk, and the
// folder is taken away only when its marker says it is the app's.
func Uninstall(s Setup) error {
	l, err := layoutOf(s.App, s.Place, nil)
	if err != nil {
		return err
	}
	mine, there, err := ours(l.Dir, l.markerDir(), s.App)
	if err != nil {
		return err
	}
	if there && !mine {
		return errors.New("installer: " + l.Dir + " is not " + s.App.Name + "'s")
	}
	var errs []error
	for _, m := range l.Outside {
		errs = append(errs, removeIf(m.Path))
	}
	if to, err := os.Readlink(l.Link); l.Link != "" && err == nil && to == l.Exe {
		errs = append(errs, os.Remove(l.Link))
	}
	if l.system == Windows {
		errs = append(errs, removeIf(l.Shortcut), runScript(s, windowsUninstall(s.App)))
	}
	if s.Cache != "" && cacheOK(s) == nil {
		if cache, _, err := ours(s.Cache, ".", s.App); cache {
			errs = append(errs, os.RemoveAll(s.Cache))
		} else {
			errs = append(errs, err)
		}
	}
	if mine {
		errs = append(errs, os.RemoveAll(l.Dir))
	}
	return errors.Join(errs...)
}

// removeIf removes name, when it is there.
func removeIf(name string) error {
	if err := os.Remove(name); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}
