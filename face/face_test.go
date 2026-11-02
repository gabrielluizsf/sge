package face

import "testing"

// Eyes blink at a person's pace: never sooner than two seconds apart nor
// later than six, each blink shutting them, then opening them, in about a
// sixth of a second; the same way every run of a seed.
func TestBlinksAtAPersonsPace(t *testing.T) {
	const dt = 1.0 / 120
	b, again := NewBlinker(7), NewBlinker(7)
	var blinks []float64
	shut, was := 0, LidsOpen
	for i := range int(60 / dt) {
		b.Update(dt)
		again.Update(dt)
		if b.Lids() != again.Lids() {
			t.Fatal("a seed blinks two ways")
		}
		if b.Lids() == LidsShut {
			shut++
		}
		if b.Lids() != LidsOpen && was == LidsOpen {
			blinks = append(blinks, float64(i)*dt)
		}
		was = b.Lids()
	}
	if len(blinks) < 60/6 || len(blinks) > 60/2+1 {
		t.Fatalf("%d blinks in a minute", len(blinks))
	}
	for i := 1; i < len(blinks); i++ {
		if gap := blinks[i] - blinks[i-1]; gap < blinkLeast || gap > blinkMost+0.2 {
			t.Errorf("blinks %.2f s apart", gap)
		}
	}
	if per := float64(shut) * dt / float64(len(blinks)); per < 0.04 || per > 0.08 {
		t.Errorf("eyes shut %.3f s a blink", per)
	}
}

// The mouth takes the shapes of the letters said, holding each a moment,
// and shuts when the words are all said.
func TestTheMouthSaysTheLetters(t *testing.T) {
	for r, want := range map[rune]Shape{'a': Open, 'Ã': Open, 'é': Open, 'o': Round, 'ú': Round, 'm': Shut, 't': Ajar, ' ': Shut, '!': Shut} {
		if got := ShapeOf(r); got != want {
			t.Errorf("%q is %v, not %v", r, got, want)
		}
	}
	var s Speaker
	s.Update(0.01, 'a', true)
	if s.Shape() != Open {
		t.Fatal("a shut mouth did not open for an a")
	}
	s.Update(0.01, 'o', true)
	if s.Shape() != Open {
		t.Error("the mouth changed shape before holding the last")
	}
	s.Update(hold, 'o', true)
	if s.Shape() != Round {
		t.Error("the mouth did not round for an o")
	}
	s.Update(0.01, 0, false)
	if s.Shape() != Shut {
		t.Error("the mouth stayed open with everything said")
	}
}
