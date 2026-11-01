//go:build linux && amd64

package gl

import (
	"unsafe"

	"github.com/gabrielluizsf/sge/native"
)

// arg is an argument to a C function: an integer, or Go memory pinned
// for the call.
type arg struct {
	i uintptr
	p unsafe.Pointer
}

// n is the integer v as an argument.
func n[T ~int | ~int32 | ~uint32 | ~uintptr](v T) arg { return arg{i: uintptr(v)} }

// ptr is the Go memory at p as an argument.
func ptr(p unsafe.Pointer) arg { return arg{p: p} }

// call calls fn with args in a frame of its own, and is what it returns:
// the args are all worked out before, so a call among them never shares
// the frame.
func call(fn native.Proc, args ...arg) uintptr {
	var f native.Frame
	for _, a := range args {
		if a.p != nil {
			f.Ptr(a.p)
		} else {
			f.Int(a.i)
		}
	}
	return f.Call(fn)
}

// callFloats calls fn with C floats.
func callFloats(fn native.Proc, fs ...float32) {
	var f native.Frame
	for _, v := range fs {
		f.Float32(v)
	}
	f.Call(fn)
}

// pointer is p as an unsafe.Pointer, for ptr.
func pointer[T any](p *T) unsafe.Pointer { return unsafe.Pointer(p) }

// cstring is s, NUL-terminated, for C.
func cstring(s string) unsafe.Pointer {
	b := make([]byte, len(s)+1)
	copy(b, s)
	return unsafe.Pointer(&b[0])
}
