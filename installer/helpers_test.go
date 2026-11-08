package installer

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

// farmApp is the app the tests install.
var farmApp = App{ID: "com.example.farm", Name: "Farm", Version: "1.2.0", Publisher: "Example", Exe: "farm", Icon: "icon.png", Game: true}

// gameFiles are a small game's files: its program, icon and a picture.
func gameFiles(t *testing.T) fstest.MapFS {
	return fstest.MapFS{
		"farm":              {Data: []byte("#!/bin/sh\necho farm\n"), Mode: 0o644},
		"icon.png":          {Data: testPNG(t, 300, 200)},
		"assets/tree.png":   {Data: []byte("tree")},
		"assets/deep/a.txt": {Data: []byte("a")},
	}
}

// homeOn is a place on system whose home is a new folder.
func homeOn(t *testing.T, system string) Place {
	home := t.TempDir()
	return Place{System: system, Home: home, Data: filepath.Join(home, ".local", "share"), Caches: filepath.Join(home, ".cache")}
}

// hostSetup is a Setup installing app from files on this system's own
// layout, into a new home, never running the system's programs: the
// tests that run them are each system's own.
func hostSetup(t *testing.T, app App, files fstest.MapFS) Setup {
	return Setup{App: app, Files: files, Place: testPlace(t), Run: func(string, []string, []string) error { return nil }}
}

// where is the layout s installs to.
func where(t *testing.T, s Setup) layout {
	l, err := layoutOf(s.App, s.Place, nil)
	if err != nil {
		t.Fatal(err)
	}
	return l
}

func read(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func exists(name string) bool {
	_, err := os.Lstat(name)
	return err == nil
}
