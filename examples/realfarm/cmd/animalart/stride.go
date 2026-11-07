package main

import (
	"math"
	"os"
	"sort"

	"github.com/gabrielluizsf/sge/animal"
)

// measureHorse is how fast, in metres a second, a horse's body goes for its
// planted hooves not to slide in the model's own walk: how fast its hind
// hooves go back while they are on the ground, in the median.
func measureHorse(root *os.Root) (float64, error) {
	m, err := open(root, "draft_horse_walk")
	if err != nil {
		return 0, err
	}
	r, err := animal.NewRig(m)
	if err != nil {
		return 0, err
	}
	fwd := r.Toward(r.Bone("spine.001_00"), r.Bone("skull_015"))
	l := math.Hypot(fwd[0], fwd[2])
	fwd = animal.Vec{fwd[0] / l, 0, fwd[2] / l}
	n := 64
	period := 1.3333
	times := make([]float64, n+1)
	for i := range times {
		times[i] = period * float64(i) / float64(n)
	}
	var speeds []float64
	for _, h := range []string{"r_hoof.L_end_045", "r_hoof.R_end_046", "f_hoof.L_end_043", "f_hoof.R_end_044"} {
		path := r.Track("draft_Walk", r.Bone(h), times)
		low := path[0][1]
		for _, p := range path {
			low = min(low, p[1])
		}
		for i := 1; i <= n; i++ {
			if path[i][1] < low+0.02 && path[i-1][1] < low+0.02 {
				d := (path[i-1][0]-path[i][0])*fwd[0] + (path[i-1][2]-path[i][2])*fwd[2]
				speeds = append(speeds, d/(period/float64(n)))
			}
		}
	}
	sort.Float64s(speeds)
	return speeds[len(speeds)/2], nil
}
