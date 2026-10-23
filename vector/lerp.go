package vector

// Lerp moves from v towards o by t, where 0 is v and 1 is o.
func (v Vector3) Lerp(o Vector3, t float64) Vector3 {
	return v.Add(o.Sub(v).Scale(t))
}
