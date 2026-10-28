package main

import (
	"github.com/gabrielluizsf/sge/pathfind"
	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/vector"
)

// The wild's layers.
const (
	wildGround = iota
	wildWater
	wildThings
)

// onFarm reports whether cell col, row is the farm's, not the wild's.
func (f *farm) onFarm(col, row int) bool {
	return col >= 0 && row >= 0 && col < f.ground.Cols && row < f.ground.Rows
}

// cellOf is the cell under the world point p, on the farm or off it.
func (f *farm) cellOf(p vector.Vector3) (col, row int) { return f.wild.CellOf(p.X, p.Y) }

// center is the world point in the middle of cell c.
func (f *farm) center(c pathfind.Cell) vector.Vector3 { return f.ground.CellCenter(c.Col, c.Row) }

// groundAt is the ground's tile on cell col, row, and setGround sets it.
func (f *farm) groundAt(col, row int) int {
	if f.onFarm(col, row) {
		return f.ground.At(col, row)
	}
	return f.wild.Gen.Tile(wildGround, col, row)
}

func (f *farm) setGround(col, row, id int) {
	if f.onFarm(col, row) {
		f.ground.Set(col, row, id)
		return
	}
	f.wild.Set(wildGround, col, row, id)
}

// thingAt is what stands on cell col, row, and setThing sets it.
func (f *farm) thingAt(col, row int) int {
	if f.onFarm(col, row) {
		return f.things.At(col, row)
	}
	return f.wild.Gen.Tile(wildThings, col, row)
}

func (f *farm) setThing(col, row, id int) {
	if f.onFarm(col, row) {
		f.things.Set(col, row, id)
		return
	}
	f.wild.Set(wildThings, col, row, id)
}

// setCrop shows tile id growing on cell col, row: on the farm's crops,
// or among the wild's things, where a garden has no tree.
func (f *farm) setCrop(col, row, id int) {
	if f.onFarm(col, row) {
		f.crops.Set(col, row, id)
		return
	}
	f.wild.Set(wildThings, col, row, id)
}

// open reports whether one may walk across cell col, row without
// swimming: nothing solid stands there and it is dry.
func (f *farm) open(col, row int) bool {
	if f.tiles.IsSolid(f.thingAt(col, row)) {
		return false
	}
	if f.onFarm(col, row) {
		return f.water.At(col, row) == tilemap.Empty
	}
	return !f.wild.Gen.In("water", col, row)
}

// cellKey is cell col, row as one number, a key JavaScript maps compare
// fast.
func cellKey(col, row int) int {
	const reach = 1 << 20
	return (col+reach)*(2*reach) + row + reach
}
