package main

import (
	"testing"

	"github.com/gabrielluizsf/sge/farm"
	"github.com/gabrielluizsf/sge/human"
)

// The game holds a climb as long as its motion lasts, and stands one down where
// its last picture has them.
func TestTheGameKeepsToTheClimbsMotion(t *testing.T) {
	if farm.MountSeconds != human.MountLength {
		t.Errorf("the game takes %v seconds, the motion %v", farm.MountSeconds, human.MountLength)
	}
	if farm.DismountSeconds != human.DismountLength {
		t.Errorf("the game takes %v seconds down, the motion %v", farm.DismountSeconds, human.DismountLength)
	}
	if farm.MountAside != human.MountAside {
		t.Errorf("the game sets one down %v cells away, the motion %v metres", farm.MountAside, human.MountAside)
	}
}
