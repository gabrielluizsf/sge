package vector

import "testing"

var sinkV Vector3

// BenchmarkVectorMath chains the operations a script moving an object does.
func BenchmarkVectorMath(b *testing.B) {
	p, q := Vector3{1, 2, 3}, Vector3{4, 5, 6}
	for b.Loop() {
		d := q.Sub(p)
		if d.SqrLength() > 0 {
			p = p.Add(d.Normalized().Scale(0.01))
		}
	}
	sinkV = p
}

// BenchmarkRotate turns a vector by a quaternion built from an axis and angle.
func BenchmarkRotate(b *testing.B) {
	v := Vector3{1, 0, 0}
	for b.Loop() {
		v = AxisAngle(Vector3{Z: 1}, 0.01).Rotate(v)
	}
	sinkV = v
}
