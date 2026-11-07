package farm

import (
	"github.com/gabrielluizsf/sge/random"
	"github.com/gabrielluizsf/sge/town"
)

// stock puts in the pen of town t's stable, the first time it wakes, the
// animals for sale: the same in any game, from the town's seed.
func (h *herd) stock(t *town.Town) {
	r := townPen(t)
	if h.stable[t] || r.cols == 0 {
		return
	}
	h.stable[t] = true
	p := &pen{kind: "sheep", r: r, town: t}
	h.pens = append(h.pens, p)
	luck := random.New(t.Seed + 29)
	for _, kind := range []string{"sheep", "sheep", "horse"} {
		c := h.add(kind, luck.Intn(animalKinds[kind]), p)
		c.wait = 1 + 3*luck.Float()
	}
}
