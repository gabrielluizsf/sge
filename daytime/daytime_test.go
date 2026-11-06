package daytime

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
)

// The clock goes round, counting the days.
func TestTheClockGoesRound(t *testing.T) {
	c := &Clock{Hour: 22, Speed: 1}
	for range 5 {
		c.Update(nil, 1)
	}
	if c.Day != 1 || c.Hour != 3 {
		t.Errorf("day %d hour %v", c.Day, c.Hour)
	}
	if !c.Between(22, 6) || c.Between(6, 22) || !c.Between(0, 4) {
		t.Error("Between is wrong round midnight")
	}
}

// Noon is light, midnight dark, and dusk goes from one to the other.
func TestTheLightFollowsTheSun(t *testing.T) {
	light := func(h float64) float64 { return (&Clock{Hour: h}).Light() }
	if light(12) != 1 || light(0) != 0 || light(23) != 0 {
		t.Fatal("noon or midnight is wrong")
	}
	last := 1.0
	for h := Dusk; h <= Night; h += 0.25 {
		if l := light(h); l > last {
			t.Fatalf("the dusk brightens at %v", h)
		} else {
			last = l
		}
	}
	if l := light(6); l <= 0 || l >= 1 {
		t.Errorf("dawn is %v", l)
	}
}

// By day the shade changes nothing; by night it turns the screen moonlit
// blue and dark.
func TestTheShadeDarkensTheNight(t *testing.T) {
	c := &Clock{Hour: 12}
	s := Shade{Clock: c}
	cv, _ := canvas.NewCanvas(4, 4)
	cv.Clear(canvas.RGB(200, 200, 200))
	s.Paint(nil, cv)
	if cv.At(1, 1) != canvas.RGB(200, 200, 200) {
		t.Fatal("the day was shaded")
	}
	c.Hour = 1
	s.Paint(&game.Object{}, cv)
	if p := cv.At(1, 1); p.R() >= 100 || p.B() <= p.R() {
		t.Errorf("night is %v", p)
	}
}

// One wakes at six: the light is bright and warm, a sunrise, not a dull
// veil, which is how the morning looked when the screen was covered in a
// mix of night and glow.
func TestTheMorningIsBrightAndWarm(t *testing.T) {
	a := (&Clock{Hour: 6}).Ambient()
	if a[0] < 0.85 || a[1] < 0.7 || a[2] < 0.6 {
		t.Errorf("six in the morning is dim: %v", a)
	}
	if a[0] <= a[2] {
		t.Errorf("six in the morning is not warm: %v", a)
	}
	if n := (&Clock{Hour: 8}).Ambient(); n != (Ambient{1, 1, 1}) {
		t.Errorf("eight is not full day: %v", n)
	}
}

// The evening goes through gold, then rose, into blue night, and no hour
// jumps: the light changes a little at a time all round the clock.
func TestTheLightChangesSmoothlyRoundTheClock(t *testing.T) {
	gold := (&Clock{Hour: 18}).Ambient()
	if gold[0] < 0.9 || gold[0]-gold[2] < 0.2 {
		t.Errorf("six in the evening is not golden: %v", gold)
	}
	night := (&Clock{Hour: 23}).Ambient()
	if night[2] <= night[0] {
		t.Errorf("the night is not blue: %v", night)
	}
	last := (&Clock{Hour: 0}).Ambient()
	const step = 0.1
	for h := step; h <= 24; h += step {
		a := (&Clock{Hour: h}).Ambient()
		for k := range a {
			if d := a[k] - last[k]; d > 0.1 || d < -0.1 {
				t.Fatalf("the light jumps at %.1f: %v to %v", h, last, a)
			}
		}
		last = a
	}
}

// The shade multiplies: what is twice as bright stays twice as bright, so
// the picture keeps its contrast, and black stays black, unlike a colour
// laid over it. Alpha is left alone.
func TestTheShadeKeepsTheContrast(t *testing.T) {
	cv, _ := canvas.NewCanvas(3, 1)
	cv.Put(0, 0, canvas.RGB(100, 100, 100))
	cv.Put(1, 0, canvas.RGB(200, 200, 200))
	cv.Put(2, 0, canvas.RGB(0, 0, 0))
	Shade{Clock: &Clock{Hour: 6}}.Paint(nil, cv)
	lo, hi := cv.At(0, 0), cv.At(1, 0)
	if d := int(hi.R()) - 2*int(lo.R()); d < -2 || d > 2 {
		t.Errorf("%d is not twice %d", hi.R(), lo.R())
	}
	if cv.At(2, 0) != canvas.RGB(0, 0, 0) {
		t.Errorf("black became %v", cv.At(2, 0))
	}
	if lo.A() != 255 {
		t.Errorf("alpha %d", lo.A())
	}
}

// A whole screen is shaded in a small part of a frame.
func BenchmarkShade(b *testing.B) {
	cv, _ := canvas.NewCanvas(1280, 720)
	s := Shade{Clock: &Clock{Hour: 6}}
	for b.Loop() {
		cv.Clear(canvas.RGB(120, 160, 90))
		s.Paint(nil, cv)
	}
}

// A schedule says what to do at each hour.
func TestASchedulePicksTheSlotOfTheHour(t *testing.T) {
	s := Schedule{{8, 18, "work"}, {18, 21, "stroll"}, {21, 7, "sleep"}}
	for h, want := range map[float64]string{9: "work", 18.5: "stroll", 23: "sleep", 3: "sleep", 7.5: ""} {
		if got := s.At(&Clock{Hour: h}); got != want {
			t.Errorf("at %v: %q, not %q", h, got, want)
		}
	}
}
