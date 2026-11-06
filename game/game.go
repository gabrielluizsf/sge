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
	// LazyReach, when above zero, is how many pixels beyond its position a
	// sprite that makes its image on first use may reach: one not made
	// yet, and standing further than that outside the canvas, is not
	// made, as it could show nothing, until it comes near. What is drawn
	// is the same; only the cost of pictures never seen is saved.
	LazyReach float64
	// Camera is the point of the world drawn at the canvas's top-left
	// corner; objects fixed to the Screen ignore it.
	Camera vector.Vector3
	// Projection, when set, is how the world is seen, as iso.New sees it
	// from a corner: the world stays a grid where objects move and
	// collide, and only what is drawn turns. Objects of one Z are then
	// drawn nearer over farther, by where they stand (Object.Footing).
	Projection Projector
	// Reveals, when set, tells each frame what to keep in sight in a
	// projected world: tall tiles and Fades objects nearer than one of
	// them, covering it, are drawn faded.
	Reveals func() []Reveal
	// Sound, when set, is the mixer Run plays to the speakers while the
	// window is open: the world's sounds. Where there are no speakers the
	// game runs in silence.
	Sound *sound.Mixer
	// Paused stops the world's time: scripts do not update, but for those
	// of objects that are Always, while animations go on moving and
	// everything is still drawn, as a game's menu holds it.
	Paused bool

	objects []*Object
	window  *antui.Window
	time    float64
	quit    bool

	viewW, viewH int
	items        []item   // what Draw draws, reused every frame
	covers       []Cover  // where the Reveals stand, for this frame
	last, next   frameKey // of the frame Changed last answered true for, and of this one
}
