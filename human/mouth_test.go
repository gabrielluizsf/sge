package human

import (
	"testing"

	"github.com/gabrielluizsf/sge/gltf"
)

func TestInsideMouth(t *testing.T) {
	for name, want := range map[string]bool{
		"teeth/teeth_base/teeth_base.mhclo": true,
		"tongue/tongue01/tongue01.mhclo":    true,
		"hair/short02/short02.mhclo":        false,
		"eyes/high-poly/high-poly.mhclo":    false,
	} {
		if got := insideMouth(name); got != want {
			t.Errorf("insideMouth(%q) = %v, want %v", name, got, want)
		}
	}
}

// TestMouthShade: the mouth's inside is lit by no sky, so teeth are dimmed
// and warm (red above blue), never the blue the open sky would make them.
func TestMouthShade(t *testing.T) {
	m := gltf.Material{BaseColor: [4]float64{1, 1, 1, 1}}
	shadeMouth(&m)
	c := m.BaseColor
	if c[0] >= 0.8 || c[0] <= c[1] || c[1] <= c[2] || c[3] != 1 {
		t.Errorf("mouth shade %v is not a dim warm ivory", c)
	}
}
