package gots

import (
	"go/ast"
	"go/types"
)

// pinned reports whether something besides the code of body itself can
// reach the variable obj while body runs: its address is taken, a pointer
// method is called on it, an array in it is sliced, or a closure refers
// to it.
func pinned(info *types.Info, body *ast.BlockStmt, obj types.Object) bool {
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.UnaryExpr:
			found = n.Op.String() == "&" && rootVar(info, n.X) == obj
		case *ast.SliceExpr:
			_, isArr := info.TypeOf(n.X).Underlying().(*types.Array)
			found = isArr && rootVar(info, n.X) == obj
		case *ast.CallExpr:
			if sel, ok := ast.Unparen(n.Fun).(*ast.SelectorExpr); ok {
				if s := info.Selections[sel]; s != nil && s.Kind() == types.MethodVal {
					if _, ptr := s.Obj().Type().(*types.Signature).Recv().Type().(*types.Pointer); ptr {
						found = rootVar(info, sel.X) == obj
					}
				}
			}
		case *ast.SelectorExpr:
			// A method value of a pointer method keeps the address.
			if s := info.Selections[n]; s != nil && s.Kind() == types.MethodVal {
				if _, ptr := s.Obj().Type().(*types.Signature).Recv().Type().(*types.Pointer); ptr {
					found = rootVar(info, n.X) == obj
				}
			}
		case *ast.FuncLit:
			ast.Inspect(n.Body, func(m ast.Node) bool {
				if id, ok := m.(*ast.Ident); ok && info.Uses[id] == obj {
					found = true
				}
				return !found
			})
		}
		return !found
	})
	return found
}

// rootVar is the variable x is stored in, through fields and array
// elements, or nil if x is not stored in a variable.
func rootVar(info *types.Info, x ast.Expr) types.Object {
	for {
		switch y := ast.Unparen(x).(type) {
		case *ast.Ident:
			return info.ObjectOf(y)
		case *ast.SelectorExpr:
			s := info.Selections[y]
			if s == nil || isPointer(s.Recv()) {
				return nil
			}
			x = y.X
		case *ast.IndexExpr:
			if _, isArr := info.TypeOf(y.X).Underlying().(*types.Array); !isArr {
				return nil
			}
			x = y.X
		default:
			return nil
		}
	}
}
