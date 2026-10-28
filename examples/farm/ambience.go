package main

import (
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/sound"
	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/vector"
)

// ambience is the sound of where the farmer is: leaves rustling as loud
// as there are trees about, water lapping as loud as there is water, and
// now and then a bird in one of the trees.
type ambience struct {
	f             *farm
	farmer        *game.Object
	leaves, water *sound.Loop
	listen, bird  float64 // seconds to the next look about, and bird
	near          around
	luck          uint32
}

// around is what is near the farmer: the trees, where they stand, and how
// many cells of water.
type around struct {
	trees []vector.Vector3
	water int
}

const (
	// hearing is how many cells away a tree or water is heard from.
	hearing = 5
	// listenEvery is how often, in seconds, the farmer's surroundings are
	// looked at again.
	listenEvery = 0.25
)

func (a *ambience) Update(_ *game.Object, dt float64) {
	if a.listen -= dt; a.listen <= 0 {
		a.listen = listenEvery
		a.near = around{}
		if a.f.room == nil {
			a.near = a.f.around(feet(a.farmer))
		}
		a.leaves.Fade(0.22 * min(1, float64(len(a.near.trees))/6))
		a.water.Fade(0.3 * min(1, float64(a.near.water)/8))
	}
	if a.bird -= dt; a.bird <= 0 {
		a.bird = 2.5 + 5.5*a.chance()
		if n := len(a.near.trees); n > 0 {
			at := a.near.trees[int(a.chance()*float64(n))%n]
			a.f.sound("chirp", &at)
		}
	}
}

// chance is the next of the ambience's random numbers, from 0 to 1.
func (a *ambience) chance() float64 {
	a.luck = a.luck*1664525 + 1013904223
	return float64(a.luck>>8) / float64(1<<24)
}

// around is what is within hearing of p: on the farm its things and
// water, beyond it the wild's.
func (f *farm) around(p vector.Vector3) around {
	var near around
	col, row := f.cellOf(p)
	for r := row - hearing; r <= row+hearing; r++ {
		for c := col - hearing; c <= col+hearing; c++ {
			tree, water := f.cellHas(c, r)
			if tree {
				near.trees = append(near.trees, f.ground.CellCenter(c, r))
			}
			if water {
				near.water++
			}
		}
	}
	return near
}

// cellHas reports whether cell col, row has a tree, and water.
func (f *farm) cellHas(col, row int) (tree, water bool) {
	if f.onFarm(col, row) {
		return f.things.At(col, row) == tileID["tree0"], f.water.At(col, row) != tilemap.Empty
	}
	return f.wild.TileAt(wildThings, col, row) == tileID["tree0"], f.wild.TileAt(wildWater, col, row) != tilemap.Empty
}
