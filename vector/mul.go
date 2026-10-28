package vector

// Mul is the component-wise product of v and o.
func (v Vector3) Mul(o Vector3) Vector3 {
	return Vector3{v.X * o.X, v.Y * o.Y, v.Z * o.Z}
}
