//go:build !(linux && amd64)

package screenpick

// Open is ErrUnavailable here: only X11 on Linux is done.
func Open() (Source, error) { return nil, ErrUnavailable }
