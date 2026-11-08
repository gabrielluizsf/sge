package installer

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path"
)

// ours is whether dir is a real folder (not a link) whose marker, at
// markerDir in it, says it is app's; and whether it is there at all.
func ours(dir, markerDir string, app App) (mine, there bool, err error) {
	info, err := os.Lstat(dir)
	if errors.Is(err, os.ErrNotExist) {
		return false, false, nil
	}
	if err != nil {
		return false, true, err
	}
	if !info.IsDir() {
		return false, true, nil
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		return false, true, err
	}
	defer root.Close()
	f, err := root.Open(path.Join(markerDir, marker))
	if err != nil {
		return false, true, nil
	}
	defer f.Close()
	got, err := io.ReadAll(io.LimitReader(f, 128))
	return err == nil && bytes.Equal(got, markerFile(app, "").Data), true, err
}

// markerDir is the folder in Dir that holds the marker.
func (l layout) markerDir() string {
	return path.Dir(l.Inside[0].Path)
}

// clear takes dir away when it is the app's or not there, and is an
// error when something else is there: never the player's own folder.
func clear(dir string, l layout, app App) error {
	mine, there, err := ours(dir, l.markerDir(), app)
	if err != nil || !there {
		return err
	}
	if !mine {
		return errors.New("installer: " + dir + " is there and is not " + app.Name + "'s")
	}
	return os.RemoveAll(dir)
}
