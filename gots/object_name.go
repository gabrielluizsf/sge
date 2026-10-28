package gots

import "go/types"

// importAlias is the module alias p is imported under, importing it on
// first use.
func (e *emitter) importAlias(p *types.Package) string {
	if alias, ok := e.imports[p.Path()]; ok {
		return alias
	}
	alias := mangle(p.Name())
	for i := 1; e.pkgNames[alias]; i++ {
		alias = mangle(p.Name()) + "$" + string(rune('0'+i))
	}
	e.pkgNames[alias] = true
	e.imports[p.Path()] = alias
	e.importOrder = append(e.importOrder, p.Path())
	return alias
}

// isPackageLevel reports whether obj is declared at package scope.
func isPackageLevel(obj types.Object) bool {
	return obj.Pkg() != nil && obj.Parent() == obj.Pkg().Scope()
}

// objName is how emitted code refers to a package-level object or a local.
func (e *emitter) objName(obj types.Object) string {
	if obj.Pkg() != nil && obj.Pkg() != e.tpkg && isPackageLevel(obj) {
		return e.importAlias(obj.Pkg()) + "." + mangle(obj.Name())
	}
	if tn, ok := obj.(*types.TypeName); ok {
		if n, ok := e.localTypes[tn]; ok {
			return n
		}
	}
	if isPackageLevel(obj) {
		return mangle(obj.Name())
	}
	if n, ok := e.names[obj]; ok {
		return n
	}
	return e.declare(obj)
}

// declare hands a local its binding name.
func (e *emitter) declare(obj types.Object) string {
	var n string
	if e.sc == nil {
		n = mangle(obj.Name())
	} else {
		n = e.sc.fresh(obj.Name())
	}
	e.names[obj] = n
	return n
}

// varRef is how a variable's value is read and written: through its Cell
// when its address is taken.
func (e *emitter) varRef(obj types.Object) string {
	n := e.objName(obj)
	if e.boxed[obj] {
		return n + ".v"
	}
	return n
}
