package gots

import (
	"go/ast"
	"go/types"
)

// readonlyParam reports whether fn only reads its i-th parameter: never
// changes it, takes its address, slices it, or lets a closure keep it.
// Storing or returning it copies it anyway.
func (e *emitter) readonlyParam(fn *types.Func, i int) bool {
	d, ok := e.funcDeclOf(fn)
	if !ok {
		return false
	}
	k := 0
	for _, field := range d.decl.Type.Params.List {
		names := field.Names
		if len(names) == 0 {
			names = []*ast.Ident{nil}
		}
		for _, id := range names {
			if k == i {
				if id == nil || id.Name == "_" {
					return true
				}
				obj := d.info.Defs[id]
				return obj != nil && !mutates(d.info, d.decl.Body, obj) && !pinned(d.info, d.decl.Body, obj)
			}
			k++
		}
	}
	return false
}
