package transform

import "github.com/gabrielluizsf/sge/vector"

// Translate moves the transform by delta.
func (t *Transform) Translate(delta vector.Vector3) {
	t.Position = t.Position.Add(delta)
}
