package face

import "github.com/gabrielluizsf/sge/random"

// A blink: the seconds the lids take shutting, staying shut and opening,
// and the least and most seconds between blinks.
const (
	blinkShutting = 0.05
	blinkShut     = 0.06
	blinkOpening  = 0.06
	blinkLeast    = 2.0
	blinkMost     = 6.0
)

// Blinker is a pair of eyes blinking at a person's pace, each wait drawn
// from its seed, so a face blinks the same way every run.
type Blinker struct {
	src  *random.Source
	wait float64 // seconds to the next blink
	t    float64 // seconds into a blink, -1 between blinks
}

// NewBlinker is eyes that open, and first blink in a moment drawn from
// seed.
func NewBlinker(seed uint32) *Blinker {
	b := &Blinker{src: random.New(seed), t: -1}
	b.wait = b.pause()
	return b
}

// pause is a wait between blinks.
func (b *Blinker) pause() float64 {
	return blinkLeast + (blinkMost-blinkLeast)*b.src.Float()
}

// Update lets dt seconds go by.
func (b *Blinker) Update(dt float64) {
	if b.t < 0 {
		if b.wait -= dt; b.wait <= 0 {
			b.t = 0
		}
		return
	}
	if b.t += dt; b.t >= blinkShutting+blinkShut+blinkOpening {
		b.t, b.wait = -1, b.pause()
	}
}

// Lids is how far shut the eyelids are now.
func (b *Blinker) Lids() int {
	switch {
	case b.t < 0:
		return LidsOpen
	case b.t < blinkShutting:
		return LidsHalf
	case b.t < blinkShutting+blinkShut:
		return LidsShut
	}
	return LidsHalf
}
