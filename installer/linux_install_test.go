//go:build linux || js

package installer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallingOnLinuxPutsTheGameInTheMenu(t *testing.T) {
	p := homeOn(t, Linux)
	s := Setup{App: farmApp, Files: gameFiles(t), Place: p, Uninstaller: []byte("uninstaller")}
	if err := Install(s); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(p.Data, "com.example.farm")
	if got := read(t, filepath.Join(dir, "assets/deep/a.txt")); got != "a" {
		t.Errorf("a file holds %q", got)
	}
	if info, err := os.Stat(filepath.Join(dir, "farm")); err != nil || info.Mode().Perm()&0o100 == 0 {
		t.Errorf("the program is not executable: %v %v", info, err)
	}
	if got := read(t, filepath.Join(dir, "uninstall")); got != "uninstaller" {
		t.Errorf("the uninstaller holds %q", got)
	}
	entry := read(t, filepath.Join(p.Data, "applications", "com.example.farm.desktop"))
	if !strings.Contains(entry, `Exec="`+filepath.Join(dir, "farm")+`"`) || !strings.Contains(entry, "Name=Farm\n") {
		t.Errorf("the entry is wrong:\n%s", entry)
	}
	if to, err := os.Readlink(filepath.Join(p.Home, ".local/bin/com.example.farm")); err != nil || to != filepath.Join(dir, "farm") {
		t.Errorf("the link goes to %q, %v", to, err)
	}
	if exists(dir+".installing") || exists(dir+".old") {
		t.Error("a folder of the installing was left behind")
	}
}

// Installing again, as an update does, replaces the old files whole: a
// file the new version no longer has is gone.
func TestInstallingAgainReplacesTheOldVersion(t *testing.T) {
	p := homeOn(t, Linux)
	files := gameFiles(t)
	s := Setup{App: farmApp, Files: files, Place: p}
	if err := Install(s); err != nil {
		t.Fatal(err)
	}
	delete(files, "assets/tree.png")
	files["assets/new.png"] = files["assets/deep/a.txt"]
	if err := Install(s); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(p.Data, "com.example.farm")
	if exists(filepath.Join(dir, "assets/tree.png")) || !exists(filepath.Join(dir, "assets/new.png")) {
		t.Error("the new version did not replace the old")
	}
}

func TestUninstallingTakesAwayAllItMade(t *testing.T) {
	p := homeOn(t, Linux)
	cache := filepath.Join(p.Caches, "farm")
	s := Setup{App: farmApp, Files: gameFiles(t), Place: p, Cache: cache}
	s.Steps = []Step{{Name: "bake", Run: func(string) error {
		os.MkdirAll(cache, 0o755)
		return os.WriteFile(filepath.Join(cache, "sprite.png"), []byte("baked"), 0o644)
	}}}
	if err := Install(s); err != nil {
		t.Fatal(err)
	}
	saves := filepath.Join(p.Home, ".config", "farm")
	os.MkdirAll(saves, 0o755)
	if err := Uninstall(s); err != nil {
		t.Fatal(err)
	}
	for _, gone := range []string{filepath.Join(p.Data, "com.example.farm"), filepath.Join(p.Data, "applications", "com.example.farm.desktop"), filepath.Join(p.Home, ".local/bin/com.example.farm"), cache} {
		if exists(gone) {
			t.Errorf("%s is still there", gone)
		}
	}
	if !exists(saves) || !exists(p.Caches) {
		t.Error("uninstalling took away what was not the game's")
	}
}

// A folder of the player's own where the game would go, with no marker,
// is never replaced nor taken away.
func TestAFolderNotTheAppsIsLeftAlone(t *testing.T) {
	p := homeOn(t, Linux)
	dir := filepath.Join(p.Data, "com.example.farm")
	os.MkdirAll(dir, 0o755)
	os.WriteFile(filepath.Join(dir, "mine.txt"), []byte("mine"), 0o644)
	s := Setup{App: farmApp, Files: gameFiles(t), Place: p}
	if err := Install(s); err == nil {
		t.Error("installed over a folder not the app's")
	}
	if err := Uninstall(s); err == nil {
		t.Error("uninstalled a folder not the app's")
	}
	if read(t, filepath.Join(dir, "mine.txt")) != "mine" {
		t.Error("the player's file was touched")
	}
}

// A link of the player's own by the game's name is left as it is.
func TestALinkElsewhereIsLeftAlone(t *testing.T) {
	p := homeOn(t, Linux)
	bin := filepath.Join(p.Home, ".local/bin")
	os.MkdirAll(bin, 0o755)
	os.Symlink("/usr/bin/true", filepath.Join(bin, "com.example.farm"))
	s := Setup{App: farmApp, Files: gameFiles(t), Place: p}
	if err := Install(s); err != nil {
		t.Fatal(err)
	}
	if err := Uninstall(s); err != nil {
		t.Fatal(err)
	}
	if to, _ := os.Readlink(filepath.Join(bin, "com.example.farm")); to != "/usr/bin/true" {
		t.Errorf("the player's link now goes to %q", to)
	}
}

func TestACacheOutsideThePlayersCacheIsRefused(t *testing.T) {
	p := homeOn(t, Linux)
	for _, cache := range []string{p.Home, p.Caches, filepath.Join(p.Caches, ".."), "farm", filepath.Join(p.Home, "farm")} {
		s := Setup{App: farmApp, Files: gameFiles(t), Place: p, Cache: cache}
		if err := Install(s); err == nil {
			t.Errorf("the cache %q is taken", cache)
		}
	}
	if exists(filepath.Join(p.Data, "com.example.farm")) {
		t.Error("the game was installed before its cache was refused")
	}
}
