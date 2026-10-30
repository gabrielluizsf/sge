package farm

import (
	"math"

	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

// bot plays the farm by itself, for the demo: it walks to a cell, faces
// a way and does something there, then picks the next thing to do. It
// walks across, then up or down, which the farm's layout keeps clear.
type bot struct {
	task *task
	// tends is how many plots it keeps, the first ones of the field.
	tends int
}

// task is a place to stand, a way to face there, and what to do then.
type task struct {
	to   vector.Vector3 // where the farmer's position should be
	face string
	do   func(f *farm, o *game.Object) bool // whether to use what is held
}

func (b *bot) control(f *farm, o *game.Object, _ float64) intent {
	if b.task == nil {
		b.task = b.next(f)
		if b.task == nil {
			return intent{}
		}
	}
	d := b.task.to.Sub(o.Transform.Position)
	switch {
	case math.Abs(d.X) > 0.5:
		return intent{dir: vector.Vector3{X: d.X}, upTo: math.Abs(d.X)}
	case math.Abs(d.Y) > 0.5:
		return intent{dir: vector.Vector3{Y: d.Y}, upTo: math.Abs(d.Y)}
	}
	t := b.task
	b.task = nil
	return intent{face: t.face, use: t.do(f, o)}
}

// standIn is where the farmer's position, their feet, is in the middle
// of cell col, row.
func standIn(f *farm, col, row int) vector.Vector3 {
	return f.ground.CellCenter(col, row)
}
