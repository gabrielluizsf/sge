package vector

// Quaternion is a rotation in 3D space. The zero value is not a rotation; use
// [Identity] for "not rotated".
type Quaternion struct {
	X, Y, Z, W float64
}
