//go:build gots_overlay

// Overlay for weak pointers, which JavaScript could hold with WeakRef but
// gots does not yet: a weak pointer here keeps what it points to, so it
// never turns nil.

package weak

import "unsafe"

func runtime_registerWeakPointer(p unsafe.Pointer) unsafe.Pointer { return p }

func runtime_makeStrongFromWeak(h unsafe.Pointer) unsafe.Pointer { return h }
