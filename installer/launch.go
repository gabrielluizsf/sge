package installer

import (
	"os"
	"os/exec"
)

// launch starts the installed game, in its own folder, and lets it run
// on its own after the installer has gone.
func launch(s Setup) error {
	l, err := layoutOf(s.App, s.Place, nil)
	if err != nil {
		return err
	}
	cmd := exec.Command(l.Exe)
	cmd.Dir = l.join(l.Dir, l.Files)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = nil, nil, nil
	cmd.Env = os.Environ()
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Process.Release()
}
