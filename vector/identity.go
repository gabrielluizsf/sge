package vector

// Identity is the rotation that leaves everything where it is.
func Identity() Quaternion { return Quaternion{W: 1} }
