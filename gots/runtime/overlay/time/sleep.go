//go:build gots_overlay

// Overlay for sleep.go and tick.go: timers run their function from
// JavaScript's setTimeout. Timers that deliver on a channel need
// goroutines, which the transpiler does not have yet.

package time

import (
	"syscall/js"
	"unsafe"
)

type Timer struct {
	C         <-chan Time
	initTimer bool
	handle    js.Value
	fn        js.Func
	f         func()
	active    bool
}

func (t *Timer) Stop() bool {
	if t.handle.IsUndefined() {
		panic("time: Stop called on uninitialized Timer")
	}
	was := t.active
	t.active = false
	js.Global().Call("clearTimeout", t.handle)
	return was
}

func (t *Timer) Reset(d Duration) bool {
	was := t.Stop()
	t.start(d)
	return was
}

func (t *Timer) start(d Duration) {
	t.active = true
	t.handle = js.Global().Call("setTimeout", t.fn, float64(d)/1e6)
}

func NewTimer(d Duration) *Timer {
	panic("gots: time.NewTimer delivers on a channel, which needs goroutines; use time.AfterFunc")
}

func After(d Duration) <-chan Time {
	panic("gots: time.After delivers on a channel, which needs goroutines; use time.AfterFunc")
}

func AfterFunc(d Duration, f func()) *Timer {
	t := &Timer{f: f}
	t.fn = js.FuncOf(func(this js.Value, args []js.Value) any {
		t.active = false
		t.f()
		return nil
	})
	t.start(d)
	return t
}

type Ticker struct {
	C          <-chan Time
	initTicker bool
}

func NewTicker(d Duration) *Ticker {
	panic("gots: time.NewTicker delivers on a channel, which needs goroutines")
}

func (t *Ticker) Stop() {}

func (t *Ticker) Reset(d Duration) {}

func Tick(d Duration) <-chan Time {
	panic("gots: time.Tick delivers on a channel, which needs goroutines")
}

// The runtime's timer hooks, replaced so the originals' bodiless
// declarations go; nothing above calls them.

func syncTimer(c chan Time) unsafe.Pointer { return nil }

func when(d Duration) int64 { return runtimeNano() + int64(d) }

func newTimer(when, period int64, f func(any, uintptr, int64), arg any, cp unsafe.Pointer) *Timer {
	panic("gots: runtime timers")
}

func stopTimer(*Timer) bool { return false }

func resetTimer(t *Timer, when, period int64) bool { return false }

func sendTime(c any, seq uintptr, delta int64) {}

func goFunc(arg any, seq uintptr, delta int64) {}
