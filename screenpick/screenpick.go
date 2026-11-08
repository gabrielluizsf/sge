// Package screenpick is an eyedropper: it tells where the mouse is on the
// whole screen, even over other windows, whether its left button is down,
// and the colour of the one pixel under it. It never keeps or returns the
// screen, only that pixel. Where it cannot work (Windows, a Wayland
// session, no X display) Open says so and nothing else is tried.
package screenpick

import (
	"errors"
	"strings"
)

// ErrUnavailable is what Open gives where the screen cannot be read.
var ErrUnavailable = errors.New("screenpick: the screen cannot be read here")

// Reading is what the screen holds under the mouse.
type Reading struct {
	X, Y    int  // the mouse, in screen pixels
	Down    bool // the left button is held
	R, G, B uint8
}

// Source reads the screen. It belongs to one goroutine.
type Source interface {
	Read() (Reading, error)
	Close()
}

// usable says why a session cannot be read, nil if it may be: it has an X
// display and is not Wayland, where an X root window is not the screen.
// getenv is os.Getenv or a stand-in.
func usable(getenv func(string) string) error {
	if strings.EqualFold(getenv("XDG_SESSION_TYPE"), "wayland") {
		return ErrUnavailable
	}
	if getenv("DISPLAY") == "" {
		return ErrUnavailable
	}
	return nil
}
