package human

import "math"

// The real body's legs, in metres, as measured at rest, and where a hip is
// from the root bone and the root bone from under the feet.
const (
	reachThigh = 0.448
	reachShank = 0.475
	hipOut     = 0.115
	hipUp      = 0.005
	hipAhead   = 0.061
	rootHeight = 1.005
	rootAhead  = -0.065
)

// rootAt is where p puts the root bone.
func rootAt(p Pose) [3]float64 {
	return [3]float64{p.Aside, rootHeight + p.Bob, rootAhead + p.Forward}
}

// legTo bends leg s of p for its ankle to be at the spot, where p's root and
// turn put its hip, as far as it reaches.
func legTo(p *Pose, s int, spot [3]float64) {
	out := sides[s].out
	r := rootAt(*p)
	v := around(yAxis, -p.Turn).turn([3]float64{spot[0] - r[0], spot[1] - r[1], spot[2] - r[2]})
	v = [3]float64{v[0] - out*hipOut, v[1] - hipUp, v[2] - hipAhead}
	d := math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
	d = min(max(d, reachShank-reachThigh+0.01), reachThigh+reachShank-0.002)
	knee := math.Acos(clampUnit((d*d - reachThigh*reachThigh - reachShank*reachShank) / (2 * reachThigh * reachShank)))
	along, back := reachThigh+reachShank*math.Cos(knee), reachShank*math.Sin(knee)
	out2 := math.Asin(clampUnit(v[0] / along))
	p.Knee[s], p.Abduct[s] = knee, out*out2
	p.Hip[s] = math.Atan2(v[2], -v[1]) + math.Atan2(back, along*math.Cos(out2))
}

// clampUnit is v held to what a sine or a cosine can be.
func clampUnit(v float64) float64 { return min(max(v, -1), 1) }
