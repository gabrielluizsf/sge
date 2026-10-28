package sound

// Clip is a mono sound: its samples, from -1 to 1, Rate a second.
type Clip struct {
	Rate    int
	Samples []float32
}

// Rate is the rate the clips the package makes up are made at.
const Rate = 44100
