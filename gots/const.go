package gots

import (
	"go/constant"
	"go/types"
	"math"
	"strconv"
)

// constant is the literal for the constant v of type t.
func (e *emitter) constant(v constant.Value, t types.Type) string {
	if isInterface(t) {
		// A constant converted to an interface type keeps its default type.
		return e.box(e.constant(v, types.Default(constType(v))), types.Default(constType(v)))
	}
	k := defaultKind(basicKind(t))
	switch {
	case v.Kind() == constant.Bool:
		return strconv.FormatBool(constant.BoolVal(v))
	case v.Kind() == constant.String:
		return jsString(constant.StringVal(v))
	case k == types.Uint64:
		return constant.ToInt(v).ExactString() + "n"
	case k == types.Complex64 || k == types.Complex128:
		re, _ := constant.Float64Val(constant.Real(v))
		im, _ := constant.Float64Val(constant.Imag(v))
		return "{ re: " + floatLit(re) + ", im: " + floatLit(im) + " }"
	case isInteger(t) || v.Kind() == constant.Int && !isFloat(t):
		if i, ok := constant.Int64Val(constant.ToInt(v)); ok {
			if i < 0 {
				return "(" + strconv.FormatInt(i, 10) + ")"
			}
			return strconv.FormatInt(i, 10)
		}
		if u, ok := constant.Uint64Val(constant.ToInt(v)); ok {
			return strconv.FormatUint(u, 10)
		}
	}
	f, _ := constant.Float64Val(constant.ToFloat(v))
	if k == types.Float32 {
		f = float64(float32(f))
	}
	return floatLit(f)
}

// constType is the untyped type a constant value has before it is given one.
func constType(v constant.Value) types.Type {
	switch v.Kind() {
	case constant.Bool:
		return types.Typ[types.UntypedBool]
	case constant.String:
		return types.Typ[types.UntypedString]
	case constant.Int:
		return types.Typ[types.UntypedInt]
	case constant.Float:
		return types.Typ[types.UntypedFloat]
	}
	return types.Typ[types.UntypedComplex]
}

func floatLit(f float64) string {
	switch {
	case math.IsInf(f, 1):
		return "Infinity"
	case math.IsInf(f, -1):
		return "(-Infinity)"
	case math.IsNaN(f):
		return "NaN"
	case f < 0 || f == 0 && math.Signbit(f):
		return "(" + strconv.FormatFloat(f, 'g', -1, 64) + ")"
	}
	return strconv.FormatFloat(f, 'g', -1, 64)
}
