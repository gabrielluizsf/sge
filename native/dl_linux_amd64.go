package native

import (
	"errors"
	"unsafe"
)

// Proc is the address of a C function.
type Proc uintptr

// Lib is a library opened with dlopen.
type Lib struct{ h uintptr }

const rtldNow = 2 // RTLD_NOW | RTLD_LOCAL

// Open loads a system library by its bare soname ("libEGL.so.1"). Paths
// are refused: the dynamic loader's own search (ld.so.cache, system
// dirs) picks the file, never the working directory.
func Open(soname string) (*Lib, error) {
	if !validSoname(soname) {
		return nil, errors.New("native: refusing library name " + soname)
	}
	var f Frame
	h := f.Ptr(unsafe.Pointer(cstr(soname))).Int(rtldNow).Call(Proc(libc.dlopen))
	if h == 0 {
		return nil, errors.New("native: dlopen " + soname + ": " + dlerror())
	}
	return &Lib{h: h}, nil
}

// Symbol looks up a function in the library.
func (l *Lib) Symbol(name string) (Proc, error) {
	var f Frame
	p := f.Int(l.h).Ptr(unsafe.Pointer(cstr(name))).Call(Proc(libc.dlsym))
	if p == 0 {
		return 0, errors.New("native: dlsym " + name + ": " + dlerror())
	}
	return Proc(p), nil
}

func dlerror() string {
	var f Frame
	return GoString(f.Call(Proc(libc.dlerror)))
}
