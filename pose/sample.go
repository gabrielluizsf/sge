package pose

import "github.com/gabrielluizsf/sge/gltf"

// Sample is a channel's value at time t: held before its first key and
// after its last, stepped or eased between them, rotations the short way
// round.
func Sample(ch *gltf.Channel, t float64) [4]float64 {
	n := len(ch.Times)
	if n == 0 || len(ch.Values) < n {
		return [4]float64{}
	}
	// The first key after t, found by halves.
	lo, hi := 0, n
	for lo < hi {
		mid := (lo + hi) / 2
		if float64(ch.Times[mid]) > t {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	if lo == 0 || lo == n {
		return wide(ch.Values[max(lo-1, 0)])
	}
	t0, t1 := float64(ch.Times[lo-1]), float64(ch.Times[lo])
	return between(ch, lo-1, (t-t0)/max(t1-t0, 1e-9))
}

// between is a channel's value f of the way from its key i to the next.
func between(ch *gltf.Channel, i int, f float64) [4]float64 {
	a := wide(ch.Values[i])
	if f <= 0 || ch.Interpolation == "STEP" {
		return a
	}
	b := wide(ch.Values[i+1])
	if ch.Path == "rotation" {
		return Slerp(a, b, f)
	}
	return [4]float64{a[0] + (b[0]-a[0])*f, a[1] + (b[1]-a[1])*f, a[2] + (b[2]-a[2])*f, a[3] + (b[3]-a[3])*f}
}

// Length is how long an animation runs: its last keyframe's time.
func Length(a *gltf.Animation) float64 {
	var l float64
	for _, ch := range a.Channels {
		if n := len(ch.Times); n > 0 {
			l = max(l, float64(ch.Times[n-1]))
		}
	}
	return l
}

func wide(v [4]float32) [4]float64 {
	return [4]float64{float64(v[0]), float64(v[1]), float64(v[2]), float64(v[3])}
}
