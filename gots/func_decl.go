package gots

import (
	"go/ast"
	"go/types"
)

// funcDecl is a function's declaration and the type information it was
// checked with.
type funcDecl struct {
	decl *ast.FuncDecl
	info *types.Info
}

// declOf is the declaration of fn, or of the generic function it is an
// instance of, among the packages gots writes; ok is false for one it
// does not have the body of.
func (p *Program) declOf(fn *types.Func) (funcDecl, bool) {
	if p.decls == nil {
		p.decls = map[*types.Func]funcDecl{}
		for _, pkg := range p.Packages {
			if pkg.Info == nil || pkg.Std && !pkg.Transpile {
				continue
			}
			for _, f := range pkg.Files {
				for _, d := range f.Decls {
					if fd, ok := d.(*ast.FuncDecl); ok && fd.Body != nil {
						if obj, ok := pkg.Info.Defs[fd.Name].(*types.Func); ok {
							p.decls[obj] = funcDecl{fd, pkg.Info}
						}
					}
				}
			}
		}
	}
	d, ok := p.decls[fn.Origin()]
	return d, ok
}
