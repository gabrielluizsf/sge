package installer

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// runSystem runs program for real, with no window of its own; an error
// carries what the program said.
func runSystem(program string, args, env []string) error {
	cmd := exec.Command(program, args...)
	cmd.Env = append(os.Environ(), env...)
	hideWindow(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil && len(out) > 0 {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out[:min(len(out), 2048)])))
	}
	return err
}
