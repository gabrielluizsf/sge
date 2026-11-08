package terrain

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/warm"
)

// The tiles of a chunk are given to the warm queue as it is made, so a
// tile of a kind not seen yet is drawn a few a frame, not all at once in
// the frame that first shows it.
func TestAChunksTilesAreWarmed(t *testing.T) {
	made := 0
	tiles := tilemap.NewLazyTileSet(tilemap.Lazy{Count: 3, Width: 4, Height: 4, Make: func(int) *canvas.Canvas {
		made++
		cv, _ := canvas.NewCanvas(4, 4)
		return cv
	}})
	gen := &Generator{Tiles: tiles, Layers: 1}
	gen.Add(Area{Name: "ground", Layer: 0, In: Everywhere, Tile: func(c, r int) int { return (c + r) & 1 }})
	q := &warm.Queue{PerFrame: 1}
	w := New("world", gen, 8, 0)
	w.Warm = q
	g := game.New()
	g.Add(w.Object)
	g.Add(game.NewObject("warm")).AddScript(q)
	cv, _ := canvas.NewCanvas(8, 8) // the view is inside a chunk
	g.Draw(cv)
	g.Step(1.0 / 60)
	if q.Waiting() == 0 {
		t.Fatal("no tile was given to the queue")
	}
	for range 4 {
		g.Step(1.0 / 60)
	}
	if q.Waiting() != 0 || made != 2 {
		t.Fatalf("%d tiles made, %d still waiting", made, q.Waiting())
	}
}
