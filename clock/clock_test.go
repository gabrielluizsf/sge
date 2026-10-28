package clock

import "testing"

func TestTickMeasuresTheGap(t *testing.T) {
	now := 10.0
	c := New(func() float64 { return now })
	now = 10.25
	if got := c.Tick(); got != 0.25 {
		t.Errorf("first Tick = %v, want 0.25", got)
	}
	now = 11
	if got := c.Tick(); got != 0.75 {
		t.Errorf("second Tick = %v, want 0.75", got)
	}
}

func TestTickNeverGoesBackwards(t *testing.T) {
	now := 5.0
	c := New(func() float64 { return now })
	now = 4
	if got := c.Tick(); got != 0 {
		t.Errorf("Tick = %v, want 0", got)
	}
}
