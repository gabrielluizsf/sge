package terrain

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/vector"
)

// bigMeadow is meadow in 32 pixel tiles, as a game at twice 16 pixels
// has them, with the grass swaying.
func bigMeadow(b *testing.B) *Generator {
	sheet, _ := canvas.NewCanvas(32*4, 32)
	sheet.FillRect(canvas.Area{Width: 32, Height: 32}, grassInk)
	sheet.FillRect(canvas.Area{X: 32, Width: 32, Height: 32}, waterInk)
	sheet.FillRect(canvas.Area{X: 72, Y: 8, Width: 16, Height: 16}, treeInk)
	sheet.FillRect(canvas.Area{X: 96, Width: 32, Height: 32}, treeInk)
	s, _ := tilemap.NewTileSet(sheet, 32, 32)
	s.Animate(grass, tilemap.Animation{Frames: []int{grass, 3}, FPS: 2, Staggered: true})
	g := meadow(b)
	g.Tiles = s
	return g
}

// BenchmarkWorldFrame is a frame of walking through an endless world in a
// maximised window: a step, making the chunks walked into, and a draw.
func BenchmarkWorldFrame(b *testing.B) {
	g := game.New()
	g.Covered = true
	w := New("world", bigMeadow(b), 16, 0, 1, 2)
	g.Add(w.Object)
	cv, _ := canvas.NewCanvas(1920, 1080)
	g.Draw(cv)
	x := 0.0
	b.ResetTimer()
	for b.Loop() {
		x += 2
		g.Camera = vector.Vector3{X: x, Y: x / 3}
		g.Step(1.0 / 60)
		g.Draw(cv)
	}
}

// BenchmarkMakeChunk is one chunk of 16 by 16 cells made, every layer.
func BenchmarkMakeChunk(b *testing.B) {
	g := game.New()
	w := New("world", bigMeadow(b), 16, 0, 1, 2)
	g.Add(w.Object)
	i := 0
	for b.Loop() {
		i++
		key, _ := cellKey(i, 0)
		w.makeChunk(key, i, 0)
	}
}
