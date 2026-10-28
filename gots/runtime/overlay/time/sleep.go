//go:build gots_overlay

// Overlay for sleep.go and tick.go: timers run on JavaScript's setTimeout
// and setInterval, delivering on their channel without waiting, as Go's
// buffered timer channels do, or calling AfterFunc's function. Sleep is a
// receive from a timer, so a sleeping goroutine lets the others run.

package time

import (
	"syscall/js"
	"unsafe"
)

type Timer struct {
	C         <-chan Time
	c         chan Time // C, to send on
	initTimer bool
	handle    js.Value
	fn        js.Func
	f         func()
	active    bool
}

// fire is what the timer does when it goes off.
func (t *Timer) fire() {
	t.active = false
	if t.f != nil {
		t.f()
		return
	}
	select {
	case t.c <- Now():
	default:
	}
}

func (t *Timer) Stop() bool {
	if t.handle.IsUndefined() {
		panic("time: Stop called on uninitialized Timer")
	}
	was := t.active
	t.active = false
	js.Global().Call("clearTimeout", t.handle)
	t.drain()
	return was
}

func (t *Timer) Reset(d Duration) bool {
	was := t.Stop()
	t.start(d)
	return was
}

// drain drops a time the timer sent but nobody received: as of Go 1.23 a
// stopped or reset timer delivers nothing stale.
func (t *Timer) drain() {
	if t.c == nil {
		return
	}
	select {
	case <-t.c:
	default:
	}
}

func (t *Timer) start(d Duration) {
	t.active = true
	t.handle = js.Global().Call("setTimeout", t.fn, float64(d)/1e6)
}

func newJSTimer(f func()) *Timer {
	t := &Timer{f: f}
	if f == nil {
		t.c = make(chan Time, 1)
		t.C = t.c
	}
	t.fn = js.FuncOf(func(this js.Value, args []js.Value) any {
		t.fire()
		return nil
	})
	return t
}

func NewTimer(d Duration) *Timer {
	t := newJSTimer(nil)
	t.start(d)
	return t
}

func After(d Duration) <-chan Time {
	return NewTimer(d).C
}

func AfterFunc(d Duration, f func()) *Timer {
	t := newJSTimer(f)
	t.start(d)
	return t
}

// Sleep pauses the goroutine for d.
func Sleep(d Duration) {
	if d <= 0 {
		return
	}
	<-After(d)
}

type Ticker struct {
	C          <-chan Time
	c          chan Time
	initTicker bool
	handle     js.Value
	fn         js.Func
}

func NewTicker(d Duration) *Ticker {
	if d <= 0 {
		panic("non-positive interval for NewTicker")
	}
	c := make(chan Time, 1)
	t := &Ticker{C: c, c: c}
	t.fn = js.FuncOf(func(this js.Value, args []js.Value) any {
		select {
		case t.c <- Now():
		default:
		}
		return nil
	})
	t.start(d)
	return t
}

func (t *Ticker) start(d Duration) {
	t.handle = js.Global().Call("setInterval", t.fn, float64(d)/1e6)
}

func (t *Ticker) Stop() {
	if !t.handle.IsUndefined() {
		js.Global().Call("clearInterval", t.handle)
	}
}

func (t *Ticker) Reset(d Duration) {
	if d <= 0 {
		panic("non-positive interval for Ticker.Reset")
	}
	t.Stop()
	t.start(d)
}

func Tick(d Duration) <-chan Time {
	if d <= 0 {
		return nil
	}
	return NewTicker(d).C
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
