package transform

import (
	"math"
	"testing"

	"github.com/gabrielluizsf/sge/vector"
)

func TestNewIsNaturalSize(t *testing.T) {
	tr := New()
	if tr.Scale != (vector.Vector3{X: 1, Y: 1, Z: 1}) {
		t.Errorf("Scale = %v, want 1,1,1", tr.Scale)
	}
	if tr.Rotation != vector.Identity() {
		t.Errorf("Rotation = %v, want identity", tr.Rotation)
	}
}

func TestTranslateAccumulates(t *testing.T) {
	tr := New()
	tr.Translate(vector.Vector3{X: 1, Y: 2})
	tr.Translate(vector.Vector3{X: 1, Z: 3})
	if want := (vector.Vector3{X: 2, Y: 2, Z: 3}); tr.Position != want {
		t.Errorf("Position = %v, want %v", tr.Position, want)
	}
}

func TestScaleByMultiplies(t *testing.T) {
	tr := New()
	tr.ScaleBy(vector.Vector3{X: 2, Y: 3, Z: 1})
	tr.ScaleBy(vector.Vector3{X: 2, Y: 1, Z: 1})
	if want := (vector.Vector3{X: 4, Y: 3, Z: 1}); tr.Scale != want {
		t.Errorf("Scale = %v, want %v", tr.Scale, want)
	}
}

func TestRotateComposes(t *testing.T) {
	tr := New()
	tr.Rotate(vector.Vector3{Z: 1}, math.Pi/2)
	tr.Rotate(vector.Vector3{Z: 1}, math.Pi/2)
	got := tr.Rotation.Rotate(vector.Vector3{X: 1})
	if math.Abs(got.X+1) > 1e-9 || math.Abs(got.Y) > 1e-9 {
		t.Errorf("two quarter turns sent X to %v, want -X", got)
	}
}

func TestCopiesAreIndependent(t *testing.T) {
	a := New()
	b := a
	b.Translate(vector.Vector3{X: 5})
	if a.Position.X != 0 {
		t.Errorf("moving a copy moved the original to %v", a.Position)
	}
}
