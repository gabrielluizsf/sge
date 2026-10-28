package vector

// SqrLength is the squared length of v, cheaper than [Vector3.Length] when
// only comparing distances.
func (v Vector3) SqrLength() float64 { return v.Dot(v) }
