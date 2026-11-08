//go:build cgo && linux && amd64

package native

// With cgo, runtime/cgo starts the threads and sets their TLS, as the
// hooks do without it; everything else is the same.
import _ "runtime/cgo"
