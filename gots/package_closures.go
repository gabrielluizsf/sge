package gots

import (
	"go/ast"
	"go/types"
)

// addPackageClosures reads the closures a package's variables are
// initialized with, var hook = func(...) {...}, each a function of its own
// and a function value: net's lookup hooks call what they are given, which
// may wait.
func (b *blocking) addPackageClosures(gd *ast.GenDecl, info *types.Info) {
	for _, spec := range gd.Specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		for _, v := range vs.Values {
			ast.Inspect(v, func(n ast.Node) bool {
				lit, ok := n.(*ast.FuncLit)
				if !ok {
					return true
				}
				b.addFunc(lit.Type, lit.Body, info)
				b.funcValue(b.key(lit.Type.Pos()), info.TypeOf(lit))
				return false
			})
		}
	}
}
