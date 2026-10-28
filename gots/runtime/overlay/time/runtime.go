//go:build gots_overlay

// Overlay for the clock time takes from the Go runtime: here it comes from
// JavaScript's Date and performance.

package time

import "syscall/js"

func now() (sec int64, nsec int32, mono int64) {
	ms := js.Global().Get("Date").Call("now").Float()
	sec = int64(ms / 1000)
	nsec = int32((ms - float64(sec)*1000) * 1e6)
	return sec, nsec, runtimeNano()
}

func runtimeNow() (sec int64, nsec int32, mono int64) { return now() }

// runtimeNano is the monotonic clock, in nanoseconds since the program
// started.
func runtimeNano() int64 {
	return int64(js.Global().Get("performance").Call("now").Float() * 1e6)
}

func runtimeIsBubbled() bool { return false }
