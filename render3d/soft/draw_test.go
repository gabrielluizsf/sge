package soft

import (
	"testing"

	"github.com/gabrielluizsf/sge/pose"
	"github.com/gabrielluizsf/sge/render3d"
	"github.com/gabrielluizsf/sge/render3d/internal/shapes"
)

// sky is the clear colour the tests draw over: sRGB 128 grey's light.
var sky = [4]float32{0.2158605, 0.2158605, 0.2158605, 1}

// front is a scene seen from the front, down the game's +Y, lit
// straight from the camera, with no ambient light.
func front(cam render3d.Camera, in ...*render3d.Instance) *render3d.Scene {
	return &render3d.Scene{
		Camera:    cam,
		Sun:       render3d.Sun{Toward: [3]float64{0, -1, 0}, Color: [3]float64{1, 1, 1}},
		Instances: in,
		Clear:     sky,
	}
}

// A textured cube, uploaded once and drawn in perspective, shows the
// face toward the camera: its left half red and its right half blue,
// fully lit; the sky round it.
func TestATexturedCubeIsDrawn(t *testing.T) {
	dev := New(64, 64)
	r := render3d.NewRenderer(dev)
	m, err := r.Upload(shapes.Cube())
	if err != nil {
		t.Fatal(err)
	}
	cam := render3d.Perspective(render3d.PerspectiveDesc{Eye: [3]float64{0, -5, 0}, Up: [3]float64{0, 0, 1}, FovY: 40, Aspect: 1, Near: 0.1, Far: 20})
	in := &render3d.Instance{Model: render3d.One(m), Place: pose.Identity}
	if err := r.Draw(front(cam, in)); err != nil {
		t.Fatal(err)
	}
	img, err := dev.Read(render3d.Screen)
	if err != nil {
		t.Fatal(err)
	}
	if c := img.NRGBAAt(27, 32); !shapes.Near(c, [3]int{255, 0, 0}) {
		t.Errorf("left of the face is %v", c)
	}
	if c := img.NRGBAAt(37, 32); !shapes.Near(c, [3]int{0, 0, 255}) {
		t.Errorf("right of the face is %v", c)
	}
	if c := img.NRGBAAt(2, 2); !shapes.Near(c, [3]int{128, 128, 128}) {
		t.Errorf("the sky is %v", c)
	}
	if s := r.Stats; s.Shown != 1 || s.Draws != 1 || s.Triangles != 12 {
		t.Errorf("stats %+v", s)
	}
}

// A bar skinned to two bones stands straight at rest; bent, its upper
// half lies to the left of the elbow, and nothing is above it.
func TestASkinnedArmBends(t *testing.T) {
	dev := New(64, 64)
	r := render3d.NewRenderer(dev)
	m, err := r.Upload(shapes.Arm())
	if err != nil {
		t.Fatal(err)
	}
	// 16 pixels a unit; the elbow, one up, at the middle of the screen.
	cam := render3d.Orthographic(render3d.OrthographicDesc{Eye: [3]float64{0, -5, 1}, Target: [3]float64{0, 0, 1}, Up: [3]float64{0, 0, 1}, Width: 4, Height: 4, Near: 0.1, Far: 10})
	in := &render3d.Instance{Model: render3d.One(m), Place: pose.Identity}
	white := [3]int{255, 255, 255}
	grey := [3]int{128, 128, 128}
	for _, c := range []struct {
		anim        string
		above, left [3]int
	}{{"", white, grey}, {"bend", grey, white}} {
		in.Animation, in.Time = c.anim, 1
		if err := r.Draw(front(cam, in)); err != nil {
			t.Fatal(err)
		}
		img, _ := dev.Read(render3d.Screen)
		if got := img.NRGBAAt(32, 24); !shapes.Near(got, c.above) {
			t.Errorf("%q: above the elbow is %v", c.anim, got)
		}
		if got := img.NRGBAAt(24, 32); !shapes.Near(got, c.left) {
			t.Errorf("%q: left of the elbow is %v", c.anim, got)
		}
		if got := img.NRGBAAt(32, 40); !shapes.Near(got, white) {
			t.Errorf("%q: below the elbow is %v", c.anim, got)
		}
	}
}

// A cube over the ground, the sun above it, shades the ground under it,
// and not the ground beside it.
func TestTheSunCastsAShadow(t *testing.T) {
	dev := New(64, 64)
	r := render3d.NewRenderer(dev)
	m, err := r.Upload(shapes.Plain())
	if err != nil {
		t.Fatal(err)
	}
	ground := &render3d.Instance{Model: render3d.One(m), Place: pose.TRS([3]float64{0, 0, -0.05}, shapes.Still, [3]float64{8, 8, 0.1})}
	block := &render3d.Instance{Model: render3d.One(m), Place: render3d.At([3]float64{0, 0, 1}, 0, 1)}
	// Seen from above, 8 pixels a unit.
	cam := render3d.Orthographic(render3d.OrthographicDesc{Eye: [3]float64{0, 0, 10}, Up: [3]float64{0, 1, 0}, Width: 8, Height: 8, Near: 0.1, Far: 20})
	s := &render3d.Scene{
		Camera:    cam,
		Instances: []*render3d.Instance{ground, block},
		Sun: render3d.Sun{Toward: [3]float64{0, 0, 1}, Color: [3]float64{0.8, 0.8, 0.8}, Ambient: [3]float64{0.2, 0.2, 0.2},
			Shadow: 64, Area: render3d.Sphere{Radius: 4}, Bias: 0.01},
	}
	// The sun a little to the side, so the block's shadow shows beside it.
	s.Sun.Toward = [3]float64{1, 0, 1}
	if err := r.Draw(s); err != nil {
		t.Fatal(err)
	}
	img, _ := dev.Read(render3d.Screen)
	// The block covers x -0.5 to 0.5; its shadow falls a unit and more
	// toward -X, past its left edge; the ground at +X is lit.
	shade, lit := img.NRGBAAt(32-12, 34), img.NRGBAAt(32+20, 34)
	if int(shade.R)+30 > int(lit.R) {
		t.Errorf("shade %v, lit %v", shade, lit)
	}
}
