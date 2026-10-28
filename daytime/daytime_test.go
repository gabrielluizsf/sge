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

// By day the shade paints nothing; by night it darkens the screen.
func TestTheShadeDarkensTheNight(t *testing.T) {
	c := &Clock{Hour: 12}
	s := Shade{Clock: c, Night: canvas.RGBA(canvas.RGBAParams{R: 10, G: 20, B: 60, A: 150}), Glow: canvas.RGBA(canvas.RGBAParams{R: 255, G: 120, B: 40, A: 80})}
	cv, _ := canvas.NewCanvas(4, 4)
	cv.Clear(canvas.RGB(200, 200, 200))
	s.Paint(nil, cv)
	if cv.At(1, 1) != canvas.RGB(200, 200, 200) {
		t.Fatal("the day was shaded")
	}
	c.Hour = 1
	s.Paint(&game.Object{}, cv)
	if p := cv.At(1, 1); p.R() >= 200 || p.B() <= p.R() {
		t.Errorf("night is %v", p)
	}
	c.Hour = 19
	if a := s.Tint().A(); a == 0 || a >= 150+40 {
		t.Errorf("dusk's alpha is %d", a)
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
