package gots

import (
	"go/ast"
	"go/token"
	"go/types"
)

// pureFunc reports whether fn, while it runs, changes nothing its caller
// could see: it writes only its own variables, held by value; starts no
// goroutine, defers nothing, makes no closure, sends or receives nothing,
// so it never yields to other code; and calls only conversions, builtins
// that change nothing, package math's functions, and pure functions. What
// fn is given by value then cannot change during the call, whoever else
// holds it, and need not be copied for fn to read it.
func (e *emitter) pureFunc(fn *types.Func) bool {
	fn = fn.Origin()
	if e.pure == nil {
		e.pure = map[*types.Func]bool{}
	}
	if p, ok := e.pure[fn]; ok {
		return p
	}
	e.pure[fn] = false // a function calling itself is not pure, as yet
	d, ok := e.funcDeclOf(fn)
	p := ok && e.pureBody(d)
	e.pure[fn] = p
	return p
}

// pureBody reports whether d's body is pure, as pureFunc has it.
func (e *emitter) pureBody(d funcDecl) bool {
	own := func(x ast.Expr) bool {
		if id, ok := ast.Unparen(x).(*ast.Ident); ok && id.Name == "_" {
			return true
		}
		v, ok := rootVar(d.info, x).(*types.Var)
		return ok && v.Pos() >= d.decl.Pos() && v.Pos() < d.decl.End()
	}
	pure := true
	ast.Inspect(d.decl.Body, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.GoStmt, *ast.DeferStmt, *ast.SendStmt, *ast.SelectStmt, *ast.FuncLit:
			pure = false
		case *ast.UnaryExpr:
			pure = pure && n.Op != token.ARROW
		case *ast.AssignStmt:
			for _, l := range n.Lhs {
				pure = pure && own(l)
			}
		case *ast.IncDecStmt:
			pure = pure && own(n.X)
		case *ast.RangeStmt:
			if n.Tok == token.ASSIGN {
				pure = pure && (n.Key == nil || own(n.Key)) && (n.Value == nil || own(n.Value))
			}
			if _, ok := d.info.TypeOf(n.X).Underlying().(*types.Chan); ok {
				pure = false
			}
		case *ast.CallExpr:
			pure = pure && e.pureCall(d.info, n)
		}
		return pure
	})
	return pure
}

// pureCall reports whether the call c changes nothing outside the
// function making it: a conversion, a builtin that writes nothing, a
// function of package math, or a pure function.
func (e *emitter) pureCall(info *types.Info, c *ast.CallExpr) bool {
	fun := ast.Unparen(c.Fun)
	if tv, ok := info.Types[fun]; ok && tv.IsType() {
		return true
	}
	var obj types.Object
	switch f := fun.(type) {
	case *ast.Ident:
		obj = info.Uses[f]
	case *ast.SelectorExpr:
		if s := info.Selections[f]; s != nil {
			if s.Kind() != types.MethodVal || types.IsInterface(s.Recv()) {
				return false
			}
			obj = s.Obj()
		} else {
			obj = info.Uses[f.Sel]
		}
	case *ast.IndexExpr, *ast.IndexListExpr: // an instance of a generic function
		return false
	}
	switch o := obj.(type) {
	case *types.Builtin:
		switch o.Name() {
		case "len", "cap", "min", "max", "real", "imag", "complex", "make", "new", "panic":
			return true
		}
		return false
	case *types.Func:
		if p := o.Pkg(); p != nil && (p.Path() == "math" || p.Path() == "math/bits") {
			return true
		}
		return e.pureFunc(o)
	}
	return false
}
