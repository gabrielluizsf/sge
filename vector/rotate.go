package vector

// Rotate turns v by the rotation q.
func (q Quaternion) Rotate(v Vector3) Vector3 {
	u := Vector3{q.X, q.Y, q.Z}
	t := u.Cross(v).Scale(2)
	return v.Add(t.Scale(q.W)).Add(u.Cross(t))
}
