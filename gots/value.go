package gots

import (
	"go/ast"
	"go/types"
)

// value is x as it is stored into a place of type to: copied when it is a
// struct or array read from storage, and put in an interface when to is one.
func (e *emitter) value(x ast.Expr, to types.Type) string {
	return e.convertValue(e.expr(x), e.typeOf(x), to, e.needsCopy(x))
}

func (e *emitter) convertValue(js string, from, to types.Type, copy bool) string {
	if from == nil {
		return js
	}
	if b, ok := from.(*types.Basic); ok && b.Kind() == types.UntypedNil {
		return "null"
	}
	if copy && isValue(from) {
		js = e.cloneOf(js, from)
	}
	if to != nil && isInterface(to) && !isInterface(from) {
		return e.box(js, from)
	}
	return js
}

// box puts a value of the concrete type t into an interface.
func (e *emitter) box(js string, t types.Type) string {
	if p, ok := t.Underlying().(*types.Pointer); ok && isStruct(p.Elem()) {
		return js
	}
	if cls := e.valClassOf(t); cls != "" {
		// What $.val would pick, named here: a construction V8 can inline.
		return "new " + cls + "(" + e.desc(t) + ", " + js + ")"
	}
	return "$.val(" + e.desc(t) + ", " + js + ")"
}

// needsCopy reports whether x reads a struct or array out of storage that
// something else can still reach, so storing it elsewhere must copy it.
func (e *emitter) needsCopy(x ast.Expr) bool {
	switch x := ast.Unparen(x).(type) {
	case *ast.Ident:
		_, isVar := e.info.ObjectOf(x).(*types.Var)
		return isVar
	case *ast.SelectorExpr:
		s := e.info.Selections[x]
		if s != nil {
			return s.Kind() == types.FieldVal
		}
		_, isVar := e.info.ObjectOf(x.Sel).(*types.Var)
		return isVar
	case *ast.IndexExpr, *ast.StarExpr, *ast.TypeAssertExpr:
		return true
	case *ast.CallExpr:
		if tv, ok := e.info.Types[x.Fun]; ok && tv.IsType() && len(x.Args) == 1 {
			to, from := e.T(tv.Type), e.typeOf(x.Args[0])
			return isStruct(to) && e.sameClass(to, from) && e.needsCopy(x.Args[0]) || isArray(to) && e.needsCopy(x.Args[0])
		}
	}
	return false
}

// sameClass reports whether struct types a and b are held in one class.
func (e *emitter) sameClass(a, b types.Type) bool {
	return e.structClass(a) == e.structClass(b)
}
