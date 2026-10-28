package gots

import (
	"go/ast"
	"go/types"
)

// expr is the JavaScript for x's value, in x's own type: no copy is made and
// nothing is put in an interface (see value for that).
func (e *emitter) expr(x ast.Expr) string {
	if tv, ok := e.info.Types[x]; ok && tv.Value != nil {
		return e.constant(tv.Value, e.T(tv.Type))
	}
	switch x := x.(type) {
	case *ast.ParenExpr:
		return e.expr(x.X)
	case *ast.Ident:
		return e.ident(x)
	case *ast.BasicLit:
		return e.unsupported(x, "literal without a constant value")
	case *ast.SelectorExpr:
		return e.selector(x)
	case *ast.IndexExpr:
		return e.index(x)
	case *ast.IndexListExpr:
		return e.expr(x.X)
	case *ast.SliceExpr:
		return e.sliceExpr(x)
	case *ast.StarExpr:
		return e.deref(x)
	case *ast.UnaryExpr:
		return e.unary(x)
	case *ast.BinaryExpr:
		return e.binary(x)
	case *ast.CallExpr:
		return e.call(x)
	case *ast.CompositeLit:
		return e.composite(x)
	case *ast.FuncLit:
		return e.funcLit(x)
	case *ast.TypeAssertExpr:
		return e.typeAssert(x)
	case *ast.KeyValueExpr:
		return e.unsupported(x, "key-value outside a composite literal")
	}
	return e.unsupported(x, "expression")
}

func (e *emitter) ident(x *ast.Ident) string {
	switch obj := e.info.ObjectOf(x).(type) {
	case *types.Nil:
		return "null"
	case *types.Var:
		return e.varRef(obj)
	case *types.Func:
		if inst, ok := e.info.Instances[x]; ok && typedNatives[nativeName(obj)] {
			return e.nativeInstance(obj, inst.TypeArgs)
		}
		if e.generic(obj) {
			if inst, ok := e.info.Instances[x]; ok {
				return e.instName(obj, inst.TypeArgs)
			}
		}
		if e.wrapped(obj) {
			return e.objName(obj) + "$g"
		}
		return e.objName(obj)
	case *types.TypeName:
		return e.objName(obj)
	case nil:
		if x.Name == "_" {
			return "undefined"
		}
	}
	return e.unsupported(x, "identifier "+x.Name)
}

// simple reports whether x can be written twice without being evaluated
// twice to a different effect.
func (e *emitter) simple(x ast.Expr) bool {
	if tv, ok := e.info.Types[x]; ok && tv.Value != nil {
		return true
	}
	switch x := x.(type) {
	case *ast.ParenExpr:
		return e.simple(x.X)
	case *ast.Ident:
		return true
	case *ast.SelectorExpr:
		if s := e.info.Selections[x]; s != nil {
			return s.Kind() == types.FieldVal && e.simple(x.X)
		}
		return true
	}
	return false
}
