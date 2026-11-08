//go:build !cgo

package native

import "unsafe"

// threadArgs mirrors runtime/cgo's ThreadStart: g0, &m.tls, mstart.
type threadArgs struct{ g, tls, fn uintptr }

// scratch is C memory for pthread/sigset structs, kept off the nosplit
// Go frames.
type scratch struct {
	attr     [64]byte
	all, old [128]byte
	lo, size uintptr
	tid      uintptr
}

var (
	setgGCC       uintptr // runtime's setg_gcc, handed to cgoInit
	keyCreated    uintptr // runtime only checks the pointer is set
	threadEntryPC uintptr // threadEntryC, set in hooks_linux_amd64.s
)

const (
	sigSetmask = 2 // SIG_SETMASK
	detached   = 1 // PTHREAD_CREATE_DETACHED
)

//go:nosplit
func ptr[T any](p *T) uintptr { return uintptr(unsafe.Pointer(p)) }

//go:nosplit
//go:norace
func newScratch(f *hookFrame) *scratch {
	p := f.c(libc.malloc, [4]uintptr{unsafe.Sizeof(scratch{})})
	if p == 0 {
		f.c(libc.abort, [4]uintptr{})
	}
	return (*scratch)(cmem(p))
}

// cgoInit records setg and gives g0 the real bounds of the main stack.
//
//go:nosplit
//go:norace
func cgoInit(g, setg uintptr) {
	setgGCC = setg
	var f hookFrame
	s := newScratch(&f)
	f.c(libc.attrInit, [4]uintptr{ptr(&s.attr)})
	f.c(libc.getattrNP, [4]uintptr{f.c(libc.self, [4]uintptr{}), ptr(&s.attr)})
	f.c(libc.attrGetstack, [4]uintptr{ptr(&s.attr), ptr(&s.lo), ptr(&s.size)})
	f.c(libc.attrDestroy, [4]uintptr{ptr(&s.attr)})
	*(*uintptr)(cmem(g)) = s.lo
	f.c(libc.free, [4]uintptr{ptr(s)})
}

// threadStart makes every new Go M a real pthread, signals blocked
// while it starts, exactly as runtime/cgo does.
//
//go:nosplit
//go:norace
func threadStart(arg uintptr) {
	var f hookFrame
	ts := f.c(libc.malloc, [4]uintptr{unsafe.Sizeof(threadArgs{})})
	if ts == 0 {
		f.c(libc.abort, [4]uintptr{})
	}
	*(*threadArgs)(cmem(ts)) = *(*threadArgs)(cmem(arg))
	s := newScratch(&f)
	f.c(libc.sigfillset, [4]uintptr{ptr(&s.all)})
	f.c(libc.sigmask, [4]uintptr{sigSetmask, ptr(&s.all), ptr(&s.old)})
	f.c(libc.attrInit, [4]uintptr{ptr(&s.attr)})
	f.c(libc.attrSetdetach, [4]uintptr{ptr(&s.attr), detached})
	f.c(libc.attrGetstacksize, [4]uintptr{ptr(&s.attr), ptr(&s.size)})
	(*[2]uintptr)(cmem((*threadArgs)(cmem(ts)).g))[1] = s.size
	err := f.c(libc.create, [4]uintptr{ptr(&s.tid), ptr(&s.attr), threadEntryPC, ts})
	f.c(libc.sigmask, [4]uintptr{sigSetmask, ptr(&s.old), 0})
	f.c(libc.attrDestroy, [4]uintptr{ptr(&s.attr)})
	f.c(libc.free, [4]uintptr{ptr(s)})
	if err != 0 {
		f.c(libc.abort, [4]uintptr{})
	}
}

// threadEntry runs on the new pthread: it sets g and returns mstart,
// which threadEntryC calls itself (mstart is Go code and does not keep
// the C callee-saved registers, so it must not return into Go here).
//
//go:nosplit
//go:norace
func threadEntry(v uintptr) uintptr {
	var f hookFrame
	ts := *(*threadArgs)(cmem(v))
	f.c(libc.free, [4]uintptr{v})
	f.c(setgGCC, [4]uintptr{ts.g})
	return ts.fn
}
