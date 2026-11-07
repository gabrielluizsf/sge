package sound

import "testing"

// A quiet sample passes as it is; a loud one stays in range, louder input
// still louder, and from both sides alike.
func TestLimitBendsOnlyWhatIsLoud(t *testing.T) {
	if limit(0.5) != 0.5 || limit(-0.5) != -0.5 {
		t.Error("a quiet sample changed")
	}
	last := float32(0)
	for x := float32(0.7); x < 20; x += 0.05 {
		y := limit(x)
		if y > 1 || y < last {
			t.Fatalf("limit(%v) = %v after %v", x, y, last)
		}
		if limit(-x) != -y {
			t.Fatalf("limit(%v) is not the opposite of limit(%v)", -x, x)
		}
		last = y
	}
}

// constClip is a second of one level.
func constClip(v float32) *Clip {
	c := &Clip{Rate: 100, Samples: make([]float32, 100)}
	for i := range c.Samples {
		c.Samples[i] = v
	}
	return c
}

// A sound repeated fast follows itself: the one before fades out, and the
// two no more than begin to overlap. Unnamed sounds still overlap.
func TestASoundAloneFollowsItself(t *testing.T) {
	m := NewMixer(100)
	v := Voice{Gain: 1, Alone: "tick"}
	m.Play(constClip(0.2), v)
	m.Play(constClip(0.2), v)
	block := make([]float32, 40) // 20 frames: 0.2 s, past the fade
	m.Fill(block)
	if len(m.playing) != 1 {
		t.Fatalf("%d playing after the fade", len(m.playing))
	}
	free := NewMixer(100)
	free.Play(constClip(0.2), Voice{Gain: 1})
	free.Play(constClip(0.2), Voice{Gain: 1})
	free.Fill(block)
	if len(free.playing) != 2 {
		t.Errorf("%d playing with no name", len(free.playing))
	}
}

// A sound of one name does not touch another's.
func TestAloneOnlyFadesItsOwnName(t *testing.T) {
	m := NewMixer(100)
	m.Play(constClip(0.2), Voice{Gain: 1, Alone: "tick"})
	m.Play(constClip(0.2), Voice{Gain: 1, Alone: "open"})
	m.Fill(make([]float32, 40))
	if len(m.playing) != 2 {
		t.Errorf("%d playing", len(m.playing))
	}
}
