package gots

import "go/types"

// textField reports whether field i of st is a string TypeScript code
// reads and writes as text: an exported string field of an exported
// struct of a public package. Go code keeps the bytes in name$s, and the
// class gives TypeScript a name property converting both ways.
func (e *emitter) textField(st *types.Struct, i int) bool {
	f := st.Field(i)
	return f.Exported() && !f.Embedded() && isString(f.Type()) && e.publicStruct(st)
}

// publicStruct reports whether st is the struct of an exported type of a
// public package.
func (e *emitter) publicStruct(st *types.Struct) bool {
	if e.publicStructs == nil {
		e.publicStructs = map[*types.Struct]bool{}
		mark := func(pkg *types.Package) {
			scope := pkg.Scope()
			for _, n := range scope.Names() {
				if tn, ok := scope.Lookup(n).(*types.TypeName); ok && tn.Exported() && !tn.IsAlias() {
					if s, ok := tn.Type().Underlying().(*types.Struct); ok {
						e.publicStructs[s] = true
					}
				}
			}
		}
		for path := range e.prog.Public {
			if p := e.prog.Packages[path]; p != nil && p.Types != nil {
				mark(p.Types)
			}
		}
		if e.prog.Public[e.tpkg.Path()] {
			mark(e.tpkg)
		}
	}
	return e.publicStructs[st]
}

// fieldJS is the property field i of st is kept in.
func (e *emitter) fieldJS(st *types.Struct, i int) string {
	if e.textField(st, i) {
		return fieldName(st.Field(i).Name()) + "$s"
	}
	return fieldJS(st, i)
}
