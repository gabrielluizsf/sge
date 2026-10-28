// Package game runs a world of objects: each frame their scripts run,
// their animations move on, the ones that touch are told, and they are
// drawn in order onto an antui window, or any canvas.
package game

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/sound"
	"github.com/gabrielluizsf/sge/vector"
)

// Game is the world and, while it runs, the window it is shown in.
type Game struct {
	// Background is what the canvas is cleared to before each frame.
	Background canvas.Color
	// FixedStep, when above zero, is the seconds every frame of Run steps
	// the world by, instead of the time the frame took: the same run
	// every time, for physics and replays.
	FixedStep float64
	// Covered says the world draws over the whole canvas every frame, as
	// an endless terrain does, so the canvas is not cleared first.
	Covered bool
	// Camera is the point of the world drawn at the canvas's top-left
	// corner; objects fixed to the Screen ignore it.
	Camera vector.Vector3
	// Projection, when set, is how the world is seen, as iso.New sees it
	// from a corner: the world stays a grid where objects move and
	// collide, and only what is drawn turns. Objects of one Z are then
	// drawn nearer over farther, by where they stand (Object.Footing).
	Projection Projector
	// Sound, when set, is the mixer Run plays to the speakers while the
	// window is open: the world's sounds. Where there are no speakers the
	// game runs in silence.
	Sound *sound.Mixer

	objects []*Object
	window  *antui.Window
	time    float64
	quit    bool

	viewW, viewH int
	items        []item // what Draw draws, reused every frame
}
