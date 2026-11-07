package kit

import (
	"slices"

	"github.com/gabrielluizsf/sge/human"
)

// smiths are the villagers that may be drawn as a smith, which villagerOf
// picks from by trade and gender.
var smiths = []string{"beto", "mei", "iara"}

// farmers are the villagers that may be drawn as a farmer, or are the
// player: those who water a plot.
var farmers = []string{"joao", "ana", "iara", "player"}

// Motions are the motions a grown-up villager has: everyone's, the
// smiths' strokes at the anvil and the farmers' pours from the can.
func (v Villager) Motions() []human.Motion {
	out := slices.Clone(human.Motions)
	if slices.Contains(smiths, v.ID) {
		out = append(out, human.Crafts...)
	}
	if slices.Contains(farmers, v.ID) {
		out = append(out, human.Farming...)
	}
	out = append(out, human.Mining...) // everyone grown swings the pickaxe at the rock
	if v.ID != "player" {
		out = append(out, human.Sleeping...)         // the others sleep in their beds
		out = append(out, human.SleepingTogether...) // and a couple in a double one
	}
	if v.ID == "player" {
		out = append(out, human.Sitting...) // only the player sits on a chair
		out = append(out, human.Eating...)  // and eats
		out = append(out, human.Drawing...) // and fills the can at a well or a river
	}
	return out
}
