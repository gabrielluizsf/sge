//go:build gots_overlay

// Overlay for the clock and sleep time takes from the Go runtime: here they
// come from JavaScript's Date, performance and Atomics.

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

var sleeper js.Value

// Sleep blocks the thread for d, which is all a JavaScript program with no
// goroutines can do with it.
func Sleep(d Duration) {
	if d <= 0 {
		return
	}
	if sleeper.IsUndefined() {
		buf := js.Global().Get("SharedArrayBuffer").New(4)
		sleeper = js.Global().Get("Int32Array").New(buf)
	}
	js.Global().Get("Atomics").Call("wait", sleeper, 0, 0, float64(d)/1e6)
}
