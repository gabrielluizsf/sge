package kit

import (
	"slices"

	"github.com/gabrielluizsf/sge/human"
)

// smiths are the villagers that may be drawn as a smith, which villagerOf
// picks from by trade and gender.
var smiths = []string{"beto", "mei", "iara"}

// Motions are the motions a grown-up villager has: everyone's, and the
// smiths' strokes at the anvil.
func (v Villager) Motions() []human.Motion {
	if slices.Contains(smiths, v.ID) {
		return append(slices.Clone(human.Motions), human.Crafts...)
	}
	return human.Motions
}
