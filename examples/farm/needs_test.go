package main

import (
	"testing"

	"github.com/gabrielluizsf/sge/examples/farm/kit/bakes"
)

// Everything the farm draws is baked when the game starts: each piece of
// furniture's model, and the forge, the stumps and the stove package
// building builds, not only what stands outside.
func TestEverythingDrawnIsBaked(t *testing.T) {
	for name, f := range furnishings {
		if !needed(bakes.Task{Kind: "model", ID: f.id}) {
			t.Errorf("the %s, %s, is never baked", name, f.id)
		}
	}
	for _, id := range []string{"stall", "anvil", "well", "fence_straight", "forge", "stump", "stove"} {
		if !needed(bakes.Task{Kind: "building", ID: id}) {
			t.Errorf("building %s is never baked", id)
		}
	}
}
