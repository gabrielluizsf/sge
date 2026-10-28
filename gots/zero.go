package gots

import (
	"fmt"
	"go/types"
)

// zero is an expression for a fresh zero value of t.
func (e *emitter) zero(t types.Type) string {
	switch u := t.Underlying().(type) {
	case *types.Basic:
		switch {
		case u.Info()&types.IsBoolean != 0:
			return "false"
		case u.Info()&types.IsString != 0:
			return `""`
		case u.Kind() == types.Uint64:
			return "0n"
		case u.Info()&types.IsComplex != 0:
			return "{ re: 0, im: 0 }"
		case u.Info()&types.IsNumeric != 0:
			return "0"
		}
		return "null"
	case *types.Struct:
		return "new " + e.structClass(t) + "()"
	case *types.Array:
		if b := backing(u.Elem()); b != "" {
			return fmt.Sprintf("new %s(%d)", b, u.Len())
		}
		return fmt.Sprintf("$.fillArray(%d, () => %s)", u.Len(), e.zero(u.Elem()))
	}
	return "null"
}

// zeroFn is a function making zero values of t, or "" when a typed array
// already holds them.
func (e *emitter) zeroFn(t types.Type) string {
	if backing(t) != "" {
		return ""
	}
	return "() => " + e.zero(t)
}
