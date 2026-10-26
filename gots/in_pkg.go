package gots

import "go/types"

// inPkg is desc, the descriptor of the named type t, saying the package
// t is declared in, as reflect's PkgPath and gob's names for registered
// types need.
func (e *emitter) inPkg(t *types.Named, desc string) string {
	pkg := t.Obj().Pkg()
	if pkg == nil {
		return desc
	}
	path := pkg.Path()
	if pkg.Name() == "main" {
		path = "main" // as Go's reflect has it, whatever the command's path
	}
	return "$.inPkg(" + desc + ", " + jsString(path) + ")"
}
