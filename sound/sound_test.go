package sound

import (
	"math"
	"testing"
)

func ones(n int) *Clip {
	s := make([]float32, n)
	for i := range s {
		s[i] = 0.5
	}
	return &Clip{Rate: Rate, Samples: s}
}

// frames is how many frames of block are not silent.
func frames(block []float32) int {
	n := 0
	for i := 0; i < len(block); i += 2 {
		if block[i] != 0 || block[i+1] != 0 {
			n++
		}
	}
	return n
}

// No sample leaves -1 to 1, and one that is not a number is silence: a
// clip made wrong cannot make the speaker click.
func TestFillKeepsSamplesInRange(t *testing.T) {
	m := NewMixer(Rate)
	nan := float32(math.NaN())
	m.Play(&Clip{Rate: Rate, Samples: []float32{5, -5, nan, float32(math.Inf(1)), 0}}, Voice{Gain: 1})
	block := make([]float32, 10)
	m.Fill(block)
	for i, s := range block {
		if s != s || s > 1 || s < -1 {
			t.Fatalf("sample %d is %v", i, s)
		}
	}
}

// A clip plays once, then is gone; at twice the pitch, or into a stream of
// half the rate, it lasts half as many frames.
func TestAClipPlaysOnceAtItsPitchAndRate(t *testing.T) {
	for _, c := range []struct {
		rate  int
		pitch float64
		want  int
	}{{Rate, 1, 99}, {Rate, 2, 50}, {Rate / 2, 1, 50}} {
		m := NewMixer(c.rate)
		m.Play(ones(100), Voice{Gain: 1, Pitch: c.pitch})
		block := make([]float32, 400)
		m.Fill(block)
		if got := frames(block); got != c.want {
			t.Errorf("rate %d pitch %v: %d frames, not %d", c.rate, c.pitch, got, c.want)
		}
		m.Fill(block)
		if got := frames(block); got != 0 {
			t.Errorf("rate %d pitch %v: played again, %d frames", c.rate, c.pitch, got)
		}
	}
}

// Panned to the left a sound is all in the left, and in the middle it is
// as loud on both sides as it is played.
func TestPanSharesASoundBetweenTheSides(t *testing.T) {
	m := NewMixer(Rate)
	m.Play(ones(10), Voice{Gain: 0.5, Pan: -1})
	m.Play(ones(10), Voice{Gain: 0.5})
	block := make([]float32, 4)
	m.Fill(block)
	if l, r := block[0], block[1]; math.Abs(float64(l)-(0.25*math.Sqrt2+0.25)) > 1e-5 || math.Abs(float64(r)-0.25) > 1e-5 {
		t.Errorf("left %v, right %v", l, r)
	}
}

// With as many voices playing as a mixer may have, the oldest is cut.
func TestTheOldestVoiceIsCut(t *testing.T) {
	m := NewMixer(Rate)
	m.MaxVoices = 2
	for range 5 {
		m.Play(ones(10), Voice{Gain: 0.1})
	}
	if len(m.playing) != 2 {
		t.Fatalf("%d voices", len(m.playing))
	}
}

// A loop is silent until faded in, eases in without a jump, and goes on
// past its end.
func TestALoopFadesInAndGoesOn(t *testing.T) {
	m := NewMixer(Rate)
	l := m.Loop(ones(50))
	block := make([]float32, 200)
	m.Fill(block)
	if frames(block) != 0 {
		t.Fatal("heard before it was faded in")
	}
	l.Fade(1)
	block = make([]float32, 2*Rate)
	m.Fill(block)
	for i := 2; i < len(block); i += 2 {
		if d := math.Abs(float64(block[i] - block[i-2])); d > 0.01 {
			t.Fatalf("jumps %v at frame %d", d, i/2)
		}
	}
	if v := block[len(block)-2]; math.Abs(float64(v)-0.5) > 0.01 {
		t.Errorf("faded in to %v, not 0.5", v)
	}
}

// A seamless loop's end runs into its start as its middle runs on.
func TestSeamlessLoopsHaveNoSeam(t *testing.T) {
	s := Tone{From: 441, Seconds: 0.1}.Make()
	for i := range s {
		s[i] += float32(i) / float32(len(s)) // a ramp, which would jump
	}
	s = Seamless(s, 400)
	if d := math.Abs(float64(s[0] - s[len(s)-1])); d > 0.08 {
		t.Errorf("the seam jumps %v", d)
	}
}

// An envelope rises, and dies to a thousandth over its decay; a lowpass
// keeps a low tone and softens a high one.
func TestEnvelopesAndFilters(t *testing.T) {
	s := make([]float32, Seconds(0.2))
	for i := range s {
		s[i] = 1
	}
	Envelope{Attack: 0.01, Decay: 0.1}.Apply(s)
	if s[0] != 0 || s[Seconds(0.01)] < 0.99 || math.Abs(float64(s[Seconds(0.11)])-0.001) > 0.0005 {
		t.Errorf("envelope %v, %v, %v", s[0], s[Seconds(0.01)], s[Seconds(0.11)])
	}
	peak := func(f float64) float32 {
		s := Tone{From: f, Seconds: 0.1}.Make()
		Lowpass(s, 500)
		var p float32
		for _, x := range s[len(s)/2:] {
			p = max(p, x)
		}
		return p
	}
	if lo, hi := peak(100), peak(8000); lo < 0.9 || hi > 0.1 {
		t.Errorf("lowpass at 500Hz keeps %v of 100Hz and %v of 8kHz", lo, hi)
	}
}

func BenchmarkFill(b *testing.B) {
	m := NewMixer(Rate)
	for range 3 {
		m.Loop(NewClip(Noise(1, Rate), 0.5)).Fade(0.3)
	}
	block := make([]float32, 2*441)
	c := NewClip(Noise(2, Rate), 0.5)
	for b.Loop() {
		if len(m.playing) < 16 {
			m.Play(c, Voice{Gain: 0.5, Pitch: 1.03})
		}
		m.Fill(block)
	}
}
