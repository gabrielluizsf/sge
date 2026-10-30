package pose

import "math"

// Slerp is the rotation f of the way from a to b, the short way round.
// Rotations less than some 16 degrees apart, as an animation's keys
// mostly are, are mixed straight and made whole again, which strays from
// the arc by less than a hundredth of a degree, for a fraction of the
// trigonometry.
func Slerp(a, b [4]float64, f float64) [4]float64 {
	d := a[0]*b[0] + a[1]*b[1] + a[2]*b[2] + a[3]*b[3]
	if d < 0 {
		b, d = [4]float64{-b[0], -b[1], -b[2], -b[3]}, -d
	}
	wa, wb := 1-f, f
	if d >= 0.99 {
		q := [4]float64{wa*a[0] + wb*b[0], wa*a[1] + wb*b[1], wa*a[2] + wb*b[2], wa*a[3] + wb*b[3]}
		n := 1 / math.Sqrt(q[0]*q[0]+q[1]*q[1]+q[2]*q[2]+q[3]*q[3])
		return [4]float64{q[0] * n, q[1] * n, q[2] * n, q[3] * n}
	}
	th := math.Acos(d)
	wa, wb = math.Sin((1-f)*th)/math.Sin(th), math.Sin(f*th)/math.Sin(th)
	var q [4]float64
	var n float64
	for i := range q {
		q[i] = wa*a[i] + wb*b[i]
		n += q[i] * q[i]
	}
	n = math.Sqrt(n)
	for i := range q {
		q[i] /= n
	}
	return q
}
