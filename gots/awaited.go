package gots

import "go/ast"

// awaited is the call c, written js, as the value it gives: awaited when
// it may block. A call of a known function that blocks is awaited; one
// through a function value or an interface, whose callee may or may not
// block, waits only for a promise, in a function that can wait at all. A
// call a defer or go statement makes later is left as it is.
func (e *emitter) awaited(c *ast.CallExpr, js string) string {
	if e.bindingCall {
		return js
	}
	async := e.fn != nil && e.fn.async
	if callee := e.staticCallee(c); callee != nil {
		if !e.blocks(callee) {
			return js
		}
		if !async {
			e.errorf(c.Pos(), "unsupported: a call that blocks outside a function that can wait (a package variable's initializer, or a range-over-func body)")
			return js
		}
		return "(await " + js + ")"
	}
	if lit := e.blockingOf().closureCalled(e.info, c); lit != nil {
		b := e.blockingOf()
		if !b.blocks[b.key(lit.Pos())] {
			return js
		}
		if !async {
			e.errorf(c.Pos(), "unsupported: a call that blocks outside a function that can wait")
			return js
		}
		return "(await " + js + ")"
	}
	if !async {
		return js
	}
	return e.awaitIfPromise(js)
}
