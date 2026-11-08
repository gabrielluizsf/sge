package native

import (
	"unsafe"
	_ "unsafe" // go:linkname
)

// runtime.cgocall is pushed by the runtime itself (//go:linkname cgocall,
// "hall of shame" for purego), so -checklinkname lets any package pull it.
//
//go:linkname cgocall runtime.cgocall
//go:noescape
func cgocall(fn, arg unsafe.Pointer) int32

var callFramePC uintptr // callFrameC, set in call_amd64.s

// call runs f on the M's g0 (pthread) stack, as cgo does: the goroutine
// enters a syscall state, so the scheduler, the GC and async preemption
// leave it alone while C runs.
func call(f *Frame) { cgocall(*(*unsafe.Pointer)(unsafe.Pointer(&callFramePC)), unsafe.Pointer(f)) }
