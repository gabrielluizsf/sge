package human

// curve is a value through keys at the times at, smooth at each, which goes
// beyond none of them: a run of equal keys stays exactly still.
type curve struct {
	at, keys []float64
}

// of is the curve's value at t, within its times.
func (c curve) of(t float64) float64 {
	k := 0
	for k < len(c.at)-2 && t >= c.at[k+1] {
		k++
	}
	h := c.at[k+1] - c.at[k]
	u := (t - c.at[k]) / h
	v0, v1 := c.keys[k], c.keys[k+1]
	m0, m1 := c.slope(k)*h, c.slope(k+1)*h
	return (2*u*u*u-3*u*u+1)*v0 + (u*u*u-2*u*u+u)*m0 + (-2*u*u*u+3*u*u)*v1 + (u*u*u-u*u)*m1
}

// slope is the curve's slope at key k: none at either end, at a turning point
// or by a still run, and otherwise the harmonic mean of the two sides, which
// keeps it from overshooting.
func (c curve) slope(k int) float64 {
	if k == 0 || k == len(c.keys)-1 {
		return 0
	}
	a := (c.keys[k] - c.keys[k-1]) / (c.at[k] - c.at[k-1])
	b := (c.keys[k+1] - c.keys[k]) / (c.at[k+1] - c.at[k])
	if a*b <= 0 {
		return 0
	}
	return 2 * a * b / (a + b)
}
