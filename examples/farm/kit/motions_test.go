package kit

import (
	"slices"
	"testing"

	"github.com/gabrielluizsf/sge/human"
	"github.com/gabrielluizsf/sge/life"
)

// Every grown villager swings the pickaxe at the rock, not the axe, and
// every one that is not the player sleeps: a teen in a bed of their own,
// the grown and the elders alone or as a couple on a side of the double one.
func TestEveryVillagerMinesAndSleepsAtEveryStage(t *testing.T) {
	has := func(ms []human.Motion, m human.Motion) bool {
		return slices.ContainsFunc(ms, func(x human.Motion) bool { return x.Name == m.Name })
	}
	for _, v := range Villagers {
		for _, g := range v.Stages() {
			grown := g.Stage >= life.Adult
			if grown && !has(g.Motions, human.Mine) {
				t.Errorf("%s as a %s does not mine", v.ID, g.Stage)
			}
			if v.ID == "player" {
				continue
			}
			if g.Stage >= life.Child && !has(g.Motions, human.Sleep) {
				t.Errorf("%s as a %s does not sleep", v.ID, g.Stage)
			}
			if couple := has(g.Motions, human.SleepLeft) && has(g.Motions, human.SleepRight); couple != grown {
				t.Errorf("%s as a %s: sleeps as a couple %v", v.ID, g.Stage, couple)
			}
		}
	}
}
