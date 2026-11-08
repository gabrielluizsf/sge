package particle

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

var testCfg = Config{Speed: 20, Lift: 60, Gravity: 200, Rest: 1, Spread: 0.2, Size: 2, Grow: 1, Colors: []canvas.Color{0xFFAA0000, 0xFF880000}}

func run(s *System, seconds float64) {
	for t := 0.0; t < seconds; t += 1.0 / 60 {
		s.Update(nil, 1.0/60)
	}
}

// Specks thrown up come down on the ground and lie there, not vanish in the air.
func TestSpecksLandBeforeTheyFade(t *testing.T) {
	s := New(8, testCfg, 1)
	s.Emit(vector.Vector3{X: 5, Y: 5, Z: 10}, vector.Vector3{}, 4)
	run(s, 0.9)
	if s.Live() != 4 {
		t.Fatalf("live %d, want 4 lying on the ground", s.Live())
	}
	for _, p := range s.specks[:4] {
		if p.z != 0 || p.vz != 0 {
			t.Fatalf("a speck is still in the air: %+v", p)
		}
	}
}

// Specks are gone once their rest is over.
func TestSpecksFadeAway(t *testing.T) {
	s := New(8, testCfg, 1)
	s.Emit(vector.Vector3{}, vector.Vector3{}, 5)
	run(s, 3)
	if s.Live() != 0 {
		t.Fatalf("live %d, want 0", s.Live())
	}
}

// A full pool reuses its oldest speck, and never holds more than its size.
func TestAFullPoolReusesTheOldest(t *testing.T) {
	s := New(3, testCfg, 1)
	s.Emit(vector.Vector3{}, vector.Vector3{}, 10)
	if s.Live() != 3 {
		t.Fatalf("live %d, want 3", s.Live())
	}
}

// The same seed throws the same specks.
func TestTheSameSeedThrowsTheSame(t *testing.T) {
	a, b := New(6, testCfg, 9), New(6, testCfg, 9)
	a.Emit(vector.Vector3{}, vector.Vector3{X: 1}, 6)
	b.Emit(vector.Vector3{}, vector.Vector3{X: 1}, 6)
	run(a, 0.3)
	run(b, 0.3)
	for i := range a.specks {
		if a.specks[i] != b.specks[i] {
			t.Fatalf("speck %d parted: %+v %+v", i, a.specks[i], b.specks[i])
		}
	}
}

// Nothing is allocated, whatever the system does.
func TestNothingIsAllocatedAfterNew(t *testing.T) {
	s := New(64, testCfg, 1)
	n := testing.AllocsPerRun(50, func() {
		s.Emit(vector.Vector3{}, vector.Vector3{X: 1}, 16)
		s.Update(nil, 1.0/60)
	})
	if n != 0 {
		t.Fatalf("%v allocations a frame", n)
	}
}

func BenchmarkEmitAndUpdate(b *testing.B) {
	s := New(256, testCfg, 1)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if i%30 == 0 {
			s.Emit(vector.Vector3{}, vector.Vector3{X: 1}, 24)
		}
		s.Update(nil, 1.0/60)
	}
}

// shown is a system with one speck thrown and a canvas it is drawn on.
func shown(t *testing.T) (*System, *game.Object, *canvas.Canvas) {
	t.Helper()
	g := game.New()
	o := game.NewObject("specks")
	g.Add(o)
	cfg := Config{Gravity: 100, Rest: 10, Size: 12, Colors: []canvas.Color{canvas.RGB(200, 0, 0)}}
	s := New(1, cfg, 1)
	s.Emit(vector.Vector3{X: 30, Y: 30, Z: 5}, vector.Vector3{}, 1)
	cv, err := canvas.NewCanvas(120, 120)
	if err != nil {
		t.Fatal(err)
	}
	return s, o, cv
}

// painted counts the pixels drawn on cv, and their box.
func painted(cv *canvas.Canvas) (n, x0, y0, x1, y1 int) {
	x0, y0 = cv.Width, cv.Height
	for y := range cv.Height {
		for x := range cv.Width {
			if cv.At(x, y).A() != 0 {
				n++
				x0, y0, x1, y1 = min(x0, x), min(y0, y), max(x1, x+1), max(y1, y+1)
			}
		}
	}
	return n, x0, y0, x1, y1
}

// In the air a speck is drawn by the system and not by Lying; on the
// ground the other way round: what lies there goes under what stands on it.
func TestAirAndGroundAreDrawnApart(t *testing.T) {
	s, o, cv := shown(t)
	Lying{System: s}.Paint(o, cv)
	if n, _, _, _, _ := painted(cv); n != 0 {
		t.Errorf("Lying drew %d pixels of a speck in the air", n)
	}
	s.Paint(o, cv)
	if n, _, _, _, _ := painted(cv); n == 0 {
		t.Error("the speck in the air was not drawn")
	}
	run(s, 2)
	s2, _ := canvas.NewCanvas(120, 120)
	s.Paint(o, s2)
	if n, _, _, _, _ := painted(s2); n != 0 {
		t.Errorf("the system drew %d pixels of a speck on the ground", n)
	}
	Lying{System: s}.Paint(o, s2)
	if n, _, _, _, _ := painted(s2); n == 0 {
		t.Error("the speck on the ground was not drawn")
	}
}

// A stain is round, not a box, and wider than tall: blood in square specks
// looked like a heap of bricks.
func TestAStainIsRoundAndFlat(t *testing.T) {
	s, o, cv := shown(t)
	run(s, 2)
	Lying{System: s}.Paint(o, cv)
	n, x0, y0, x1, y1 := painted(cv)
	w, h := x1-x0, y1-y0
	if w <= h {
		t.Errorf("the stain is %d by %d, not flat", w, h)
	}
	if n >= w*h*9/10 || n < w*h/2 {
		t.Errorf("%d pixels in a box of %d by %d: not round", n, w, h)
	}
}

// A stain holds its colour, then fades out before it is gone.
func TestAStainHoldsItsColourAndThenFades(t *testing.T) {
	p := speck{rest: 10}
	if a := p.opacity(); a != 255 {
		t.Errorf("new: %d", a)
	}
	p.lain = 6
	if a := p.opacity(); a != 255 {
		t.Errorf("at 60%%: %d", a)
	}
	p.lain = 8.5
	if a := p.opacity(); a <= 0 || a >= 255 {
		t.Errorf("at 85%%: %d", a)
	}
}
