package gots

import "go/types"

// small reports whether integers of kind k fit JavaScript's 32-bit
// operators, which then do Go's wrapping for free.
func small(k types.BasicKind) bool {
	switch k {
	case types.Int8, types.Int16, types.Int32, types.Uint8, types.Uint16, types.Uint32, types.Uint:
		return true
	}
	return false
}

func unsigned(k types.BasicKind) bool {
	switch k {
	case types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64, types.Uintptr:
		return true
	}
	return false
}

func width(k types.BasicKind) int {
	switch k {
	case types.Int8, types.Uint8:
		return 8
	case types.Int16, types.Uint16:
		return 16
	case types.Int32, types.Uint32, types.Uint, types.Float32:
		return 32
	}
	return 64
}

// wrap brings the result of an operation back into the range of kind k.
func wrap(k types.BasicKind, js string) string {
	switch k {
	case types.Int8:
		return "((" + js + ") << 24 >> 24)"
	case types.Int16:
		return "((" + js + ") << 16 >> 16)"
	case types.Int32:
		return "((" + js + ") | 0)"
	case types.Uint8:
		return "((" + js + ") & 255)"
	case types.Uint16:
		return "((" + js + ") & 65535)"
	case types.Uint32, types.Uint:
		return "((" + js + ") >>> 0)"
	case types.Uint64:
		return "$.int.u64(" + js + ")"
	case types.Float32:
		return "$Math.fround(" + js + ")"
	}
	return js
}
