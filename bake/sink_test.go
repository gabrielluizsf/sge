package bake

import (
	"bytes"
	"testing"
)

// A shot sunk below the ground is the model stood that much lower: a
// horse wading is baked as if its node stood in the water, the same to
// the pixel, its feet held where they are.
func TestASunkShotIsTheModelStoodLower(t *testing.T) {
	o := options()
	o.Water = &Water{Color: [3]float64{0.02, 0.1, 0.12}, Clear: 0.3}
	sunk := Render(box(1, 2, [3]float64{0.8, 0.2, 0.2}), Shot{Sink: 1}, o)
	lowered := box(1, 2, [3]float64{0.8, 0.2, 0.2})
	lowered.Nodes[0].Translation = [3]float64{0, -1, 0}
	want := Render(lowered, Shot{}, o)
	if sunk.Pivot != want.Pivot || sunk.Image.Bounds() != want.Image.Bounds() || !bytes.Equal(sunk.Image.Pix, want.Image.Pix) {
		t.Errorf("sunk %v at %v, lowered %v at %v", sunk.Image.Bounds(), sunk.Pivot, want.Image.Bounds(), want.Pivot)
	}
	if dry := Render(box(1, 2, [3]float64{0.8, 0.2, 0.2}), Shot{}, o); bytes.Equal(dry.Image.Pix, sunk.Image.Pix) {
		t.Error("sinking changed nothing")
	}
}
