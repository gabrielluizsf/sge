//go:build gots_overlay

// Overlay for the depth JSON may nest to: Go follows 10000 levels on a
// goroutine stack that grows, but a JavaScript stack is fixed, and JSON from
// outside nested a few thousand deep would overflow it and end the program
// (the stack exhaustion of CVE-2026-56859 and CVE-2026-33818, in encoding/xml
// and encoding/asn1). Past 1000 levels json reports an error instead.

package json

const maxNestingDepth = 1000
