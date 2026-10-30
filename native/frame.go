package native

import (
	"math"
	"runtime"
	"unsafe"
)

// Frame is one C call: up to 16 integer/pointer arguments (6 in
// registers, the rest on the stack) and 8 float arguments. The layout of
// the first fields is read by call_amd64.s; do not reorder them.
type Frame struct {
	fn     uintptr
	ints   [16]uintptr
	floats [8]uint64
	ret    uintptr
	fret   uint64
	ni, nf int
	pin    runtime.Pinner
}

// Int appends an integer argument.
func (f *Frame) Int(v uintptr) *Frame {
	f.ints[f.ni] = v
	f.ni++
	return f
}

// Ptr appends a pointer argument, pinned until the call returns.
func (f *Frame) Ptr(p unsafe.Pointer) *Frame {
	if p != nil {
		f.pin.Pin(p)
	}
	return f.Int(uintptr(p))
}

// Float32 appends a C float argument.
func (f *Frame) Float32(v float32) *Frame {
	f.floats[f.nf] = uint64(math.Float32bits(v))
	f.nf++
	return f
}

// Float64 appends a C double argument.
func (f *Frame) Float64(v float64) *Frame {
	f.floats[f.nf] = math.Float64bits(v)
	f.nf++
	return f
}

// Call calls fn with the arguments appended so far, then resets f.
func (f *Frame) Call(fn Proc) uintptr {
	if fn == 0 {
		panic("native: call through a nil Proc")
	}
	f.fn = uintptr(fn)
	call(f)
	f.pin.Unpin()
	f.ni, f.nf = 0, 0
	return f.ret
}

// Float returns the float/double result of the last call.
func (f *Frame) Float() float64 { return math.Float64frombits(f.fret) }
