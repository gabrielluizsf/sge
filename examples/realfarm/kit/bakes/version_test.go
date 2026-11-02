package bakes

import (
	"testing"

	"github.com/gabrielluizsf/sge/human"
)

// Only the chop is baked again when it changes: every other motion keeps
// its name, and so what the shelf already holds of it.
func TestChopHasItsOwnVersion(t *testing.T) {
	if versionOf(human.Stand) != "v3" || versionOf(human.Walk) != "v3" {
		t.Errorf("the other motions' versions changed: %s, %s", versionOf(human.Stand), versionOf(human.Walk))
	}
	if versionOf(human.Chop) == versionOf(human.Stand) {
		t.Error("the chop's version is the others'")
	}
}

// The forge has its own version too, and it is not the chop's, so
// changing one bakes only that one again.
func TestForgeHasItsOwnVersion(t *testing.T) {
	if versionOf(human.Forge) == versionOf(human.Chop) || versionOf(human.Forge) == versionOf(human.Stand) {
		t.Errorf("the forge's version %s is another's", versionOf(human.Forge))
	}
}

// The hoe and the pour each have their own, apart from every other.
func TestHoeAndPourHaveTheirOwnVersions(t *testing.T) {
	seen := map[string]string{}
	for _, mo := range []human.Motion{human.Stand, human.Chop, human.Forge, human.Hoe, human.Water, human.Pick} {
		v := versionOf(mo)
		if other, ok := seen[v]; ok {
			t.Errorf("%s and %s are both %s", mo.Name, other, v)
		}
		seen[v] = mo.Name
	}
}
