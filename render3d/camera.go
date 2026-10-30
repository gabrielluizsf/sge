package render3d

import (
	"math"

	"github.com/gabrielluizsf/sge/pose"
)

// Camera is where the world is seen from: ViewProj takes the world into
// clip space, and Eye is where the camera stands, for lighting.
type Camera struct {
	ViewProj pose.Mat
	Eye      [3]float64
}

// PerspectiveDesc is a camera at Eye looking at Target, Up its up, FovY
// degrees high and Aspect times as wide, seeing from Near to Far.
type PerspectiveDesc struct {
	Eye, Target, Up [3]float64
	FovY, Aspect    float64
	Near, Far       float64
}

// Perspective is the camera p describes, far things smaller.
func Perspective(p PerspectiveDesc) Camera {
	f := 1 / math.Tan(p.FovY*math.Pi/360)
	proj := pose.Mat{0: f / p.Aspect, 5: f, 10: p.Far / (p.Near - p.Far), 11: -1, 14: p.Near * p.Far / (p.Near - p.Far)}
	return Camera{ViewProj: proj.Mul(lookAt(p.Eye, p.Target, p.Up)), Eye: p.Eye}
}

// OrthographicDesc is a camera at Eye looking at Target, Up its up,
// seeing Width by Height of the world, from Near to Far.
type OrthographicDesc struct {
	Eye, Target, Up [3]float64
	Width, Height   float64
	Near, Far       float64
}

// Orthographic is the camera o describes, far things as large as near.
func Orthographic(o OrthographicDesc) Camera {
	proj := pose.Mat{0: 2 / o.Width, 5: 2 / o.Height, 10: 1 / (o.Near - o.Far), 14: o.Near / (o.Near - o.Far), 15: 1}
	return Camera{ViewProj: proj.Mul(lookAt(o.Eye, o.Target, o.Up)), Eye: o.Eye}
}

// lookAt is the view from eye toward target, up its up: the eye at the
// origin looking down its -Z, Y up.
func lookAt(eye, target, up [3]float64) pose.Mat {
	f := unit(sub(target, eye))
	s := unit(cross(f, up))
	u := cross(s, f)
	return pose.Mat{
		s[0], u[0], -f[0], 0,
		s[1], u[1], -f[1], 0,
		s[2], u[2], -f[2], 0,
		-dot(s, eye), -dot(u, eye), dot(f, eye), 1,
	}
}

func sub(a, b [3]float64) [3]float64 { return [3]float64{a[0] - b[0], a[1] - b[1], a[2] - b[2]} }

func dot(a, b [3]float64) float64 { return a[0]*b[0] + a[1]*b[1] + a[2]*b[2] }

func cross(a, b [3]float64) [3]float64 {
	return [3]float64{a[1]*b[2] - a[2]*b[1], a[2]*b[0] - a[0]*b[2], a[0]*b[1] - a[1]*b[0]}
}

func unit(v [3]float64) [3]float64 {
	l := math.Sqrt(dot(v, v))
	if l == 0 {
		return v
	}
	return [3]float64{v[0] / l, v[1] / l, v[2] / l}
}
