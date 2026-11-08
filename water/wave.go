package water

import (
	"math"

	"github.com/gabrielluizsf/sge/random"
)

// wave is one sine wave of the sea: a whole number of crests across and
// down a tile, so it wraps, how steep it is, where it starts, and which
// way it travels.
type wave struct {
	kx, ky float64
	slope  float64
	phase  float64
	dir    float64
}

// sharp is how much the crests are drawn up and the troughs flattened,
// as Gerstner's waves are, by bending the sine's phase.
const sharp = 0.45

// waves are the waves of a sea of seed: a dozen, between two and eight
// crests across each cell of a map span cells across, most travelling the
// same way, steepness their slope together.
func waves(seed uint32, steepness float64, span int) []wave {
	r := random.New(seed*2246822519 + 7)
	out := make([]wave, 12)
	for i := range out {
		a := 0.9*(r.Float()-0.5)*math.Pi + 0.6
		k := float64(span) * (2 + 6*r.Float()*r.Float())
		kx, ky := math.Round(k*math.Cos(a)), math.Round(k*math.Sin(a))
		if kx == 0 && ky == 0 {
			kx = 2
		}
		dir := 1.0
		if r.Float() < 0.25 {
			dir = -1
		}
		out[i] = wave{kx: kx, ky: ky, slope: steepness / float64(span) * (0.4 + 0.6*r.Float()) / (1 + 0.08*math.Hypot(kx, ky)) / 9, phase: 2 * math.Pi * r.Float(), dir: dir}
	}
	return out
}

// slopeAt is the surface's slope along u and along v at uv, turned t of
// a loop on, from 0 to 1: after a whole loop each wave is where it began.
func slopeAt(ws []wave, uv [2]float64, t float64) (float64, float64) {
	u, v := uv[0], uv[1]
	var sx, sy float64
	for _, w := range ws {
		th := 2*math.Pi*(w.kx*u+w.ky*v+w.dir*t) + w.phase
		bent := th + sharp*math.Sin(th)
		s := w.slope * math.Cos(bent) * (1 + sharp*math.Cos(th))
		n := math.Hypot(w.kx, w.ky)
		sx += s * w.kx / n
		sy += s * w.ky / n
	}
	return sx, sy
}
