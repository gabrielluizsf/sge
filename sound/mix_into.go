package sound

// MixInto adds src into dst where at says, at its gain, as far as dst
// goes.
func MixInto(dst, src []float32, at Place) {
	for i, x := range src {
		if j := at.Sample + i; j >= 0 && j < len(dst) {
			dst[j] += x * float32(at.Gain)
		}
	}
}

// Place is where and how loud a sound is mixed into another: from its
// Sample on, at Gain.
type Place struct {
	Sample int
	Gain   float64
}
