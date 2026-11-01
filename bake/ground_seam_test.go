package bake

import (
	"image"
	"image/draw"
	"testing"
)

// Ground tiles laid out as a map lays them, three by three, leave no
// hole: every pixel of the middle cell's diamond, its corners where four
// tiles meet too, is covered.
func TestGroundTilesLeaveNoHoles(t *testing.T) {
	m := box(1, 1, [3]float64{0.4, 0.6, 0.3})
	o := options()
	canvas := image.NewNRGBA(image.Rect(-160, -32, 160, 160))
	for r := range 3 {
		for c := range 3 {
			s := Ground(m, Patch{Span: 2, Col: c % 2, Row: r % 2}, o)
			at := image.Pt((c-r)*32-s.Pivot.X, (c+r)*16-s.Pivot.Y)
			draw.Draw(canvas, s.Image.Bounds().Add(at), s.Image, image.Point{}, draw.Over)
		}
	}
	// The middle cell, from (1, 1) to (2, 2): its diamond on the screen.
	inside := func(x, y float64) bool {
		u, v := x/32, y/16 // u = c-r, v = c+r, the cell's own coordinates
		return v >= 2+absf(u) && v <= 4-absf(u)
	}
	holes, checked := 0, 0
	for y := canvas.Rect.Min.Y; y < canvas.Rect.Max.Y; y++ {
		for x := canvas.Rect.Min.X; x < canvas.Rect.Max.X; x++ {
			if inside(float64(x)+0.5, float64(y)+0.5) {
				checked++
			}
			if inside(float64(x)+0.5, float64(y)+0.5) && canvas.NRGBAAt(x, y).A != 255 {
				holes++
			}
		}
	}
	if checked != 32*64/2 {
		t.Fatalf("the middle cell's diamond is %d pixels, not %d", checked, 32*64/2)
	}
	if holes > 0 {
		t.Errorf("%d pixels of the middle cell are holes", holes)
	}
}

func absf(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
