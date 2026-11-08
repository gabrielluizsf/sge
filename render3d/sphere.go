package render3d

import (
	"math"

	"github.com/gabrielluizsf/sge/pose"
)

// Sphere is a ball around a model, for culling and choosing its detail.
type Sphere struct {
	Center [3]float64
	Radius float64
}

// placed is s moved by m, its radius grown by m's largest scale.
func (s Sphere) placed(m pose.Mat) Sphere {
	sc := 0.0
	for c := range 3 {
		sc = max(sc, m[c*4]*m[c*4]+m[c*4+1]*m[c*4+1]+m[c*4+2]*m[c*4+2])
	}
	return Sphere{Center: m.Point(s.Center), Radius: s.Radius * math.Sqrt(sc)}
}

// onScreen is how far across the screen s is, in pixels, and how deep,
// 0 nearest to 1, seen by c on a screen of size pixels.
func onScreen(c *Camera, s Sphere, size [2]int) (pixels, depth float64) {
	m := &c.ViewProj
	p := s.Center
	z := m[2]*p[0] + m[6]*p[1] + m[10]*p[2] + m[14]
	w := m[3]*p[0] + m[7]*p[1] + m[11]*p[2] + m[15]
	w = max(math.Abs(w), 1e-6)
	sx := math.Sqrt(m[0]*m[0]+m[4]*m[4]+m[8]*m[8]) * float64(size[0]) / 2
	sy := math.Sqrt(m[1]*m[1]+m[5]*m[5]+m[9]*m[9]) * float64(size[1]) / 2
	return 2 * s.Radius * max(sx, sy) / w, z / w
}
