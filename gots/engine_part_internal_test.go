package gots

import (
	"slices"
	"testing"
)

// Every engine package is in exactly one part of a split, so a CI that
// splits the engine among jobs tests each package once.
func TestTheEngineSplitsWhole(t *testing.T) {
	t.Setenv("GOTS_ENGINE", "bake,water")
	some := enginePart()
	t.Setenv("GOTS_ENGINE", "-bake,water")
	rest := enginePart()
	if !slices.Equal(some, []string{"bake", "water"}) || len(some)+len(rest) != len(enginePackages) || slices.Contains(rest, "bake") {
		t.Errorf("split into %v and %d others of %d", some, len(rest), len(enginePackages))
	}
	t.Setenv("GOTS_ENGINE", "")
	if len(enginePart()) != len(enginePackages) {
		t.Error("no split does not test every package")
	}
}
