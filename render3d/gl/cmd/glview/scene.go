package main

import (
	"github.com/gabrielluizsf/sge/pose"
	"github.com/gabrielluizsf/sge/render3d"
	"github.com/gabrielluizsf/sge/render3d/internal/shapes"
)

// view is what glview shows: a textured cube spinning, an arm bending,
// on a ground, in the sun.
type view struct {
	scene     render3d.Scene
	cube, arm *render3d.Instance
}

// newView uploads glview's models to r and is the scene showing them.
func newView(r *render3d.Renderer) (*view, error) {
	cube, err := r.Upload(shapes.Cube())
	if err != nil {
		return nil, err
	}
	arm, err := r.Upload(shapes.Arm())
	if err != nil {
		return nil, err
	}
	plain, err := r.Upload(shapes.Plain())
	if err != nil {
		return nil, err
	}
	v := &view{
		cube: &render3d.Instance{Model: render3d.One(cube)},
		arm:  &render3d.Instance{Model: render3d.One(arm), Animation: "bend"},
	}
	ground := &render3d.Instance{Model: render3d.One(plain), Place: pose.TRS([3]float64{0, 0, -0.05}, shapes.Still, [3]float64{8, 8, 0.1})}
	v.scene = render3d.Scene{
		Instances: []*render3d.Instance{ground, v.cube, v.arm},
		Sun: render3d.Sun{Toward: [3]float64{0.6, -0.8, 1.2}, Color: [3]float64{1, 0.95, 0.85}, Ambient: [3]float64{0.25, 0.27, 0.3},
			Shadow: 2048, Area: render3d.Sphere{Radius: 4}, Bias: 0.003},
		Clear: [4]float32{0.25, 0.4, 0.6, 1},
	}
	return v, nil
}

// at sets the view t seconds in, seen on a screen w by h.
func (v *view) at(t float64, size [2]int) {
	v.cube.Place = render3d.At([3]float64{-1.2, 0, 1}, t*60, 1)
	v.arm.Place = render3d.At([3]float64{1.2, 0, 0}, 0, 1)
	v.arm.Time = t
	v.scene.Camera = render3d.Perspective(render3d.PerspectiveDesc{Eye: [3]float64{0, -6, 3}, Target: [3]float64{0, 0, 0.9}, Up: [3]float64{0, 0, 1},
		FovY: 40, Aspect: float64(size[0]) / float64(max(size[1], 1)), Near: 0.1, Far: 50})
}
