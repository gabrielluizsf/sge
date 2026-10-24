package terrain

import "github.com/gabrielluizsf/sge/game"

// New is an endless world of gen's areas in chunks of size cells, its
// layers drawn at z, one a layer; an object named name to add to a game.
func New(name string, gen *Generator, size int, z ...int) *World {
	w := &World{Object: game.NewObject(name), Gen: gen, Size: max(size, 1), Z: z, Margin: 1, chunks: map[int]*chunk{}}
	w.AddScript(stream{w})
	return w
}
