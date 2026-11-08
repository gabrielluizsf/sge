package installer

import (
	"fmt"
	"os"
	"slices"
)

// Main is an installer program's main, s its app (and its Steps and
// Cache), look its window: with no arguments it installs the files it
// carries in a window, then starts the game; with --quiet it shows no
// window and starts nothing; with --uninstall it takes the game away.
// It is the exit code.
func Main(s Setup, look Window) int {
	args := os.Args[1:]
	var err error
	if slices.Contains(args, "--uninstall") {
		err = uninstallMain(s)
	} else {
		err = installMain(s, look, slices.Contains(args, "--quiet"))
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

// installMain installs the files the running program carries.
func installMain(s Setup, look Window, quiet bool) error {
	self, err := os.Executable()
	if err != nil {
		return err
	}
	p, err := Open(self)
	if err != nil {
		return err
	}
	defer p.Close()
	if s.Place, err = Here(); err != nil {
		return err
	}
	s.Files, s.Uninstaller = p.Files, p.Program
	if quiet {
		return Install(s)
	}
	if err := installWindow(s, look); err != nil {
		return err
	}
	return launch(s)
}
