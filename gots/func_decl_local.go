package gots

import (
	"go/ast"
	"go/types"
)

// funcDeclOf is the declaration of fn: among the module's own files first,
// which a package written with its tests was checked again with, then
// among every package gots writes.
func (e *emitter) funcDeclOf(fn *types.Func) (funcDecl, bool) {
	fn = fn.Origin()
	if e.localDecls == nil {
		e.localDecls = map[*types.Func]funcDecl{}
		for _, f := range e.files {
			for _, d := range f.Decls {
				if fd, ok := d.(*ast.FuncDecl); ok && fd.Body != nil {
					if obj, ok := e.info0.Defs[fd.Name].(*types.Func); ok {
						e.localDecls[obj] = funcDecl{fd, e.info0}
					}
				}
			}
		}
	}
	if d, ok := e.localDecls[fn]; ok {
		return d, true
	}
	return e.prog.declOf(fn)
}
