package installer

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// markCache marks the Setup's Cache as the app's, so uninstalling takes
// it away; with no Cache, it does nothing.
func (w *work) markCache() error {
	if w.s.Cache == "" {
		return nil
	}
	if err := os.MkdirAll(w.s.Cache, 0o755); err != nil {
		return err
	}
	root, err := os.OpenRoot(w.s.Cache)
	if err != nil {
		return err
	}
	defer root.Close()
	m := markerFile(w.s.App, "")
	return root.WriteFile(m.Path, m.Data, os.FileMode(m.Mode))
}

// cacheOK is an error unless the Setup's Cache is a folder inside the
// player's cache folder: never the cache folder itself, nor home.
func cacheOK(s Setup) error {
	c, top := filepath.Clean(s.Cache), filepath.Clean(s.Place.Caches)
	rel, err := filepath.Rel(top, c)
	if s.Place.Caches == "" || !filepath.IsAbs(c) || err != nil || rel == "." || !filepath.IsLocal(rel) || strings.HasPrefix(rel, "..") {
		return errors.New("installer: Cache must be a folder inside " + s.Place.Caches)
	}
	return nil
}
