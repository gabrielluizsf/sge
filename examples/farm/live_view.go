package main

import (
	"math"

	"github.com/gabrielluizsf/sge/render3d"
)

// liveView is where the person is seen from: a character creator's
// camera, in front of them at the height of the eyes, mildly in
// perspective, framing the whole body or, smoothly, the head. The person
// turns only when Turn says.
type liveView struct {
	yaw    float64 // degrees, counterclockwise seen from above
	head   float64 // 0 frames the body, 1 the head
	goal   float64 // what head moves to
	height float64 // the person's height, metres
}

// The camera's measures.
const (
	viewFov      = 24.0 // degrees high: mild perspective
	headShare    = 0.25 // how much of the height the head's frame shows
	bodyShare    = 1.16 // how much of the height the body's frame shows
	eyeShare     = 0.935
	headCentre   = 0.915 // the head's middle, as a share of the height
	bodyCentre   = 0.5
	bodyEyeShare = 0.72 // where the camera stands for the body: low, so it looks level
	frameRate    = 7.0  // how fast the framing moves, 1/s
	turnPerPixel = 0.5  // degrees a pixel of drag turns the person
)

// Turn turns the person by deg degrees.
func (v *liveView) Turn(deg float64) { v.yaw = math.Mod(v.yaw+deg, 360) }

// Frame sets what the camera moves to frame: the head or the whole body.
func (v *liveView) Frame(head bool) {
	v.goal = 0
	if head {
		v.goal = 1
	}
}

// Step moves the framing dt seconds on, and reports whether it is still
// moving.
func (v *liveView) Step(dt float64) bool {
	d := v.goal - v.head
	if math.Abs(d) < 1e-3 {
		moved := v.head != v.goal
		v.head = v.goal
		return moved
	}
	v.head += d * (1 - math.Exp(-dt*frameRate))
	return true
}

// Camera is the camera that sees the person from the front on a screen
// aspect times as wide as high.
func (v *liveView) Camera(aspect float64) render3d.Camera {
	h := max(v.height, 0.5)
	t := v.head
	lerp := func(a, b float64) float64 { return a + (b-a)*t }
	shown := h * lerp(bodyShare, headShare)
	dist := shown / 2 / math.Tan(viewFov*math.Pi/360)
	return render3d.Perspective(render3d.PerspectiveDesc{
		Eye:    [3]float64{0, -dist, h * lerp(bodyEyeShare, eyeShare)},
		Target: [3]float64{0, 0, h * lerp(bodyCentre, headCentre)},
		Up:     [3]float64{0, 0, 1}, FovY: viewFov, Aspect: aspect, Near: 0.1, Far: dist + 10,
	})
}
