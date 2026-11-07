package main

import (
	"testing"

	"github.com/gabrielluizsf/sge/examples/farm/kit/bakes"
	"github.com/gabrielluizsf/sge/farm"
)

// The facing a chair of each rotation turns the sitter to is baked at the
// yaw that chair's own picture is, so the chair is not seen to jump when
// its tile gives way to the sitter's.
func TestASitterIsBakedAtTheYawOfTheirChair(t *testing.T) {
	for h, dir := range []string{"downleft", "upleft", "upright", "downright"} {
		yaw := float64(turnOf(dir)) * 360 / bakes.Facings
		if yaw != bakes.Headings[h] {
			t.Errorf("heading %d: %s is baked at yaw %v, the chair at %v", h, dir, yaw, bakes.Headings[h])
		}
	}
	if _, ok := testArt(t).(farm.SeatArt); !ok {
		t.Error("the realistic look does not draw the chair with the sitter")
	}
}
