package farm

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/vector"
)

// BenchmarkFrame is a frame of the demo drawn into a window's canvas: a
// step of the world and every layer, the farmer and the HUD drawn.
func BenchmarkFrame(b *testing.B) {
	g, f := newGame(&bot{tends: 10})
	cv, err := canvas.NewCanvas(f.ground.Cols*tile, f.ground.Rows*tile)
	if err != nil {
		b.Fatal(err)
	}
	for range 600 {
		g.Step(frame)
	}
	b.ResetTimer()
	for range b.N {
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
