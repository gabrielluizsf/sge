package sound

// NewClip is a clip of samples made at Rate, its loudness scaled so its
// peak is peak.
func NewClip(samples []float32, peak float64) *Clip {
	var top float32
	for _, x := range samples {
		top = max(top, x, -x)
	}
	if top > 0 {
		g := float32(peak) / top
		for i := range samples {
			samples[i] *= g
		}
	}
	return &Clip{Rate: Rate, Samples: samples}
}
