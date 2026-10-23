package gots

import (
	"go/ast"
	"go/types"
)

// assertOkStmt writes v, ok = x.(T) as statements, rather than a pair
// made and taken apart: a test, then each variable stored. A new v that
// is only read gets the interface's value itself, which never changes.
func (e *emitter) assertOkStmt(w *writer, lhs []ast.Expr, r *ast.TypeAssertExpr, define bool) {
	t := e.typeOf(r.Type)
	x, ok := e.newTmp(), e.newTmp()
	w.line("const %s = %s;", x, e.expr(r.X))
	w.line("const %s = %s;", ok, e.holds(x, t))
	got := e.unboxed(x, t)
	if v := e.newVar(lhs[0], define); v != nil && isValue(t) && e.fn != nil && e.fn.body != nil &&
		!mutates(e.info, e.fn.body, v) && !pinned(e.info, e.fn.body, v) {
		got = x + ".v"
	}
	vals := []string{"(" + ok + " ? " + got + " : " + e.zero(t) + ")", ok}
	froms := []types.Type{t, types.Typ[types.Bool]}
	for i, l := range lhs {
		if id, isId := l.(*ast.Ident); isId && id.Name == "_" {
			continue
		}
		e.store(w, l, e.convertValue(vals[i], froms[i], e.typeOf(l), false), define)
	}
}
