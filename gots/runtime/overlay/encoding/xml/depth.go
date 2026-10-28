//go:build gots_overlay

// Overlay for the depth XML may nest to when decoded: Go's limits are
// sized for a goroutine stack that grows, but a JavaScript stack is fixed,
// and XML from outside nested a few thousand deep would overflow it and
// end the program (CVE-2026-56859's stack exhaustion, by another route).
// Code V8 has not optimized yet runs out at about 475 levels, so past 256
// Unmarshal reports an error instead.

package xml

const (
	maxUnmarshalDepth     = 256
	maxUnmarshalDepthWasm = 256
)
