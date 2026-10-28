package town

import (
	"math"

	"github.com/gabrielluizsf/sge/game"
)

// Keeper is the script that keeps the towns near the camera awake: when
// a town comes within Range cells of the camera, Wake makes its objects,
// its people say, and adds them to the game; when it is half as far again
// away, they are destroyed, and Wake makes them anew if it comes back. The
// game keeps what must last between, as Wake's objects are let go.
type Keeper struct {
	Layout *Layout
	// Range is how many cells from the camera a town wakes.
	Range int
	// Cell is the side of a cell in the world's units.
	Cell  float64
	Wake  func(t *Town) []*game.Object
	awake map[*Town][]*game.Object
	order []*Town // awake, in the order woken, to let go the same way every run
}

func (k *Keeper) Update(o *game.Object, _ float64) {
	g := o.Game()
	col, row := int(math.Floor(g.Camera.X/k.Cell)), int(math.Floor(g.Camera.Y/k.Cell))
	if k.awake == nil {
		k.awake = map[*Town][]*game.Object{}
	}
	for _, t := range k.Layout.Near(col, row, k.Range) {
		if _, ok := k.awake[t]; ok {
			continue
		}
		objs := k.Wake(t)
		for _, w := range objs {
			g.Add(w)
		}
		k.awake[t] = objs
		k.order = append(k.order, t)
	}
	far := k.Range * 3 / 2
	kept := k.order[:0]
	for _, t := range k.order {
		if abs(t.Center.Col-col) <= far && abs(t.Center.Row-row) <= far {
			kept = append(kept, t)
			continue
		}
		for _, w := range k.awake[t] {
			w.Destroy()
		}
		delete(k.awake, t)
	}
	k.order = kept
}

// Awake reports whether town t's objects are in the game.
func (k *Keeper) Awake(t *Town) bool {
	_, ok := k.awake[t]
	return ok
}
