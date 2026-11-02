package gots

import (
	"go/ast"
	"go/types"
)

// dyingLocal reports whether the result r of return s may be the value
// it names itself, not a copy: a struct or array variable of the
// function's own, declared in its body, not a parameter, which a caller
// may have lent it, returned once, of the result's very type, which no
// closure or pointer keeps, so with the function gone nothing else holds
// it. A copy is an allocation V8 cannot avoid.
func (e *emitter) dyingLocal(s *ast.ReturnStmt, r ast.Expr, to types.Type) bool {
	id, ok := ast.Unparen(r).(*ast.Ident)
	if !ok || e.fn == nil || e.fn.body == nil || !isValue(to) || isInterface(to) || !types.Identical(e.typeOf(r), to) {
		return false
	}
	v, ok := e.info.Uses[id].(*types.Var)
	if !ok || v.IsField() || e.boxed[v] || v.Pos() < e.fn.body.Pos() || v.Pos() >= e.fn.body.End() || pinned(e.info, e.fn.body, v) {
		return false
	}
	n := 0
	for _, x := range s.Results {
		ast.Inspect(x, func(m ast.Node) bool {
			if i, ok := m.(*ast.Ident); ok && e.info.Uses[i] == v {
				n++
			}
			return true
		})
	}
	return n == 1
}
