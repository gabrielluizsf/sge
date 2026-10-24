package gots

import (
	"go/ast"
	"go/types"
)

// sendStmt is ch <- v, waiting until the channel takes v; the runtime
// hands back a promise only when it has to wait.
func (e *emitter) sendStmt(w *writer, s *ast.SendStmt) {
	if e.fn == nil || !e.fn.async {
		w.line("%s;", e.unsupported(s, "a send outside a function that can wait"))
		return
	}
	elem := e.typeOf(s.Chan).Underlying().(*types.Chan).Elem()
	p := e.fn.scratch
	w.line("if ((%s = $.ch.send(%s, %s)) !== undefined) await %s;", p, e.expr(s.Chan), e.value(s.Value, elem), p)
}

// recvExpr is <-ch: the pair of what was received and whether the
// channel is still open, once one comes.
func (e *emitter) recvExpr(x *ast.UnaryExpr) string {
	if e.fn == nil || !e.fn.async {
		return e.unsupported(x, "a receive outside a function that can wait")
	}
	return e.awaitIfPromise("$.ch.recv(" + e.expr(x.X) + ")")
}

// awaitIfPromise is js, awaited only when it is a promise: a channel
// operation that needs no waiting gives its result at once, and awaiting
// it anyway costs a turn of the microtask queue, five times the operation.
func (e *emitter) awaitIfPromise(js string) string {
	p := e.fn.scratch
	return "((" + p + " = " + js + ") instanceof Promise ? await " + p + " : " + p + ")"
}
