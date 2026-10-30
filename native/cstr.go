package native

import "unsafe"

// cstr returns a NUL-terminated heap copy of s.
func cstr(s string) *byte {
	b := make([]byte, len(s)+1)
	copy(b, s)
	return &b[0]
}

// GoString copies a NUL-terminated C string.
func GoString(p uintptr) string {
	if p == 0 {
		return ""
	}
	n := 0
	for *(*byte)(unsafe.Add(cmem(p), n)) != 0 {
		n++
	}
	return string(unsafe.Slice((*byte)(cmem(p)), n))
}
