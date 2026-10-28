//go:build gots_overlay

// Overlay for sys_unix.go: there are no zoneinfo files to read, so named
// time zones other than UTC and Local fail to load.

package time

import "errors"

var errNoFiles = errors.New("time: no zoneinfo files under gots")

func interrupt() {}

func open(name string) (uintptr, error) { return 0, errNoFiles }

func read(fd uintptr, buf []byte) (int, error) { return 0, errNoFiles }

func closefd(fd uintptr) {}

func preadn(fd uintptr, buf []byte, off int) error { return errNoFiles }
