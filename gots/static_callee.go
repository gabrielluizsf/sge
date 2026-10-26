package gots

import (
	"go/ast"
	"go/types"
)

// staticCallee is the function or method c calls, when that is known
// without running it: not a function value, not a method of an interface.
func (e *emitter) staticCallee(c *ast.CallExpr) *types.Func {
	fn := staticCalleeIn(e.info, c)
	if fn != nil && fn.Type().(*types.Signature).Recv() != nil && isInterface(e.recvType(fn)) {
		return nil
	}
	return fn
}

// staticCalleeIn is staticCallee with info's types, before substitution.
func staticCalleeIn(info *types.Info, c *ast.CallExpr) *types.Func {
	fun := ast.Unparen(c.Fun)
	// A generic function with its type arguments written out, here or
	// in another package.
	switch f := fun.(type) {
	case *ast.IndexExpr:
		if _, ok := info.Instances[instanceIdent(f.X)]; ok {
			fun = ast.Unparen(f.X)
		}
	case *ast.IndexListExpr:
		fun = ast.Unparen(f.X)
	}
	switch f := fun.(type) {
	case *ast.Ident:
		fn, _ := info.Uses[f].(*types.Func)
		return fn
	case *ast.SelectorExpr:
		if s := info.Selections[f]; s != nil {
			fn, _ := s.Obj().(*types.Func)
			if s.Kind() != types.MethodVal || fn == nil || types.IsInterface(fn.Type().(*types.Signature).Recv().Type()) {
				return nil
			}
			return fn
		}
		fn, _ := info.Uses[f.Sel].(*types.Func)
		return fn
	}
	return nil
}

// instanceIdent is the name x calls a function by, F or pkg.F, or nil.
func instanceIdent(x ast.Expr) *ast.Ident {
	switch x := ast.Unparen(x).(type) {
	case *ast.Ident:
		return x
	case *ast.SelectorExpr:
		return x.Sel
	}
	return nil
}
