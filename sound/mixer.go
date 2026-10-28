package sound

import "sync"

// Mixer mixes what is playing into a device's stream. Its methods may be
// called from the game while the device's goroutine fills.
type Mixer struct {
	mu      sync.Mutex
	rate    int
	playing []*playing
	loops   []*Loop
	// MaxVoices is how many clips may play at once; the oldest is cut
	// when another starts. 0 is 32.
	MaxVoices int
}

// playing is a clip being played: where in it the mix is, and how fast
// that moves on.
type playing struct {
	clip        *Clip
	pos, step   float64
	left, right float32
}

// NewMixer is a mixer for a stream of rate frames a second.
func NewMixer(rate int) *Mixer { return &Mixer{rate: rate} }

// SetRate sets the stream's rate, to what the device runs at.
func (m *Mixer) SetRate(rate int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rate = rate
}
