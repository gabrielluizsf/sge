package vignette

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

var red = Config{Color: canvas.RGB(192, 24, 32), Alpha: 150, From: 0.68}

// None in the middle, much at the edge, and no step to see between.
func TestItGrowsSmoothlyFromTheMiddle(t *testing.T) {
	const w, h = 320, 180
	v := New(w, h, red)
	if a := v.AlphaAt(w/2, h/2); a != 0 {
		t.Errorf("the middle has alpha %d", a)
	}
	if a := v.AlphaAt(0, h/2); a < 110 {
		t.Errorf("the edge has alpha only %d", a)
	}
	for x := 1; x < w; x++ {
		d := int(v.AlphaAt(x, h/2)) - int(v.AlphaAt(x-1, h/2))
		if d < -8 || d > 8 {
			t.Fatalf("alpha jumps by %d at x %d", d, x)
		}
	}
}

// Painting leaves what it does not cover, reddens the edge as much as the
// colour's share says, keeps the canvas's alpha and fades with opacity.
func TestPaintBlendsTheColour(t *testing.T) {
	const w, h = 64, 36
	v := New(w, h, red)
	base := canvas.RGB(20, 200, 60)
	cv, _ := canvas.NewCanvas(w, h)
	cv.Clear(base)
	v.Paint(cv, 1)
	if cv.At(w/2, h/2) != base {
		t.Errorf("the middle became %v", cv.At(w/2, h/2))
	}
	for _, at := range [][2]int{{0, 0}, {w - 1, 0}, {0, h - 1}, {w - 1, h - 1}, {0, h / 2}, {w - 1, h / 2}, {w / 2, 0}, {w / 2, h - 1}} {
		p := cv.At(at[0], at[1])
		a := int(v.AlphaAt(at[0], at[1]))
		want := (20*(255-a) + 192*a) / 255
		if d := int(p.R()) - want; d < -2 || d > 2 {
			t.Errorf("at %v red is %d, want about %d", at, p.R(), want)
		}
		wantG := (200*(255-a) + 24*a) / 255
		if d := int(p.G()) - wantG; d < -2 || d > 2 {
			t.Errorf("at %v green is %d, want about %d", at, p.G(), wantG)
		}
		if p.A() != 255 {
			t.Errorf("at %v alpha is %d", at, p.A())
		}
	}
	half, _ := canvas.NewCanvas(w, h)
	half.Clear(base)
	v.Paint(half, 0.5)
	if half.At(0, h/2).R() >= cv.At(0, h/2).R() {
		t.Error("half strength is as red as full")
	}
	none, _ := canvas.NewCanvas(w, h)
	none.Clear(base)
	v.Paint(none, 0)
	if none.At(0, 0) != base {
		t.Error("no strength painted")
	}
}

// Every covered pixel is painted, on both sides of every row, however the
// rows are cut: the shortcut past the middle skips nothing.
func TestPaintMissesNothingItCovers(t *testing.T) {
	for _, size := range [][2]int{{64, 36}, {65, 37}, {9, 9}, {1, 1}, {2, 3}} {
		w, h := size[0], size[1]
		v := New(w, h, Config{Color: canvas.RGB(255, 255, 255), Alpha: 255, From: 0.1})
		cv, _ := canvas.NewCanvas(w, h)
		cv.Clear(canvas.RGB(0, 0, 0))
		v.Paint(cv, 1)
		for y := range h {
			for x := range w {
				got := int(cv.At(x, y).R())
				if d := got - int(v.AlphaAt(x, y)); d < -1 || d > 1 {
					t.Fatalf("%dx%d at %d,%d: %d, alpha %d", w, h, x, y, got, v.AlphaAt(x, y))
				}
			}
		}
	}
}

// A canvas of another size is left alone, not read out of its bounds.
func TestAnotherSizeIsLeftAlone(t *testing.T) {
	v := New(64, 36, red)
	cv, _ := canvas.NewCanvas(32, 18)
	cv.Clear(canvas.RGB(1, 2, 3))
	v.Paint(cv, 1)
	if cv.At(0, 0) != canvas.RGB(1, 2, 3) {
		t.Error("painted another size")
	}
}

func BenchmarkPaint(b *testing.B) {
	v := New(1280, 720, red)
	cv, _ := canvas.NewCanvas(1280, 720)
	cv.Clear(canvas.RGB(90, 140, 70))
	for b.Loop() {
		v.Paint(cv, 0.5)
	}
}
