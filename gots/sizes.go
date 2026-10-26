package gots

import "go/types"

// wasmSizes are the sizes packages are checked with: js/wasm's, but for
// uint, which is 32 bits as on 386 and arm. A JavaScript number holds 53
// bits exactly, so a 64-bit uint would lose its high bits; at 32 the
// standard library takes the paths it has for 32-bit machines (math/big's
// and crypto's words among them), whose arithmetic numbers do exactly.
// int stays 64 bits, a number, as gots keeps it, and so does uintptr,
// which the runtime has as wide as a pointer.
var wasmSizes types.Sizes = sizes{types.SizesFor("gc", "wasm")}

type sizes struct{ wasm types.Sizes }

// word reports whether t is uint.
func word(t types.Type) bool { return basicKind(t) == types.Uint }

func (s sizes) Alignof(t types.Type) int64 {
	switch u := t.Underlying().(type) {
	case *types.Array:
		return s.Alignof(u.Elem())
	case *types.Struct:
		a := int64(1)
		for i := range u.NumFields() {
			a = max(a, s.Alignof(u.Field(i).Type()))
		}
		return a
	}
	if word(t) {
		return 4
	}
	return s.wasm.Alignof(t)
}

func (s sizes) Offsetsof(fields []*types.Var) []int64 {
	offs := make([]int64, len(fields))
	var o int64
	for i, f := range fields {
		o = align(o, s.Alignof(f.Type()))
		offs[i] = o
		o += s.Sizeof(f.Type())
	}
	return offs
}

func (s sizes) Sizeof(t types.Type) int64 {
	switch u := t.Underlying().(type) {
	case *types.Array:
		if u.Len() == 0 {
			return 0
		}
		return align(s.Sizeof(u.Elem()), s.Alignof(u.Elem())) * u.Len()
	case *types.Struct:
		n := u.NumFields()
		if n == 0 {
			return 0
		}
		fields := make([]*types.Var, n)
		for i := range n {
			fields[i] = u.Field(i)
		}
		offs := s.Offsetsof(fields)
		size := offs[n-1] + s.Sizeof(fields[n-1].Type())
		if size > 0 && s.Sizeof(fields[n-1].Type()) == 0 {
			// As gc pads a final zero-sized field, so a pointer to it
			// stays inside the struct.
			size++
		}
		return align(size, s.Alignof(t))
	}
	if word(t) {
		return 4
	}
	return s.wasm.Sizeof(t)
}

func align(x, a int64) int64 { return (x + a - 1) / a * a }
