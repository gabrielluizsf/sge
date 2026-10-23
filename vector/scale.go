package vector

// Scale multiplies every component of v by s.
func (v Vector3) Scale(s float64) Vector3 {
	return Vector3{v.X * s, v.Y * s, v.Z * s}
}
