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
	switch f := ast.Unparen(c.Fun).(type) {
	case *ast.Ident:
		fn, _ := info.Uses[f].(*types.Func)
		return fn
	case *ast.IndexExpr: // an instantiated generic function
		if id, ok := ast.Unparen(f.X).(*ast.Ident); ok {
			fn, _ := info.Uses[id].(*types.Func)
			return fn
		}
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
