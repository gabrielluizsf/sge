package game

import (
	"math"

	"github.com/gabrielluizsf/sge/vector"
)

// Follow is a script that points the camera at its object, keeping it in
// the middle of the view as it moves and as the window is resized. With a
// Lag of some seconds the camera eases after it, covering most of the way
// in that time, instead of moving with it.
type Follow struct {
	Lag float64
}

// Start points the camera at the object at once: a lagging camera eases
// after it as it moves, not from wherever the camera was.
func (f Follow) Start(o *Object) {
	o.Game().Camera = aim(o)
}

func (f Follow) Update(o *Object, dt float64) {
	g := o.Game()
	want := aim(o)
	if f.Lag <= 0 {
		g.Camera = want
		return
	}
	// Two thirds of the way in Lag seconds, at any frame rate.
	g.Camera = g.Camera.Lerp(want, 1-math.Pow(1.0/3, dt/f.Lag))
}

// aim is where the camera is when o is in the middle of the view.
func aim(o *Object) vector.Vector3 {
	w, h := o.Game().View()
	return o.Transform.Position.Sub(vector.Vector3{X: float64(w) / 2, Y: float64(h) / 2})
}
