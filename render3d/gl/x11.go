//go:build linux && amd64

package gl

import (
	"errors"

	"github.com/gabrielluizsf/sge/native"
)

// openX11 is a connection of the device's own to the X display named
// name, as DISPLAY spells it: an X Display*.
func openX11(name string) (uintptr, error) {
	lib, err := native.Open("libX11.so.6")
	if err != nil {
		return 0, err
	}
	open, err := lib.Symbol("XOpenDisplay")
	if err != nil {
		return 0, err
	}
	dpy := call(open, ptr(cstring(name)))
	if dpy == 0 {
		return 0, errors.New("gl: XOpenDisplay " + name + " failed")
	}
	return dpy, nil
}
