package main

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// A dark thing on the dark panel gets a pale edge round it, and its own
// pixels are left as they were: nothing it shows is lost.
func TestADarkIconIsSeenOnThePanel(t *testing.T) {
	cv, err := canvas.NewCanvas(16, 16)
	if err != nil {
		t.Fatal(err)
	}
	dark := canvas.RGB(20, 20, 22)
	for y := 6; y < 10; y++ {
		for x := 6; x < 10; x++ {
			cv.Pixels[y*cv.Stride+x] = dark
		}
	}
	out := halo(cv)
	if out.Pixels[7*out.Stride+7] != dark {
		t.Error("the icon's own pixel changed")
	}
	edge := out.Pixels[7*out.Stride+5]
	if edge.A() == 0 || edge.R() < 200 {
		t.Errorf("the edge beside it: %d %d %d %d", edge.R(), edge.G(), edge.B(), edge.A())
	}
	if far := out.Pixels[0]; far.A() != 0 {
		t.Errorf("a pixel far off has %d", far.A())
	}
}
