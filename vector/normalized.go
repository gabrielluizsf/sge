package vector

// Normalized is v scaled to length 1, or the zero vector when v has no length
// to scale.
func (v Vector3) Normalized() Vector3 {
	l := v.Length()
	if l == 0 {
		return Vector3{}
	}
	return v.Scale(1 / l)
}
