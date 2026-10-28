package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// A temporary name's pattern may not reach outside its directory.
	_, err := os.MkdirTemp("", "../x*")
	fmt.Println(err)
	_, err = os.CreateTemp("", "a/b")
	fmt.Println(err)

	// An executable found through a relative PATH entry is refused.
	dir, _ := os.MkdirTemp("", "lp")
	defer os.RemoveAll(dir)
	// The tool leaves a mark if it ever runs.
	os.WriteFile(filepath.Join(dir, "gots-tool"), []byte("#!/bin/sh\ntouch ran\n"), 0o755)
	os.Chdir(dir)
	os.Setenv("PATH", ":")
	path, err := exec.LookPath("gots-tool")
	fmt.Println(path, err, errors.Is(err, exec.ErrDot))
	_, err = exec.LookPath("gots-missing")
	fmt.Println(err, errors.Is(err, exec.ErrNotFound))

	// Command keeps LookPath's refusal, and running it runs nothing.
	cmd := exec.Command("gots-tool")
	fmt.Println(cmd.Path, cmd.Err)
	err = cmd.Run()
	fmt.Println(err, errors.Is(err, exec.ErrDot))
	out, err := exec.Command("gots-tool").Output()
	fmt.Println(len(out), errors.Is(err, exec.ErrDot))
	_, err = os.Stat("ran")
	fmt.Println("tool ran:", err == nil)
	err = exec.Command("gots-missing").Run()
	fmt.Println(err, errors.Is(err, exec.ErrNotFound))
}
