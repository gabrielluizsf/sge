package gots

import (
	"go/ast"
	"go/types"
)

// closureVars are the local variables of files that hold one closure and
// nothing else: assigned once, from a func literal, their address never
// taken. A call through one is a call of that closure, known as well as a
// call of a declared function, f := func() {...}; f() being how Go code
// names a helper.
func closureVars(files []*ast.File, info *types.Info) map[types.Object]*ast.FuncLit {
	count := map[types.Object]int{}
	lits := map[types.Object]*ast.FuncLit{}
	assign := func(id *ast.Ident, rhs ast.Expr) {
		obj := info.ObjectOf(id)
		v, ok := obj.(*types.Var)
		if !ok || v.Parent() == nil || v.Parent() == v.Pkg().Scope() {
			return
		}
		count[obj]++
		if lit, ok := ast.Unparen(rhs).(*ast.FuncLit); ok && rhs != nil {
			lits[obj] = lit
		}
	}
	for _, f := range files {
		ast.Inspect(f, func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.AssignStmt:
				for i, l := range n.Lhs {
					if id, ok := ast.Unparen(l).(*ast.Ident); ok {
						var rhs ast.Expr
						if len(n.Lhs) == len(n.Rhs) {
							rhs = n.Rhs[i]
						}
						assign(id, rhs)
					}
				}
			case *ast.ValueSpec:
				for i, id := range n.Names {
					if i < len(n.Values) && len(n.Names) == len(n.Values) {
						assign(id, n.Values[i])
					}
				}
			case *ast.RangeStmt:
				for _, x := range []ast.Expr{n.Key, n.Value} {
					if id, ok := x.(*ast.Ident); ok {
						assign(id, nil)
					}
				}
			case *ast.UnaryExpr:
				if id, ok := ast.Unparen(n.X).(*ast.Ident); ok && n.Op.String() == "&" {
					count[info.ObjectOf(id)] += 2
				}
			}
			return true
		})
	}
	out := map[types.Object]*ast.FuncLit{}
	for obj, lit := range lits {
		if count[obj] == 1 {
			out[obj] = lit
		}
	}
	return out
}
