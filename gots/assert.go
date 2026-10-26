package gots

import (
	"go/ast"
	"go/types"
)

// typeAssert is x.(T), panicking when x does not hold a T.
func (e *emitter) typeAssert(x *ast.TypeAssertExpr) string {
	v := e.expr(x.X)
	t := e.typeOf(x.Type)
	switch {
	case isInterface(t):
		return "$.assertIface(" + v + ", " + jsString(typeString(t)) + ", " + e.hoistConst("$.methodSet("+methodList(t.Underlying().(*types.Interface))+")") + ")"
	case isStructPtr(t):
		return "$.assertPtr(" + v + ", " + e.structClass(derefType(t)) + ")"
	}
	return "$.assertVal(" + v + ", " + e.desc(t) + ")"
}

func isStructPtr(t types.Type) bool {
	p, ok := t.Underlying().(*types.Pointer)
	return ok && isStruct(p.Elem())
}

// holds is a test that the interface value in the binding v holds a T.
func (e *emitter) holds(v string, t types.Type) string {
	switch {
	case isNil(t):
		return "(" + v + " === null)"
	case isInterface(t):
		return "$.implementsI(" + v + ", " + e.hoistConst("$.methodSet("+methodList(t.Underlying().(*types.Interface))+")") + ")"
	case isStructPtr(t):
		return "$.isPtrTo(" + v + ", " + e.structClass(derefType(t)) + ")"
	}
	// Inline, without instanceof: only a Val has a $t.
	return "(" + v + "?.$t === " + e.desc(t) + ")"
}

// unboxed is the T held by the interface value in the binding v, once
// holds has said it is one.
func (e *emitter) unboxed(v string, t types.Type) string {
	if isStructPtr(t) {
		return "$.unboxPtr(" + v + ")"
	}
	if isInterface(t) || isNil(t) {
		return v
	}
	return e.cloneOf(v+".v", t)
}

// mapOk is v, ok := m[k] as an expression for the pair.
func (e *emitter) mapOk(x *ast.IndexExpr) string {
	m := e.typeOf(x.X).Underlying().(*types.Map)
	v := e.newTmp()
	return "((" + v + ": any) => " + v + " === undefined ? [" + e.zero(m.Elem()) + ", false] : [" + e.cloneOf(v, m.Elem()) + ", true])(" + e.expr(x.X) + "?.get(" + e.value(x.Index, m.Key()) + "))"
}
