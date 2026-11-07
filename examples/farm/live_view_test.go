package main

import (
	"math"
	"testing"

	"github.com/gabrielluizsf/sge/render3d"
)

// ndc is the point (0, 0, z) in clip space, y from -1 (bottom) to 1 (top).
func ndcY(cam render3d.Camera, z float64) float64 {
	m := cam.ViewProj
	v := [4]float64{0, 0, z, 1}
	var out [4]float64
	for r := range 4 {
		for c := range 4 {
			out[r] += m[c*4+r] * v[c]
		}
	}
	return out[1] / out[3]
}

func settled(v *liveView) {
	for range 600 {
		v.Step(1.0 / 60)
	}
}

// The body is framed whole, feet and crown inside the picture, and the
// head close: the crown inside, the feet out.
func TestTheHeadAndTheBodyAreFramed(t *testing.T) {
	v := &liveView{height: 1.7}
	settled(v)
	body := v.Camera(0.8)
	if feet, crown := ndcY(body, 0), ndcY(body, 1.7); feet < -1 || crown > 1 || crown-feet < 1.4 {
		t.Errorf("body: feet %.2f, crown %.2f", feet, crown)
	}
	v.Frame(true)
	settled(v)
	head := v.Camera(0.8)
	if feet, crown := ndcY(head, 0), ndcY(head, 1.7); feet > -1 || crown > 1 || crown < 0.3 {
		t.Errorf("head: feet %.2f, crown %.2f", feet, crown)
	}
	if face := ndcY(head, 1.6); math.Abs(face) > 0.5 {
		t.Errorf("the face is at %.2f, not in the middle", face)
	}
}

// The framing moves smoothly: it is neither where it was nor where it is
// going a moment after it is told.
func TestTheFramingMovesSmoothly(t *testing.T) {
	v := &liveView{height: 1.7}
	v.Frame(true)
	v.Step(1.0 / 60)
	if v.head <= 0 || v.head >= 0.5 {
		t.Errorf("one frame in, the framing is at %.3f", v.head)
	}
	if !v.Step(1.0 / 60) {
		t.Error("not moving")
	}
	settled(v)
	if v.Step(1.0/60) || v.head != 1 {
		t.Errorf("settled at %v", v.head)
	}
}

// Turn turns the person, and nothing else does.
func TestOnlyTurnTurnsThePerson(t *testing.T) {
	v := &liveView{height: 1.7}
	v.Frame(true)
	settled(v)
	v.Camera(1)
	if v.yaw != 0 {
		t.Errorf("yaw %v without a turn", v.yaw)
	}
	v.Turn(370)
	if v.yaw != 10 {
		t.Errorf("yaw %v", v.yaw)
	}
}
