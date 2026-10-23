package vector

// Mul composes two rotations: the result rotates by o first, then by q.
func (q Quaternion) Mul(o Quaternion) Quaternion {
	return Quaternion{
		q.W*o.X + q.X*o.W + q.Y*o.Z - q.Z*o.Y,
		q.W*o.Y - q.X*o.Z + q.Y*o.W + q.Z*o.X,
		q.W*o.Z + q.X*o.Y - q.Y*o.X + q.Z*o.W,
		q.W*o.W - q.X*o.X - q.Y*o.Y - q.Z*o.Z,
	}
}
