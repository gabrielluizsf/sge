package sprite

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/transform"
	"github.com/gabrielluizsf/sge/vector"
)

// BenchmarkDrawRotated draws a 24x24 sprite turned and scaled.
func BenchmarkDrawRotated(b *testing.B) {
	s, err := Rect(24, 24, canvas.RGB(80, 170, 255))
	if err != nil {
		b.Fatal(err)
	}
	cv, err := canvas.NewCanvas(480, 320)
	if err != nil {
		b.Fatal(err)
	}
	t := transform.New()
	t.Position = vector.Vector3{X: 240, Y: 160}
	t.Rotate(vector.Vector3{Z: 1}, 0.5)
	t.ScaleBy(vector.Vector3{X: 2, Y: 2, Z: 1})
	for b.Loop() {
		s.Draw(cv, t)
	}
}
