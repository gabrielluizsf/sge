package human

import "math"

// smooth is the value at t, from 0 to 1 round a loop, of the smooth
// periodic curve through v's values, evenly spaced: a Catmull-Rom spline,
// which passes through each and does not overshoot a gentle run of them.
func smooth(v []float64, t float64) float64 {
	n := len(v)
	x := (t - math.Floor(t)) * float64(n)
	i := int(x)
	f := x - float64(i)
	p0, p1, p2, p3 := v[(i+n-1)%n], v[i%n], v[(i+1)%n], v[(i+2)%n]
	return p1 + 0.5*f*(p2-p0+f*(2*p0-5*p1+4*p2-p3+f*(3*(p1-p2)+p3-p0)))
}
