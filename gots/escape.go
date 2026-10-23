package gots

import (
	"go/ast"
	"go/types"
)

// findBoxed marks the local variables that must live in a Cell: those
// whose address is taken, explicitly with & or implicitly by calling a
// pointer method on them, unless they are structs or arrays, which are
// objects and pointed at directly.
func (e *emitter) findBoxed(files []*ast.File) { e.findBoxedIn(files, e.info) }

// findBoxedIn is findBoxed over another package's files, whose generic
// code this module may instantiate.
func (e *emitter) findBoxedIn(files []*ast.File, info *types.Info) {
	saved := e.info
	e.info = info
	defer func() { e.info = saved }()
	for _, f := range files {
		ast.Inspect(f, func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.UnaryExpr:
				if n.Op.String() == "&" {
					e.markBoxed(n.X)
				}
			case *ast.CallExpr:
				sel, ok := ast.Unparen(n.Fun).(*ast.SelectorExpr)
				if !ok {
					break
				}
				s := e.info.Selections[sel]
				if s == nil || s.Kind() != types.MethodVal {
					break
				}
				if _, ptrRecv := s.Obj().Type().(*types.Signature).Recv().Type().(*types.Pointer); ptrRecv && len(s.Index()) == 1 {
					if _, isPtr := s.Recv().Underlying().(*types.Pointer); !isPtr {
						e.markBoxed(sel.X)
					}
				}
			}
			return true
		})
	}
}

func (e *emitter) markBoxed(x ast.Expr) {
	id, ok := ast.Unparen(x).(*ast.Ident)
	if !ok {
		return
	}
	v, ok := e.info.ObjectOf(id).(*types.Var)
	if !ok || isPackageLevel(v) || isValue(e.varType(v)) {
		return
	}
	e.boxed[v] = true
}
