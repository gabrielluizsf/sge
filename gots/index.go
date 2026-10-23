package gots

import (
	"go/ast"
	"go/types"
)

func (e *emitter) index(x *ast.IndexExpr) string {
	t := e.typeOf(x.X)
	if _, isFunc := t.Underlying().(*types.Signature); isFunc {
		// An instantiation, f[T]: type arguments are erased.
		return e.expr(x.X)
	}
	switch u := t.Underlying().(type) {
	case *types.Map:
		m := e.expr(x.X)
		k := e.value(x.Index, u.Key())
		return "(" + m + "?.get(" + k + ") ?? " + e.zero(u.Elem()) + ")"
	case *types.Slice:
		i := e.expr(x.Index)
		if e.simple(x.X) {
			s := e.expr(x.X)
			return s + ".a[" + s + ".o + " + i + "]"
		}
		return "$.at(" + e.expr(x.X) + ", " + i + ")"
	case *types.Basic:
		return e.expr(x.X) + ".charCodeAt(" + e.expr(x.Index) + ")"
	}
	// An array, or a pointer to one.
	return e.expr(x.X) + "[" + e.expr(x.Index) + "]"
}

func (e *emitter) sliceExpr(x *ast.SliceExpr) string {
	t := e.typeOf(x.X)
	opt := func(b ast.Expr) string {
		if b == nil {
			return "undefined"
		}
		return e.expr(b)
	}
	lo := "0"
	if x.Low != nil {
		lo = e.expr(x.Low)
	}
	if isString(t) {
		if x.High == nil {
			return e.expr(x.X) + ".substring(" + lo + ")"
		}
		return e.expr(x.X) + ".substring(" + lo + ", " + e.expr(x.High) + ")"
	}
	if _, ok := t.Underlying().(*types.Slice); ok {
		return "$.sub(" + e.expr(x.X) + ", " + lo + ", " + opt(x.High) + ", " + opt(x.Max) + ")"
	}
	return "$.subArray(" + e.expr(x.X) + ", " + lo + ", " + opt(x.High) + ", " + opt(x.Max) + ")"
}
