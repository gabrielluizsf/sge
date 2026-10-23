package gots

import (
	"go/ast"
	"go/types"
)

// staticCallee is the function or method c calls, when that is known
// without running it: not a function value, not a method of an interface.
func (e *emitter) staticCallee(c *ast.CallExpr) *types.Func {
	switch f := ast.Unparen(c.Fun).(type) {
	case *ast.Ident:
		fn, _ := e.info.Uses[f].(*types.Func)
		return fn
	case *ast.IndexExpr: // an instantiated generic function
		if id, ok := ast.Unparen(f.X).(*ast.Ident); ok {
			fn, _ := e.info.Uses[id].(*types.Func)
			return fn
		}
	case *ast.SelectorExpr:
		if s := e.info.Selections[f]; s != nil {
			fn, _ := s.Obj().(*types.Func)
			if s.Kind() != types.MethodVal || fn == nil || isInterface(e.recvType(fn)) {
				return nil
			}
			return fn
		}
		fn, _ := e.info.Uses[f.Sel].(*types.Func)
		return fn
	}
	return nil
}
