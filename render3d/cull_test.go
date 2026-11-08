package render3d

import (
	"fmt"
	"math"
	"testing"

	"github.com/gabrielluizsf/sge/iso"
	"github.com/gabrielluizsf/sge/render3d/internal/shapes"
	"github.com/gabrielluizsf/sge/vector"
)

// flat is a camera looking down the game's +Y, 16 pixels a unit on a
// screen 64 square.
var flat = Orthographic(OrthographicDesc{Eye: [3]float64{0, -5, 0}, Up: [3]float64{0, 0, 1}, Width: 4, Height: 4, Near: 0.1, Far: 10})

// cubes is a renderer on a null screen 64 square, and a cube uploaded
// to it.
func cubes(t testing.TB) (*Renderer, *Model) {
	r := NewRenderer(&null{w: 64, h: 64})
	m, err := r.Upload(shapes.Cube())
	if err != nil {
		t.Fatal(err)
	}
	return r, m
}

// Instances out of the camera's sight, to a side or behind it, are not
// drawn; one straddling the screen's edge is, with the one in the
// middle, in one draw.
func TestInstancesOutOfSightAreCulled(t *testing.T) {
	r, m := cubes(t)
	at := func(x, y float64) *Instance { return &Instance{Model: One(m), Place: At([3]float64{x, y, 0}, 0, 1)} }
	s := &Scene{Camera: flat, Instances: []*Instance{at(0, 0), at(2.3, 0), at(10, 0), at(0, -9)}}
	if err := r.Draw(s); err != nil {
		t.Fatal(err)
	}
	if st := r.Stats; st.Shown != 2 || st.Culled != 2 || st.Draws != 1 || st.Triangles != 24 {
		t.Errorf("stats %+v", st)
	}
	if d := r.dev.(*null).draws; d != 1 {
		t.Errorf("%d draws on the device", d)
	}
}

// An instance shrinking and growing past its levels' thresholds changes
// level only once past them by the hysteresis.
func TestLevelsOfDetailChangeWithHysteresis(t *testing.T) {
	r, fine := cubes(t)
	coarse, err := r.Upload(shapes.Cube())
	if err != nil {
		t.Fatal(err)
	}
	lod := &LOD{Levels: []Level{{Model: fine, Pixels: 20}, {Model: coarse}}, Hysteresis: 0.2}
	in := &Instance{Model: lod}
	// The cube's ball is √3/2 round, 16 pixels a unit.
	across := 2 * math.Sqrt(3) / 2 * 16
	for i, step := range []struct {
		pixels float64
		level  int
	}{{30, 0}, {18, 0}, {15, 1}, {22, 1}, {25, 0}, {17, 0}, {10, 1}} {
		in.Place = At([3]float64{}, 0, step.pixels/across)
		r.Prepare(&Scene{Camera: flat, Instances: []*Instance{in}})
		if in.Level() != step.level {
			t.Errorf("step %d, %v pixels across: level %d, want %d", i, step.pixels, in.Level(), step.level)
		}
	}
}

// Opaque draws come before those cut by alpha, and nearer before
// farther within a material.
func TestDrawsAreSorted(t *testing.T) {
	r := NewRenderer(&null{w: 64, h: 64})
	cut := shapes.Plain()
	cut.Materials[0].AlphaMode, cut.Materials[0].AlphaCutoff = "MASK", 0.5
	leaves, err := r.Upload(cut)
	if err != nil {
		t.Fatal(err)
	}
	arm, err := r.Upload(shapes.Arm())
	if err != nil {
		t.Fatal(err)
	}
	var s Scene
	s.Camera = flat
	for _, y := range []float64{2, 0, 1} {
		s.Instances = append(s.Instances, &Instance{Model: One(leaves), Place: At([3]float64{0, y, 0}, 0, 0.5)}, &Instance{Model: One(arm), Place: At([3]float64{1, y, -1}, 0, 0.5)})
	}
	r.Prepare(&s)
	var got []string
	last := -1.0
	for i, d := range r.draws {
		kind := "arm"
		if d.model == leaves {
			kind = "leaves"
		}
		got = append(got, kind)
		if kind == "arm" && i > 0 && d.depth < last {
			t.Errorf("draw %d nearer than the one before", i)
		}
		last = d.depth
	}
	if want := "[arm arm arm leaves]"; fmt.Sprint(got) != want {
		t.Errorf("draws %v, want %s", got, want)
	}
}

// The isometric camera puts world points where the game's projection
// does, nearer ones, up the line of sight, less deep.
func TestTheIsoCameraMatchesTheProjection(t *testing.T) {
	p := iso.New(32, 64, 32)
	center := [3]float64{100, 50, 0}
	cam := Iso(IsoDesc{Projection: p, Center: center, Width: 640, Height: 480, Depth: 1000})
	pixel := func(w [3]float64) (x, y, z float64) {
		m := cam.ViewProj
		cx := m[0]*w[0] + m[4]*w[1] + m[8]*w[2] + m[12]
		cy := m[1]*w[0] + m[5]*w[1] + m[9]*w[2] + m[13]
		cz := m[2]*w[0] + m[6]*w[1] + m[10]*w[2] + m[14]
		return (cx + 1) / 2 * 640, (1 - cy) / 2 * 480, cz
	}
	o := p.Project(vector.Vector3{X: center[0], Y: center[1], Z: center[2]})
	for _, w := range [][3]float64{{100, 50, 0}, {164, 18, 0}, {90, 70, 40}} {
		want := p.Project(vector.Vector3{X: w[0], Y: w[1], Z: w[2]})
		x, y, _ := pixel(w)
		if math.Abs(x-(want.X-o.X+320)) > 1e-6 || math.Abs(y-(want.Y-o.Y+240)) > 1e-6 {
			t.Errorf("%v at %v, %v; want %v, %v", w, x, y, want.X-o.X+320, want.Y-o.Y+240)
		}
	}
	_, _, far := pixel([3]float64{100, 50, 0})
	_, _, nearer := pixel([3]float64{101, 51, 1})
	if math.Abs(far-0.5) > 1e-9 || nearer >= far {
		t.Errorf("depths %v, %v", far, nearer)
	}
}
