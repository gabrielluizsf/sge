package gots

// intrinsics are functions whose calls are written as JavaScript's own.
//
// The exact ones give the very answers Go does: IEEE fixes them, and Go
// on js/wasm takes them from wasm's instructions too.
//
// The others are JavaScript's for speed: Go's own transcendental functions,
// transpiled, reduce their arguments in uint64, a BigInt here, and run five
// times slower than Math. Their special cases are C99's in both; a result
// may differ from Go's in its last bit. Those whose special cases Go and
// JavaScript disagree on (Pow, Max, Min, Hypot) are not here.
var intrinsics = map[string]string{
	// A compiler intrinsic with no body of its own: true is 1, false 0.
	"crypto/internal/constanttime.boolToUint8": "$Number",
	// Exact.
	"math.Sqrt":  "$Math.sqrt",
	"math.Floor": "$Math.floor",
	"math.Ceil":  "$Math.ceil",
	"math.Trunc": "$Math.trunc",
	"math.Abs":   "$Math.abs",
	// Within an ulp of Go's.
	"math.Sin":   "$Math.sin",
	"math.Cos":   "$Math.cos",
	"math.Tan":   "$Math.tan",
	"math.Asin":  "$Math.asin",
	"math.Acos":  "$Math.acos",
	"math.Atan":  "$Math.atan",
	"math.Atan2": "$Math.atan2",
	"math.Sinh":  "$Math.sinh",
	"math.Cosh":  "$Math.cosh",
	"math.Tanh":  "$Math.tanh",
	"math.Exp":   "$Math.exp",
	"math.Expm1": "$Math.expm1",
	"math.Log":   "$Math.log",
	"math.Log2":  "$Math.log2",
	"math.Log10": "$Math.log10",
	"math.Log1p": "$Math.log1p",
}
