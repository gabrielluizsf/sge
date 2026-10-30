//go:build linux && amd64

package native

import (
	"runtime"
	"sync"
	"testing"
	"unsafe"
)

// libc's and libm's functions, called as any C library's: an integer and a
// string in, an integer out; a double in and out.
func TestCallsC(t *testing.T) {
	lib, err := Open("libc.so.6")
	if err != nil {
		t.Fatal(err)
	}
	strlen, err := lib.Symbol("strlen")
	if err != nil {
		t.Fatal(err)
	}
	var f Frame
	if n := f.Ptr(unsafe.Pointer(cstr("sge native"))).Call(strlen); n != 10 {
		t.Errorf("strlen is %d, not 10", n)
	}
	libm, err := Open("libm.so.6")
	if err != nil {
		t.Fatal(err)
	}
	fabs, err := libm.Symbol("fabs")
	if err != nil {
		t.Fatal(err)
	}
	f.Float64(-2.5).Call(fabs)
	if v := f.Float(); v != 2.5 {
		t.Errorf("fabs(-2.5) is %v", v)
	}
	if _, err := lib.Symbol("no_such_function_here"); err == nil {
		t.Error("found a symbol that is not there")
	}
}

// C called from many goroutines at once, each locked to its thread while
// the collector runs and threads come and go, keeps working: every
// thread is one libc made, its TLS libc's own (malloc and free use it).
func TestCallsFromManyThreads(t *testing.T) {
	lib, err := Open("libc.so.6")
	if err != nil {
		t.Fatal(err)
	}
	malloc, _ := lib.Symbol("malloc")
	free, _ := lib.Symbol("free")
	var wg sync.WaitGroup
	for range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runtime.LockOSThread() // the thread dies with the goroutine
			var f Frame
			for i := range 2000 {
				p := f.Int(uintptr(16 + i%64)).Call(malloc)
				if p == 0 {
					t.Error("malloc failed")
					return
				}
				f.Int(p).Call(free)
				if i%200 == 0 {
					runtime.GC()
				}
			}
		}()
	}
	wg.Wait()
}

func TestNilProcPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("called through a nil Proc")
		}
	}()
	var f Frame
	f.Call(0)
}
