package gots

import "go/types"

func isStruct(t types.Type) bool {
	_, ok := t.Underlying().(*types.Struct)
	return ok
}

func isArray(t types.Type) bool {
	_, ok := t.Underlying().(*types.Array)
	return ok
}

// isValue reports whether values of t are JavaScript objects with storage
// of their own, copied by $clone and assigned by $set.
func isValue(t types.Type) bool { return isStruct(t) || isArray(t) }

func isInterface(t types.Type) bool { return types.IsInterface(t) }

// isObjectPtr reports whether t is a pointer whose pointee is an object, so
// the pointer is the object itself.
func isObjectPtr(t types.Type) bool {
	p, ok := t.Underlying().(*types.Pointer)
	return ok && isValue(p.Elem())
}

func basicKind(t types.Type) types.BasicKind {
	if b, ok := t.Underlying().(*types.Basic); ok {
		return b.Kind()
	}
	return types.Invalid
}

func isString(t types.Type) bool {
	b, ok := t.Underlying().(*types.Basic)
	return ok && b.Info()&types.IsString != 0
}

func isInteger(t types.Type) bool {
	b, ok := t.Underlying().(*types.Basic)
	return ok && b.Info()&types.IsInteger != 0
}

func isFloat(t types.Type) bool {
	b, ok := t.Underlying().(*types.Basic)
	return ok && b.Info()&types.IsFloat != 0
}

func isNumeric(t types.Type) bool {
	b, ok := t.Underlying().(*types.Basic)
	return ok && b.Info()&types.IsNumeric != 0
}

// isBig reports whether t is kept in a bigint.
func isBig(t types.Type) bool { return basicKind(t) == types.Uint64 }

// defaultKind is the kind an untyped constant takes when nothing else
// decides it.
func defaultKind(k types.BasicKind) types.BasicKind {
	switch k {
	case types.UntypedBool:
		return types.Bool
	case types.UntypedInt:
		return types.Int
	case types.UntypedRune:
		return types.Int32
	case types.UntypedFloat:
		return types.Float64
	case types.UntypedComplex:
		return types.Complex128
	case types.UntypedString:
		return types.String
	}
	return k
}

// backing is the typed array that holds elements of t, or "" for a plain
// array.
func backing(t types.Type) string {
	switch defaultKind(basicKind(t)) {
	case types.Int, types.Int64, types.Uint, types.Uintptr, types.Float64:
		return "Float64Array"
	case types.Float32:
		return "Float32Array"
	case types.Int32:
		return "Int32Array"
	case types.Int16:
		return "Int16Array"
	case types.Int8:
		return "Int8Array"
	case types.Uint32:
		return "Uint32Array"
	case types.Uint16:
		return "Uint16Array"
	case types.Uint8:
		return "Uint8Array"
	case types.Uint64:
		return "BigUint64Array"
	}
	return ""
}

// backingOrNull is backing as an expression for the runtime: the
// constructor, or null for a plain array.
func backingOrNull(t types.Type) string {
	if b := backing(t); b != "" {
		return b
	}
	return "null"
}

// structOf is the struct underneath t, looking through a pointer.
func structOf(t types.Type) *types.Struct {
	if p, ok := t.Underlying().(*types.Pointer); ok {
		t = p.Elem()
	}
	s, _ := t.Underlying().(*types.Struct)
	return s
}

// derefType is t's pointee when t is a pointer.
func derefType(t types.Type) types.Type {
	if p, ok := t.Underlying().(*types.Pointer); ok {
		return p.Elem()
	}
	return t
}
