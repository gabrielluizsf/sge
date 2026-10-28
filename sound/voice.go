package sound

// Voice is how a clip is played: Gain its loudness, 1 as made; Pan from
// -1, the left, to 1, the right; Pitch how much faster it is played, 1, or
// 0, as made.
type Voice struct {
	Gain, Pan, Pitch float64
}
