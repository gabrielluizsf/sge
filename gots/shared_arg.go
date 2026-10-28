package gots

import (
	"go/ast"
	"go/types"
)

// sharedArg reports whether the i-th argument x of a call of callee may be
// passed as it is, not copied, though it is a struct or array read from a
// variable: the callee only reads its parameter, and the variable belongs
// to the function being written and nothing else can reach it during the
// call, so neither side can tell. A copy is an allocation V8 cannot avoid.
func (e *emitter) sharedArg(callee *types.Func, i int, x ast.Expr, to types.Type) bool {
	if callee == nil || e.deferring || e.fn == nil || e.fn.body == nil || !isValue(to) || isInterface(to) {
		return false
	}
	id, ok := ast.Unparen(x).(*ast.Ident)
	if !ok {
		return false
	}
	v, ok := e.info.Uses[id].(*types.Var)
	if !ok || v.IsField() || e.boxed[v] || !types.Identical(e.typeOf(x), to) {
		return false
	}
	// Declared in this function, its parameters included, not in one
	// around it.
	if v.Pos() < e.fn.ft.Pos() || v.Pos() >= e.fn.body.End() {
		return false
	}
	return !pinned(e.info, e.fn.body, v) && e.readonlyParam(callee, i)
}
