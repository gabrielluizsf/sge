package vector

// Add is the component-wise sum of v and o.
func (v Vector3) Add(o Vector3) Vector3 {
	return Vector3{v.X + o.X, v.Y + o.Y, v.Z + o.Z}
}
