//go:build !(linux && amd64)

package native

import "errors"

// ErrUnsupported is what Open says where native cannot call C yet.
var ErrUnsupported = errors.New("native: calling C is not supported on this platform yet")

// Proc is the address of a C function.
type Proc uintptr

// Lib is a library opened with dlopen.
type Lib struct{}

// Open is ErrUnsupported here.
func Open(soname string) (*Lib, error) { return nil, ErrUnsupported }

// Symbol is ErrUnsupported here.
func (l *Lib) Symbol(name string) (Proc, error) { return 0, ErrUnsupported }

func call(f *Frame) { panic(ErrUnsupported) }
