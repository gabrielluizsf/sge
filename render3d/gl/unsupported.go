//go:build !(linux && amd64)

package gl

import (
	"errors"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/sge/render3d"
)

// ErrUnsupported is what Offscreen and Window say where this backend
// cannot run yet.
var ErrUnsupported = errors.New("gl: OpenGL through EGL runs only on Linux and amd64 yet")

// Device is an OpenGL device, which cannot be made here.
type Device struct{ render3d.Device }

// Offscreen is ErrUnsupported here.
func Offscreen(w, h int) (*Device, error) { return nil, ErrUnsupported }

// Window is ErrUnsupported here.
func Window(win *antui.Window) (*Device, error) { return nil, ErrUnsupported }

// Present does nothing here.
func (d *Device) Present() {}

// Close does nothing here.
func (d *Device) Close() {}
