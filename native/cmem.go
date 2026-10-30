package native

import "unsafe"

// cmem is the C memory at address p, which Go's collector neither moves
// nor frees: an unsafe.Pointer to it, made without the uintptr
// conversion vet rightly warns about for Go memory.
//
//go:nosplit
func cmem(p uintptr) unsafe.Pointer { return *(*unsafe.Pointer)(unsafe.Pointer(&p)) }
