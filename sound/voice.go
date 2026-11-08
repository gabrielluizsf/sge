package sound

// Voice is how a clip is played: Gain its loudness, 1 as made; Pan from
// -1, the left, to 1, the right; Pitch how much faster it is played, 1, or
// 0, as made. Alone, when not empty, names what the clip is: another of
// that name still playing fades out quickly as this one starts, so a sound
// repeated fast follows itself rather than piling up.
type Voice struct {
	Gain, Pan, Pitch float64
	Alone            string
}
