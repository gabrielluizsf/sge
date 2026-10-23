package gots

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"
)

func (e *emitter) unary(x *ast.UnaryExpr) string {
	switch x.Op {
	case token.AND:
		return e.addr(x.X)
	case token.ARROW:
		return e.unsupported(x, "channel receive")
	}
	a := e.expr(x.X)
	t := e.typeOf(x)
	k := basicKind(t)
	switch x.Op {
	case token.ADD:
		return a
	case token.NOT:
		return "!" + a
	case token.SUB:
		if isFloat(t) {
			return "(-" + a + ")"
		}
		if k == types.Complex64 || k == types.Complex128 {
			return "$.cplx.neg(" + a + ")"
		}
		if k == types.Uint64 || small(k) {
			return wrap(k, "-"+a)
		}
		return "(-" + a + ")"
	case token.XOR:
		switch {
		case k == types.Uint64:
			return wrap(k, "~"+a)
		case small(k) && unsigned(k):
			return wrap(k, "~"+a)
		case small(k):
			return "(~" + a + ")"
		}
		return "(-" + a + " - 1)"
	}
	return e.unsupported(x, "unary "+x.Op.String())
}

// addr is &x.
func (e *emitter) addr(x ast.Expr) string {
	x = ast.Unparen(x)
	t := e.typeOf(x)
	switch x := x.(type) {
	case *ast.CompositeLit:
		if isValue(t) {
			return e.expr(x)
		}
		return "new $.Cell(" + e.expr(x) + ")"
	case *ast.StarExpr:
		return e.expr(x.X)
	}
	if isValue(t) {
		return e.expr(x)
	}
	switch x := x.(type) {
	case *ast.Ident:
		obj := e.info.ObjectOf(x)
		if isPackageLevel(obj) {
			n := e.objName(obj)
			return "new $.VarPtr(() => " + n + ", (v: any) => { " + e.setPackageVar(obj, "v") + "; })"
		}
		return e.objName(obj)
	case *ast.SelectorExpr:
		if e.info.Selections[x] == nil {
			// &pkg.Var
			obj := e.info.ObjectOf(x.Sel)
			n := e.objName(obj)
			return "new $.VarPtr(() => " + n + ", (v: any) => { " + e.setPackageVar(obj, "v") + "; })"
		}
		js := e.expr(x)
		dot := strings.LastIndex(js, ".")
		return "new $.FieldPtr(" + js[:dot] + ", " + jsQuote(js[dot+1:]) + ")"
	case *ast.IndexExpr:
		if _, ok := e.typeOf(x.X).Underlying().(*types.Slice); ok {
			return "$.elemPtr(" + e.expr(x.X) + ", " + e.expr(x.Index) + ")"
		}
		return "new $.IndexPtr(" + e.expr(x.X) + ", " + e.expr(x.Index) + ")"
	}
	return e.unsupported(x, "address of this expression")
}

// setPackageVar is a statement storing v in a package-level variable, which
// another module can only do through the variable's setter.
func (e *emitter) setPackageVar(obj types.Object, v string) string {
	if obj.Pkg() != e.tpkg {
		return e.importAlias(obj.Pkg()) + "." + mangle(obj.Name()) + "$set(" + v + ")"
	}
	return mangle(obj.Name()) + " = " + v
}

// deref is *p.
func (e *emitter) deref(x *ast.StarExpr) string {
	if js, ok := e.reinterpret(x); ok {
		return js
	}
	if isValue(e.typeOf(x)) {
		return e.expr(x.X)
	}
	return e.expr(x.X) + ".v"
}

// reinterpret recognizes *(*T)(unsafe.Pointer(&v)), reading the bits of a
// number as another number type of the same size, and does it the way Go
// does: through memory, here a DataView.
func (e *emitter) reinterpret(x *ast.StarExpr) (string, bool) {
	conv, ok := ast.Unparen(x.X).(*ast.CallExpr)
	if !ok || len(conv.Args) != 1 {
		return "", false
	}
	inner, ok := ast.Unparen(conv.Args[0]).(*ast.CallExpr)
	if !ok || len(inner.Args) != 1 || basicKind(e.typeOf(inner)) != types.UnsafePointer {
		return "", false
	}
	addr, ok := ast.Unparen(inner.Args[0]).(*ast.UnaryExpr)
	if !ok || addr.Op != token.AND {
		return "", false
	}
	from, to := e.typeOf(addr.X), e.typeOf(x)
	if !isNumeric(from) || !isNumeric(to) || wasmSizes.Sizeof(from) != wasmSizes.Sizeof(to) {
		return "", false
	}
	return fmt.Sprintf("$.unsafe.reinterpret(%s, %q, %q)", e.expr(addr.X), types.Typ[defaultKind(basicKind(from))].Name(), types.Typ[defaultKind(basicKind(to))].Name()), true
}
