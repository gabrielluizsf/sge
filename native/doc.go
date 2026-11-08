// Package native calls C libraries (libEGL, libGL, libX11) from Go with
// CGO_ENABLED=0. On linux/amd64 it stands in for runtime/cgo: it links
// libc dynamically, lets libc create every Go thread (so C sees real
// pthreads with initialised TLS) and turns runtime.cgocall on.
package native
