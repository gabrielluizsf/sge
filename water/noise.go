package water

import (
	"math"

	"github.com/gabrielluizsf/sge/terrain"
)

// noise is smooth value noise from 0 to 1 at uv, repeating every 1 in u
// and v, period lattice cells across: it wraps, so tiles of it meet.
func noise(uv [2]float64, period int, seed uint32) float64 {
	x, y := uv[0]*float64(period), uv[1]*float64(period)
	x0, y0 := math.Floor(x), math.Floor(y)
	fx, fy := x-x0, y-y0
	fx, fy = fx*fx*(3-2*fx), fy*fy*(3-2*fy)
	i, j := int(x0)+period*4, int(y0)+period*4
	at := func(di, dj int) float64 {
		return float64(terrain.Hash(seed, (i+di)%period, (j+dj)%period)&0xffff) / 65535
	}
	return mix(mix(at(0, 0), at(1, 0), fx), mix(at(0, 1), at(1, 1), fx), fy)
}
