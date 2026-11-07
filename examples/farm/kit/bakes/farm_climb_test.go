package bakes

import (
	"testing"

	"github.com/gabrielluizsf/sge/human"
)

// The last frame of a motion run through once is its end, not its start: its
// moment must fall before the last key's, which is kept in single precision and
// can round down.
func TestALastFrameIsTheMotionsEnd(t *testing.T) {
	for _, mo := range []human.Motion{human.Mount, human.Dismount, human.Hoe, human.Water} {
		if !mo.Once {
			continue
		}
		if at := timeOf(mo, FramesOf(mo.Name)-1); at >= float64(float32(mo.Length)) {
			t.Errorf("%s: the last frame at %v wraps past the last key at %v", mo.Name, at, float32(mo.Length))
		}
	}
}

// A climb's shadow is whole on the ground and gone in the saddle, where the
// horse's is.
func TestAClimbsShadow(t *testing.T) {
	if climbShade(human.Mount.Name, 0) != 1 || climbShade(human.Mount.Name, 1) != 0 {
		t.Error("getting on does not end with no shadow of its own")
	}
	if climbShade(human.Dismount.Name, 0) != 0 || climbShade(human.Dismount.Name, 1) != 1 {
		t.Error("getting down does not end with a whole shadow")
	}
}
