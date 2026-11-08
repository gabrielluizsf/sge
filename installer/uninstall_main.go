package installer

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// uninstallMain uninstalls s's app. On Windows a program cannot take
// its own file away while it runs, so an uninstaller in the app's folder
// first copies itself to the player's temporary folder and leaves that
// copy to do it, once it has gone.
func uninstallMain(s Setup) error {
	var err error
	if s.Place, err = Here(); err != nil {
		return err
	}
	l, err := layoutOf(s.App, s.Place, nil)
	if err != nil {
		return err
	}
	self, err := os.Executable()
	if err != nil {
		return err
	}
	if s.Place.System == Windows && inside(self, l.Dir) {
		return handOver(self)
	}
	for try := 0; ; try++ { // the first uninstaller may still be closing
		if err = Uninstall(s); err == nil || try == 20 {
			return err
		}
		time.Sleep(500 * time.Millisecond)
	}
}

// inside is whether file is in dir, names compared as Windows does.
func inside(file, dir string) bool {
	rel, err := filepath.Rel(strings.ToLower(dir), strings.ToLower(file))
	return err == nil && filepath.IsLocal(rel)
}

// handOver starts a copy of the program self, made in a new file of the
// player's temporary folder, to uninstall.
func handOver(self string) error {
	data, err := os.ReadFile(self)
	if err != nil {
		return err
	}
	f, err := os.CreateTemp("", "uninstall-*.exe")
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		return err
	}
	cmd := exec.Command(f.Name(), "--uninstall")
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Process.Release()
}
