package human

import (
	"strings"

	"github.com/gabrielluizsf/sge/gltf"
)

// insideMouth is whether the proxy in file name sits inside the mouth:
// teeth and tongue, which no sky reaches.
func insideMouth(name string) bool {
	return strings.HasPrefix(name, "teeth/") || strings.HasPrefix(name, "tongue/")
}

// mouthShade is the factor of the light that reaches the inside of a
// mouth: dim, and warm as lips and cheeks make the light that passes them.
var mouthShade = [3]float64{0.62, 0.52, 0.38}

// shadeMouth tints material m as the inside of a mouth lit through lips:
// the renderer has no occlusion there, so teeth would take the sky's blue.
func shadeMouth(m *gltf.Material) {
	for i, s := range mouthShade {
		m.BaseColor[i] *= s
	}
}
