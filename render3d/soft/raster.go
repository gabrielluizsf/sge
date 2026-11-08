package soft

import "math"

// raster fills the pixels whose middles tri covers, nearer than what is
// there, with its shading; a triangle turned away is left out unless its
// pipeline draws both sides.
func (p *pass) raster(s *state, tri [3]vout) {
	if p.depth == nil {
		return
	}
	var sc [3][3]float64 // x, y in pixels, depth
	var iw [3]float64    // 1 / w, for interpolating in perspective
	for k, v := range tri {
		if v.clip[3] <= 0 {
			return
		}
		iw[k] = 1 / v.clip[3]
		sc[k] = [3]float64{(v.clip[0]*iw[k] + 1) / 2 * float64(p.w), (1 - v.clip[1]*iw[k]) / 2 * float64(p.h), v.clip[2] * iw[k]}
	}
	area := edge(sc[0], sc[1], sc[2])
	// Y runs down the screen: a triangle counterclockwise in clip space
	// has a negative area here.
	front := area < 0
	if area == 0 || (!front && !s.double) {
		return
	}
	x0 := max(0, int(math.Floor(min(sc[0][0], sc[1][0], sc[2][0]))))
	x1 := min(p.w-1, int(math.Ceil(max(sc[0][0], sc[1][0], sc[2][0]))))
	y0 := max(0, int(math.Floor(min(sc[0][1], sc[1][1], sc[2][1]))))
	y1 := min(p.h-1, int(math.Ceil(max(sc[0][1], sc[1][1], sc[2][1]))))
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			c := [3]float64{float64(x) + 0.5, float64(y) + 0.5}
			b := [3]float64{edge(sc[1], sc[2], c) / area, edge(sc[2], sc[0], c) / area, edge(sc[0], sc[1], c) / area}
			if b[0] < 0 || b[1] < 0 || b[2] < 0 {
				continue
			}
			z := b[0]*sc[0][2] + b[1]*sc[1][2] + b[2]*sc[2][2]
			i := y*p.w + x
			if z < 0 || z > 1 || z >= float64(p.depth[i]) {
				continue
			}
			pw := [3]float64{b[0] * iw[0], b[1] * iw[1], b[2] * iw[2]}
			sum := pw[0] + pw[1] + pw[2]
			f := frag{tri: &tri, w: [3]float64{pw[0] / sum, pw[1] / sum, pw[2] / sum}, back: !front}
			if p.plot(s, f, i) {
				p.depth[i] = float32(z)
			}
		}
	}
}

// edge is twice the signed area of a, b, c.
func edge(a, b, c [3]float64) float64 {
	return (b[0]-a[0])*(c[1]-a[1]) - (b[1]-a[1])*(c[0]-a[0])
}
