package gots

import (
	"go/ast"
	"go/types"
)

// sharedArg reports whether the i-th argument x of a call of callee may be
// passed as it is, not copied, though it is a struct or array read from a
// variable, or an element or field of one held in it: the callee only reads its parameter, and the variable belongs
// to the function being written and nothing else can reach it during the
// call, so neither side can tell; or, whatever x is, callee is pure, so
// nothing runs to change x while it reads it. A copy is an allocation V8
// cannot avoid.
func (e *emitter) sharedArg(callee *types.Func, i int, x ast.Expr, to types.Type) bool {
	if callee == nil || e.deferring || e.fn == nil || e.fn.body == nil || !isValue(to) || isInterface(to) {
		return false
	}
	if types.Identical(e.typeOf(x), to) && e.readonlyParam(callee, i) && e.pureFunc(callee) {
		return true // nothing can change x while callee reads it
	}
	id, ok := valueRoot(e.info, x).(*ast.Ident)
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

// valueRoot is the variable x is a part of, held in it by value: through
// an array's elements and a struct's fields, not a slice's or a
// pointer's, whose storage others may reach.
func valueRoot(info *types.Info, x ast.Expr) ast.Expr {
	for {
		x = ast.Unparen(x)
		switch n := x.(type) {
		case *ast.IndexExpr:
			if _, ok := info.TypeOf(n.X).Underlying().(*types.Array); !ok {
				return x
			}
			x = n.X
		case *ast.SelectorExpr:
			s := info.Selections[n]
			if s == nil || s.Kind() != types.FieldVal || s.Indirect() {
				return x
			}
			if _, ok := info.TypeOf(n.X).Underlying().(*types.Struct); !ok {
				return x
			}
			x = n.X
		default:
			return x
		}
	}
}
