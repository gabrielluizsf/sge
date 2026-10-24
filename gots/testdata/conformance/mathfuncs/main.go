// math as Go has it, special cases included, where JavaScript's Math
// answers otherwise: Max and Min with NaN and infinities, Pow's edge
// cases, and the functions Math has no counterpart for.
package main

import (
	"fmt"
	"math"
)

func main() {
	nan, inf := math.NaN(), math.Inf(1)
	fmt.Printf("%.12g %.12g %.12g\n", math.Atan2(1, -1), math.Exp(1), math.Log(10))
	fmt.Println(math.Max(inf, nan), math.Min(math.Inf(-1), nan), math.Max(0, math.Copysign(0, -1)), math.Min(0, math.Copysign(0, -1)))
	fmt.Println(math.Pow(1, nan), math.Pow(nan, 0), math.Pow(-8, 1.0/3), math.Pow(2, 10), math.Pow(-2, 3))
	fmt.Println(math.Nextafter(1, 2)-1, math.Nextafter32(1, 2)-1, math.Nextafter(0, -1))
	fr, exp := math.Frexp(48)
	fmt.Println(fr, exp, math.Ldexp(0.75, 6), math.Ilogb(48), math.Logb(48))
	fmt.Printf("%.12g %.12g %.12g %.12g\n", math.Gamma(5), math.Erf(0.5), math.Erfinv(0.5), math.J0(1))
	lg, sign := math.Lgamma(-0.5)
	fmt.Printf("%.12g %d %.12g %.12g\n", lg, sign, math.Asinh(1), math.Atanh(0.5))
	fmt.Println(math.Sqrt(2), math.Floor(-1.5), math.Ceil(-1.5), math.Trunc(-1.5), math.Abs(-0.0), math.Round(-2.5), math.RoundToEven(2.5))
	fmt.Println(math.Mod(7, -3), math.Remainder(7, 3), math.Dim(3, 5), math.FMA(2, 3, 4), math.Cbrt(27))
	fmt.Println(math.Float64bits(1.5), math.Float64frombits(0x4004000000000000), math.Float32bits(1.5), math.Signbit(math.Copysign(0, -1)))
	fmt.Println(math.MaxInt32, math.MinInt8, math.MaxUint32, math.SmallestNonzeroFloat64, math.Ln10, math.Log2E)
	f := math.Sqrt
	fmt.Println(f(9), math.IsInf(1/math.Copysign(0, -1), -1))
}
