package gots

import "go/types"

// keyExpr is the map key of js, a value of type t, with the equality Go
// gives it: a struct by its fields, a pointer by identity, an interface
// by its dynamic type and value.
func (e *emitter) keyExpr(js string, t types.Type) string {
	switch u := t.Underlying().(type) {
	case *types.Struct:
		return js + ".$key()"
	case *types.Array:
		return "$.keyOfType(" + e.desc(t) + ", " + js + ")"
	case *types.Pointer, *types.Chan, *types.Map, *types.Signature:
		return "$.keyOfPtr(" + js + ")"
	case *types.Basic:
		if u.Kind() == types.UnsafePointer || u.Info()&types.IsComplex != 0 {
			return "$.keyOfType(" + e.desc(t) + ", " + js + ")"
		}
	}
	return "$.keyOf(" + js + ")"
}
