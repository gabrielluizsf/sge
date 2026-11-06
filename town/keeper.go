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
	Cell float64
	Wake func(t *Town) []*game.Object
	// Ahead, when set, makes a town wake when its middle is within as far
	// as the camera's view reaches and Ahead cells more, if that is more
	// than Range, so a town is awake and its people made before they come
	// into view, whatever the window's size.
	Ahead int
	// Sleep, when set, is told of a town as its objects are let go, to
	// drop what only they used: Wake makes it anew if the town comes back.
	Sleep func(t *Town)
	// PerFrame is how many towns wake in a frame: more
	// wait for the next, so a day's walk across towns never holds one frame
	// up waking several. With none set, all that are near wake at once.
	PerFrame int
	awake    map[*Town][]*game.Object
	order    []*Town // awake, in the order woken, to let go the same way every run
}

func (k *Keeper) Update(o *game.Object, _ float64) {
	g := o.Game()
	col, row := int(math.Floor(g.Camera.X/k.Cell)), int(math.Floor(g.Camera.Y/k.Cell))
	if k.awake == nil {
		k.awake = map[*Town][]*game.Object{}
	}
	woke := 0
	reach := k.reach(g)
	for _, t := range k.Layout.Near(col, row, reach) {
		if _, ok := k.awake[t]; ok {
			continue
		}
		if k.PerFrame > 0 && woke >= k.PerFrame {
			break
		}
		woke++
		objs := k.Wake(t)
		for _, w := range objs {
			g.Add(w)
		}
		k.awake[t] = objs
		k.order = append(k.order, t)
	}
	far := reach * 3 / 2
	kept := k.order[:0]
	for _, t := range k.order {
		if abs(t.Center.Col-col) <= far && abs(t.Center.Row-row) <= far {
			kept = append(kept, t)
			continue
		}
		for _, w := range k.awake[t] {
			w.Destroy()
		}
		if k.Sleep != nil {
			k.Sleep(t)
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
