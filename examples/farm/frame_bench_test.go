package farm

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

// windowed is a game warmed up on a 1280 by 720 canvas, the window the
// farm opens: stepped and drawn until the lazy tiles are painted and the
// pictures made, so a benchmark measures a frame in steady state.
func windowed(b *testing.B, ctl Controller) (*game.Game, *canvas.Canvas) {
	g, _ := newGame(ctl)
	cv, err := canvas.NewCanvas(1280, 720)
	if err != nil {
		b.Fatal(err)
	}
	for range 600 {
		g.Step(frame)
		g.Draw(cv)
	}
	b.ResetTimer()
	return g, cv
}

// BenchmarkFrame is a frame of the demo drawn into a window's canvas: a
// step of the world and every layer, the farmer and the HUD drawn, the
// bot tending the farm.
func BenchmarkFrame(b *testing.B) {
	g, cv := windowed(b, &bot{tends: 10})
	for range b.N {
		g.Step(frame)
		g.Draw(cv)
	}
}

// BenchmarkFrameIdle is a frame of the farm with the player standing
// still: only what sways moves.
func BenchmarkFrameIdle(b *testing.B) {
	g, cv := windowed(b, idle{})
	for range b.N {
		g.Step(frame)
		g.Draw(cv)
	}
}

// BenchmarkStepIdle is just the step of a frame of the farm with the player
// standing still: the clocks of the maps and the scripts, no drawing.
func BenchmarkStepIdle(b *testing.B) {
	g, _ := windowed(b, idle{})
	for range b.N {
		g.Step(frame)
	}
}

// BenchmarkFrameScroll is a frame of the farmer walking to and fro in the
// farm, the camera following, so the view scrolls a pixel or so a frame.
func BenchmarkFrameScroll(b *testing.B) {
	g, cv := windowed(b, idle{})
	o := g.Find("farmer")
	dx := 1.5
	for i := range b.N {
		if i%200 == 0 {
			dx = -dx
		}
		o.Transform.Position.X += dx
		g.Step(frame)
		g.Draw(cv)
	}
}

// BenchmarkWalkTheWild is a frame of the farmer walking off into the wild
// in a maximised window: the camera following, chunks of land made ahead
// of it, and everything drawn.
func BenchmarkWalkTheWild(b *testing.B) {
	g, _ := newGame(idle{})
	cv, err := canvas.NewCanvas(1920, 1080)
	if err != nil {
		b.Fatal(err)
	}
	g.Draw(cv)
	o := g.Find("farmer")
	o.Transform.Position = vector.Vector3{Y: -300}
	for b.Loop() {
		o.Transform.Position.X -= 2
		g.Step(frame)
		g.Draw(cv)
	}
}

// BenchmarkWildChunk is the tiles of a chunk of the wild worked out, every
// layer of its 16 by 16 cells, somewhere new each time.
func BenchmarkWildChunk(b *testing.B) {
	_, f := newGame(idle{})
	gen := f.wild.Gen
	i := 0
	for b.Loop() {
		i++
		for l := range gen.Layers {
			for r := range 16 {
				for c := range 16 {
					gen.Tile(l, i*16+c, 100+r)
				}
			}
		}
	}
}
