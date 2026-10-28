package transform

import "github.com/gabrielluizsf/sge/vector"

// ScaleBy multiplies the transform's scale by factor, axis by axis.
func (t *Transform) ScaleBy(factor vector.Vector3) {
	t.Scale = t.Scale.Mul(factor)
}
