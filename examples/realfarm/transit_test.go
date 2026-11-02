package main

import (
	"testing"

	"github.com/gabrielluizsf/sge/examples/farm"
	"github.com/gabrielluizsf/sge/human"
)

// The baked dive and rise last as long as the farm shows them for, so the
// farmer is still while they play and no longer.
func TestDiveAndRiseLastAsLongAsTheFarmShowsThem(t *testing.T) {
	if human.Dive.Length != farm.DiveSeconds || human.Rise.Length != farm.RiseSeconds {
		t.Errorf("dive %v and rise %v seconds, not %v and %v", human.Dive.Length, human.Rise.Length, farm.DiveSeconds, farm.RiseSeconds)
	}
	if !human.Dive.Once || !human.Rise.Once {
		t.Error("a dive or a rise loops")
	}
}
