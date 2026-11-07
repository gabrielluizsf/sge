package main

import (
	"os"
	"path/filepath"
)

// cacheDir is where the baked sprites are kept: sge-realfarm in the
// user's cache, outside the repository. With no cache of the user's, it
// is a new folder of the user's alone in the temporary directory, never
// one whose name another user could guess and make first, links and
// all, for the baker to write through.
func cacheDir() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return os.MkdirTemp("", "sge-realfarm-")
	}
	return filepath.Join(dir, "sge-realfarm"), nil
}
