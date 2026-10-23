package vector

// Sub is the component-wise difference of v and o.
func (v Vector3) Sub(o Vector3) Vector3 {
	return Vector3{v.X - o.X, v.Y - o.Y, v.Z - o.Z}
}
