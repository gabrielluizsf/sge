package gots

import (
	"go/ast"
	"go/types"
)

// A function that blocks only through what it calls by interface or
// function value, say bufio.Reader.Read over a socket, has a twin: the
// function itself, for code that cannot wait, which never waits, and
// name$a, for code that can, which awaits what it calls. Engine code
// drawing a frame stays synchronous; a goroutine serving HTTP waits on
// its socket.

// twinFunc reports whether the function or closure of type ft has a twin.
func (e *emitter) twinFunc(ft *ast.FuncType) bool {
	if ft == nil || !ft.Pos().IsValid() {
		return false
	}
	k := e.blockingOf().key(ft.Pos())
	return e.mayBlockOf().blocks[k] && !e.blockingOf().blocks[k]
}

// twin reports whether calling fn may block through what it calls by
// interface, though fn itself is not one that blocks for sure: whether it
// has a twin.
func (e *emitter) twin(fn *types.Func) bool {
	if twinNatives[nativeName(fn)] {
		return true
	}
	if blockingNatives[nativeName(fn)] {
		return false
	}
	s := fn.Origin().Scope()
	if s == nil || !s.Pos().IsValid() {
		return false
	}
	k := e.blockingOf().key(s.Pos())
	return e.mayBlockOf().blocks[k] && !e.blockingOf().blocks[k]
}

// ifaceMayBlock reports whether a method fn of an interface may block in
// some type that has it.
func (e *emitter) ifaceMayBlock(fn *types.Func) bool {
	return e.mayBlockOf().blocks[ifaceKey(fn)]
}

// canWait reports whether the code being written may wait: it is in an
// async function, or it is the call a go statement makes.
func (e *emitter) canWait() bool {
	return e.goCall || e.fn != nil && e.fn.async
}

// functionTwins writes function and, for a function with a twin, the twin
// named prefix$a, which may wait.
func (e *emitter) functionTwins(w *writer, prefix string, sig *types.Signature, ft *ast.FuncType, body *ast.BlockStmt, prologue []string, fresh bool) {
	e.function(w, prefix, sig, ft, body, prologue, false, fresh)
	if e.twinFunc(ft) {
		e.waiting = true
		e.function(w, prefix+"$a", sig, ft, body, prologue, false, fresh)
	}
}

// entry is the function main, a test or a benchmark is run as: its twin,
// if it has one, as they all may wait.
func (e *emitter) entry(name string) string {
	if fn, ok := e.tpkg.Scope().Lookup(name).(*types.Func); ok && e.twin(fn) {
		return name + "$a"
	}
	return name
}
