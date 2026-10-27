package sound

import "math"

// gains is pan's left and right gains for a sound of gain: equal power,
// so a sound moving across is as loud in the middle as at either side.
func gains(gain, pan float64) (left, right float32) {
	pan = math.Max(-1, math.Min(1, pan))
	a := (pan + 1) * math.Pi / 4
	return float32(gain * math.Cos(a) * math.Sqrt2), float32(gain * math.Sin(a) * math.Sqrt2)
}
