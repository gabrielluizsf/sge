//go:build gots_overlay

// Overlay for the functions internal/bytealg implements in assembly or
// takes from the runtime.

package bytealg

import "unsafe"

func MakeNoZero(n int) []byte { return make([]byte, n) }

func Compare(a, b []byte) int {
	n := min(len(a), len(b))
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			if a[i] < b[i] {
				return -1
			}
			return 1
		}
	}
	switch {
	case len(a) < len(b):
		return -1
	case len(a) > len(b):
		return 1
	}
	return 0
}

func abigen_runtime_cmpstring(a, b string) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

func abigen_runtime_memequal(a, b unsafe.Pointer, size uintptr) bool {
	panic("gots: memequal on raw pointers")
}

func abigen_runtime_memequal_varlen(a, b unsafe.Pointer) bool {
	panic("gots: memequal on raw pointers")
}

func IndexByte(b []byte, c byte) int {
	for i, x := range b {
		if x == c {
			return i
		}
	}
	return -1
}

func IndexByteString(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}
