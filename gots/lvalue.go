package gots

import (
	"go/ast"
	"go/types"
)

// lvalue is a place a value can be stored: how to read it and how to write
// to it.
type lvalue struct {
	get   string
	set   func(v string) string
	t     types.Type
	blank bool
}

// hoist is x as an expression that can be written twice: x itself when
// that is safe, or a constant bound before the statement.
func (e *emitter) hoist(w *writer, x ast.Expr) string {
	if e.simple(x) {
		return e.expr(x)
	}
	t := e.newTmp()
	w.line("const %s = %s;", t, e.expr(x))
	return t
}

// place is the lvalue x denotes. Operands with effects are evaluated once,
// in statements written to w ahead of the store.
func (e *emitter) place(w *writer, x ast.Expr) lvalue {
	x = ast.Unparen(x)
	t := e.typeOf(x)
	inPlace := func(get string) lvalue {
		return lvalue{get: get, t: t, set: func(v string) string {
			if isValue(t) {
				return e.setOf(get, v, t)
			}
			return get + " = " + v
		}}
	}
	switch x := x.(type) {
	case *ast.Ident:
		if x.Name == "_" {
			return lvalue{blank: true, set: func(v string) string { return v }}
		}
		obj := e.info.ObjectOf(x)
		if isPackageLevel(obj) && obj.Pkg() != e.tpkg {
			get := e.objName(obj)
			return lvalue{get: get, t: t, set: func(v string) string {
				if isValue(t) {
					return e.setOf(get, v, t)
				}
				return e.setPackageVar(obj, v)
			}}
		}
		return inPlace(e.varRef(obj))
	case *ast.SelectorExpr:
		if s := e.info.Selections[x]; s != nil {
			base := e.hoist(w, x.X)
			get, _ := e.fieldPath(base, s.Recv(), s.Index())
			return inPlace(get)
		}
		obj := e.info.ObjectOf(x.Sel)
		get := e.objName(obj)
		return lvalue{get: get, t: t, set: func(v string) string {
			if isValue(t) {
				return e.setOf(get, v, t)
			}
			return e.setPackageVar(obj, v)
		}}
	case *ast.IndexExpr:
		switch u := e.typeOf(x.X).Underlying().(type) {
		case *types.Map:
			m := e.hoist(w, x.X)
			k := e.value(x.Index, u.Key())
			if !e.simple(x.Index) {
				kt := e.newTmp()
				w.line("const %s = %s;", kt, k)
				k = kt
			}
			return lvalue{
				get: "(" + m + "?.get(" + k + ") ?? " + e.zero(u.Elem()) + ")",
				t:   t,
				set: func(v string) string { return "$.mapSet(" + m + ", " + k + ", " + v + ")" },
			}
		case *types.Slice:
			s := e.hoist(w, x.X)
			i := e.hoist(w, x.Index)
			return inPlace(s + ".a[" + s + ".o + " + i + "]")
		}
		a := e.hoist(w, x.X)
		i := e.hoist(w, x.Index)
		return inPlace(a + "[" + i + "]")
	case *ast.StarExpr:
		p := e.hoist(w, x.X)
		if isValue(t) {
			return inPlace(p)
		}
		return inPlace(p + ".v")
	}
	e.unsupported(x, "assignment to this expression")
	return lvalue{blank: true, set: func(v string) string { return v }}
}
