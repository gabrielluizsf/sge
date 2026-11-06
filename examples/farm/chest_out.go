package farm

import (
	"github.com/gabrielluizsf/sge/collider"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/sprite"
)

// yardChest is a chest standing out of doors, on a cell of the farm or of
// the wild, with its object in the world. Chests are kept here, not in
// the tiles of the wild, so they stay wherever they were put whether that
// part of the world is on screen or not.
type yardChest struct {
	p   piece
	obj *game.Object
}

// chestAt is the chest out of doors on cell col, row, or nil.
func (f *farm) chestAt(col, row int) *yardChest {
	for _, y := range f.yard {
		if y.p.col == col && y.p.row == row {
			return y
		}
	}
	return nil
}

// putChest stands the chest p out of doors, its object added to the game
// if there is one on screen.
func (f *farm) putChest(p piece) *yardChest {
	y := &yardChest{p: p, obj: f.newChestObject(p)}
	f.yard = append(f.yard, y)
	if g := f.player.Game(); g != nil && f.room == nil {
		g.Add(y.obj)
	}
	return y
}

// takeChest lifts the chest y off the ground.
func (f *farm) takeChest(y *yardChest) {
	y.obj.Destroy()
	for i, c := range f.yard {
		if c == y {
			f.yard = append(f.yard[:i], f.yard[i+1:]...)
			return
		}
	}
}

// newChestObject is the object of chest p: its picture, turned as it is,
// standing on its cell among the people and the trees.
func (f *farm) newChestObject(p piece) *game.Object {
	o := game.NewObject("chest")
	o.Z = zThings
	img := pieceFull(f.art, p.name, p.heading)
	// A tile stands by its bottom middle on its cell's bottom corner, a
	// quarter tile below the middle the object stands on.
	o.Sprite = &sprite.Sprite{Image: img, PivotX: 0.5, PivotY: float64(img.Height-isoW/4) / float64(img.Height)}
	o.Transform.Position = f.ground.CellCenter(p.col, p.row)
	return o
}

// chestWalls is the chests out of doors as walls: nobody walks through.
type chestWalls struct{ f *farm }

// Blocks reports whether c touches a chest.
func (w chestWalls) Blocks(c collider.Collider) bool {
	for _, y := range w.f.yard {
		at := w.f.ground.CellCenter(y.p.col, y.p.row)
		box := collider.Box{X: at.X - tile/2, Y: at.Y - tile/2, Width: tile, Height: tile}
		if collider.Collides(c, box) {
			return true
		}
	}
	return false
}

// Stride is half the chests' thickness.
func (chestWalls) Stride() float64 { return tile / 2 }
