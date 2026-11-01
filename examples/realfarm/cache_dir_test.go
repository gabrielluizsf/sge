package main

import (
	"os"
	"path/filepath"
	"testing"
)

// With no cache of the user's, the sprites go to a new folder only the
// user may open, not to a shared one anyone could have made first.
func TestCacheWithoutAUserCacheIsPrivateAndNew(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", "")
	t.Setenv("HOME", "")
	shared := filepath.Join(os.TempDir(), "sge-realfarm")
	a, err := cacheDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(a)
	b, err := cacheDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(b)
	if a == shared || a == b {
		t.Fatalf("folders %q and %q: a name another user can guess", a, b)
	}
	if fi, err := os.Stat(a); err != nil || fi.Mode().Perm() != 0o700 {
		t.Errorf("%q is %v %v, not the user's alone", a, fi.Mode(), err)
	}
}

// With a cache of the user's, the sprites are kept in it.
func TestCacheIsTheUsersOwn(t *testing.T) {
	home := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", home)
	if dir, err := cacheDir(); err != nil || dir != filepath.Join(home, "sge-realfarm") {
		t.Errorf("%q %v", dir, err)
	}
}
