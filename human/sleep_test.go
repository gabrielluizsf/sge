package human

import (
	"math"
	"testing"

	"github.com/gabrielluizsf/sge/gltf"
)

// A sleeper lies flat and breathes: four poses, no two alike in the chest,
// the body turned onto its back and let down to the bed.
func TestTheSleeperLiesFlatAndBreathes(t *testing.T) {
	if len(Sleep.Poses) != 4 || Sleep.Name != "sleep" || Sleep.Once {
		t.Fatalf("sleep: %d poses, once %v", len(Sleep.Poses), Sleep.Once)
	}
	seen := map[float64]bool{}
	for _, p := range Sleep.Poses {
		if math.Abs(math.Abs(p.Pitch)-math.Pi/2) > 0.01 || p.Pitch >= 0 {
			t.Errorf("pitch %v is not on the back", p.Pitch)
		}
		seen[p.Lean] = true
	}
	if len(seen) < 3 {
		t.Errorf("the chest rises to only %d heights", len(seen))
	}
}

// The bed is a root of its own, so it stays flat while the body turns, and
// shows only in the sleeping.
func TestTheBedIsARootShownOnlyWhileSleeping(t *testing.T) {
	m := seatPerson()
	m.Animations[1].Name = "sleep"
	err := Lay(m, seatChair(), BedPlace{At: [3]float64{0, -0.4, 0}, Yaw: math.Pi / 2, Scale: 0.9, During: "sleep"})
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Roots) != 2 || len(m.Nodes[m.Roots[0]].Children) != 0 {
		t.Fatalf("roots %v, the person's children %v", m.Roots, m.Nodes[m.Roots[0]].Children)
	}
	n := m.Nodes[m.Roots[1]]
	if n.Scale[0] != tiny || n.Translation[1] != -0.4 || math.Abs(n.Rotation[1]-math.Sin(math.Pi/4)) > 1e-9 {
		t.Errorf("the bed at rest is %+v", n)
	}
	if len(m.Animations[0].Channels) != 1 {
		t.Error("the bed shows while standing")
	}
	if last := m.Animations[1].Channels[len(m.Animations[1].Channels)-1]; last.Path != "scale" || last.Values[0][0] != 0.9 {
		t.Errorf("sleeping: %+v", last)
	}
}

// Lying in a motion the person does not have is an error.
func TestABedInAMissingMotionIsAnError(t *testing.T) {
	if err := Lay(seatPerson(), &gltf.Model{Roots: []int{0}, Nodes: []gltf.Node{{Mesh: -1}}}, BedPlace{Scale: 1, During: "sleep"}); err == nil {
		t.Error("no error")
	}
}

// A couple's two sides of the double bed are the sleep under names of their
// own, so a bed is laid under each by the motion's name and the poses stay one.
func TestTheSidesOfTheDoubleBedAreTheSleepUnderOtherNames(t *testing.T) {
	if len(SleepingTogether) != 2 || SleepLeft.Name == SleepRight.Name || SleepLeft.Name == Sleep.Name {
		t.Fatalf("names %q %q %q", Sleep.Name, SleepLeft.Name, SleepRight.Name)
	}
	for _, m := range SleepingTogether {
		if len(m.Poses) != len(Sleep.Poses) || m.Length != Sleep.Length || m.Once != Sleep.Once || m.Poses[1] != Sleep.Poses[1] {
			t.Errorf("%s is not the sleep", m.Name)
		}
	}
}
