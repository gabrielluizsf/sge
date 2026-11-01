package gl

import (
	"testing"

	"github.com/gabrielluizsf/sge/gltf"
	"github.com/gabrielluizsf/sge/pose"
	"github.com/gabrielluizsf/sge/render3d"
	"github.com/gabrielluizsf/sge/render3d/internal/shapes"
)

// offscreen is an offscreen device w by h, or the test skipped where
// there is no GL 4.2 to make one with.
func offscreen(tb testing.TB, w, h int) *Device {
	tb.Helper()
	d, err := Offscreen(w, h)
	if err != nil {
		tb.Skip("no OpenGL 4.2 here:", err)
	}
	tb.Cleanup(d.Close)
	return d
}

// scene is a scene of render3d/soft's tests: its model, and the scene
// that shows it, uploaded.
type scene struct {
	name  string
	model *gltf.Model
	show  func(m *render3d.Model) *render3d.Scene
}

// sky is the clear colour the scenes are drawn over: sRGB 128 grey.
var sky = [4]float32{0.2158605, 0.2158605, 0.2158605, 1}

// lit is a scene seen by cam, lit from the front, down the game's +Y,
// with no ambient light.
func lit(cam render3d.Camera, in *render3d.Instance) *render3d.Scene {
	return &render3d.Scene{Camera: cam, Sun: render3d.Sun{Toward: [3]float64{0, -1, 0}, Color: [3]float64{1, 1, 1}},
		Instances: []*render3d.Instance{in}, Clear: sky}
}

// cube is a textured cube in perspective, its face toward the camera
// red on the left and blue on the right.
func cube(m *render3d.Model) *render3d.Scene {
	cam := render3d.Perspective(render3d.PerspectiveDesc{Eye: [3]float64{0, -5, 0}, Up: [3]float64{0, 0, 1}, FovY: 40, Aspect: 1, Near: 0.1, Far: 20})
	return lit(cam, &render3d.Instance{Model: render3d.One(m), Place: render3d.At([3]float64{0, 0, 0}, 30, 1)})
}

// arm is a two-bone arm, sixteen pixels a unit, played anim one second.
func arm(anim string) func(m *render3d.Model) *render3d.Scene {
	return func(m *render3d.Model) *render3d.Scene {
		cam := render3d.Orthographic(render3d.OrthographicDesc{Eye: [3]float64{0, -5, 1}, Target: [3]float64{0, 0, 1}, Up: [3]float64{0, 0, 1}, Width: 4, Height: 4, Near: 0.1, Far: 10})
		return lit(cam, &render3d.Instance{Model: render3d.One(m), Place: pose.Identity, Animation: anim, Time: 1})
	}
}

// shade is a block over the ground, seen from above, the sun to its side
// casting its shadow by a map.
func shade(m *render3d.Model) *render3d.Scene {
	ground := &render3d.Instance{Model: render3d.One(m), Place: pose.TRS([3]float64{0, 0, -0.05}, shapes.Still, [3]float64{8, 8, 0.1})}
	block := &render3d.Instance{Model: render3d.One(m), Place: render3d.At([3]float64{0, 0, 1}, 0, 1)}
	cam := render3d.Orthographic(render3d.OrthographicDesc{Eye: [3]float64{0, 0, 10}, Up: [3]float64{0, 1, 0}, Width: 8, Height: 8, Near: 0.1, Far: 20})
	return &render3d.Scene{
		Camera:    cam,
		Instances: []*render3d.Instance{ground, block},
		Sun: render3d.Sun{Toward: [3]float64{1, 0, 1}, Color: [3]float64{0.8, 0.8, 0.8}, Ambient: [3]float64{0.2, 0.2, 0.2},
			Shadow: 64, Area: render3d.Sphere{Radius: 4}, Bias: 0.01},
	}
}

// scenes are the scenes soft's tests draw.
func scenes() []scene {
	return []scene{
		{"textured cube", shapes.Cube(), cube},
		{"arm at rest", shapes.Arm(), arm("")},
		{"arm bent", shapes.Arm(), arm("bend")},
		{"sun's shadow", shapes.Plain(), shade},
	}
}
