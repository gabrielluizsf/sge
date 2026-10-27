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
	// Exact: see rt/bits.ts. uint is 32 bits.
	"math/bits.Add32":                                      "$.bits.add32",
	"math/bits.Add":                                        "$.bits.add32",
	"math/bits.Sub32":                                      "$.bits.sub32",
	"math/bits.Sub":                                        "$.bits.sub32",
	"math/bits.Mul32":                                      "$.bits.mul32",
	"math/bits.Mul":                                        "$.bits.mul32",
	"math/bits.RotateLeft32":                               "$.bits.rotl32",
	"math/bits.RotateLeft":                                 "$.bits.rotl32",
	"math/bits.LeadingZeros32":                             "$Math.clz32",
	"math/bits.LeadingZeros":                               "$Math.clz32",
	"math/bits.Len32":                                      "$.bits.len32",
	"math/bits.Len":                                        "$.bits.len32",
	"math/bits.TrailingZeros32":                            "$.bits.trailingZeros32",
	"math/bits.TrailingZeros":                              "$.bits.trailingZeros32",
	"math/bits.OnesCount32":                                "$.bits.onesCount32",
	"math/bits.OnesCount":                                  "$.bits.onesCount32",
	"crypto/internal/fips140/aes/gcm.ghashAbsorb":          "$.ghashAbsorb",
	"crypto/internal/fips140/aes/gcm.aeadSeal":             "$.aeadSeal",
	"crypto/internal/fips140/aes/gcm.aeadOpen":             "$.aeadOpen",
	"vendor/golang.org/x/crypto/chacha20poly1305.aeadSeal": "$.aeadSeal",
	"vendor/golang.org/x/crypto/chacha20poly1305.aeadOpen": "$.aeadOpen",
	"golang.org/x/crypto/chacha20poly1305.aeadSeal":        "$.aeadSeal",
	"golang.org/x/crypto/chacha20poly1305.aeadOpen":        "$.aeadOpen",
	"crypto/ecdh.x25519Native":                             "$.x25519",
	"crypto/internal/fips140/sha3.keccakF1600Generic":      "$.keccakF1600",
	"crypto/internal/fips140/sha512.blockGeneric":          "$.sha512Block",
	"crypto/internal/fips140/bigmod.addMulVVW":             "$.bits.addMulVVW",
	"crypto/internal/fips140/bigmod.montgomeryWords":       "$.bits.montgomeryWords",
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
