//go:build gots_overlay

// Overlay for the floor functions js/wasm has in assembly: Go's own pure
// versions. Calls of Floor, Ceil and Trunc are JavaScript's Math anyway
// (see intrinsics.go), which gives the same answers.

package math

func archFloor(x float64) float64 { return floor(x) }

func archCeil(x float64) float64 { return ceil(x) }

func archTrunc(x float64) float64 { return trunc(x) }
