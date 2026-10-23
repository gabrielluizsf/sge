package vector

// Cross is the cross product of v and o, perpendicular to both.
func (v Vector3) Cross(o Vector3) Vector3 {
	return Vector3{
		v.Y*o.Z - v.Z*o.Y,
		v.Z*o.X - v.X*o.Z,
		v.X*o.Y - v.Y*o.X,
	}
}
