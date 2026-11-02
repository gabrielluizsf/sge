package human

import (
	"image"
	"testing"

	"github.com/gabrielluizsf/sge/gltf"
)

// pick is a tool of one node, mesh, material and image.
func pick() *gltf.Model {
	return &gltf.Model{
		Roots:     []int{0},
		Nodes:     []gltf.Node{{Name: "pick", Mesh: 0, Skin: -1, Rotation: [4]float64{3: 1}, Scale: [3]float64{1, 1, 1}}},
		Meshes:    []gltf.Mesh{{Primitives: []gltf.Primitive{{Material: 0}}}},
		Materials: []gltf.Material{{Texture: 0, MetalRough: -1, Normal: -1, Occlusion: -1}},
		Images:    []image.Image{image.NewNRGBA(image.Rect(0, 0, 1, 1))},
	}
}

// What is held shows only in the animations it is held in, as a child of
// the person's root, and brings its own mesh, material and image along.
func TestHeldToolIsTheirsAndShowsWhenHeld(t *testing.T) {
	m := armed(t)
	if err := Animate(m, Chop, Stand); err != nil {
		t.Fatal(err)
	}
	meshes, mats, images := len(m.Meshes), len(m.Materials), len(m.Images)
	if err := Hold(m, pick(), Grip{Hands: [2]string{"wrist.L", "wrist.R"}, During: []string{"chop"}, End: -0.4, Facing: zAxis}); err != nil {
		t.Fatal(err)
	}
	grip := boneOf(m, "grip")
	if grip < 0 || childNamed(m, m.Roots[0], "grip") != grip {
		t.Fatal("the tool is not a child of the root")
	}
	if m.Nodes[grip].Scale[0] > 0.001 {
		t.Error("the tool shows in no animation")
	}
	tool := m.Nodes[grip].Children[0]
	if n := m.Nodes[tool]; n.Mesh != meshes || m.Meshes[n.Mesh].Primitives[0].Material != mats || m.Materials[mats].Texture != images {
		t.Errorf("the tool's mesh, material and image are not its own: %+v", n)
	}
	shown := map[string]bool{}
	for _, a := range m.Animations {
		for _, ch := range a.Channels {
			if ch.Node == grip && ch.Path == "scale" && ch.Values[0][0] == 1 {
				shown[a.Name] = true
			}
		}
	}
	if !shown["chop"] || shown["stand"] {
		t.Errorf("the tool shows in %v", shown)
	}
}

// The hand at the handle's end never lets go: in every pose the palm is
// on the handle's axis, where End says.
func TestHeldToolStaysInTheLeadHand(t *testing.T) {
	m := armed(t)
	if err := Animate(m, Chop); err != nil {
		t.Fatal(err)
	}
	g := Grip{Hands: [2]string{"wrist.L", "wrist.R"}, During: []string{"chop"}, End: -0.4, Facing: zAxis}
	a := &m.Animations[0]
	for k := range a.Channels[0].Times {
		at, q := handle(m, a, g, k)
		w := worldAt(m, a, k)
		n := boneOf(m, "wrist.L")
		hand := add(w[n].at, w[n].turn.turn(palm(m, n)))
		end := add(at, q.turn([3]float64{g.End, handleY, 0}))
		if d := sub(hand, end); dot3(d, d) > 1e-8 {
			t.Fatalf("pose %d: the palm is %v from where it holds", k, d)
		}
	}
}

// A tool cannot be put in a hand that is not there, nor shown in a motion
// that does not exist.
func TestHoldRefusesWhatIsNotThere(t *testing.T) {
	m := armed(t)
	if err := Animate(m, Chop); err != nil {
		t.Fatal(err)
	}
	if Hold(m, pick(), Grip{Hands: [2]string{"wrist.X", "wrist.R"}}) == nil {
		t.Error("held in a hand that is not there")
	}
	if Hold(m, pick(), Grip{Hands: [2]string{"wrist.L", "wrist.R"}, During: []string{"dance"}}) == nil {
		t.Error("shown in a motion that is not there")
	}
}

// A tool held in one hand goes where the hand does: its origin at Offset
// from the wrist, turned as the wrist and Turn say, in every pose.
func TestToolInOneHandFollowsTheWrist(t *testing.T) {
	m := armed(t)
	if err := Animate(m, Chop); err != nil {
		t.Fatal(err)
	}
	g := Grip{Hands: [2]string{"wrist.R", "wrist.R"}, During: []string{"chop"}, Offset: [3]float64{0.1, -0.2, 0.05}, Turn: [4]float64{0, 0.7071068, 0, 0.7071068}}
	a := &m.Animations[0]
	for k := range a.Channels[0].Times {
		at, q := handle(m, a, g, k)
		w := worldAt(m, a, k)[boneOf(m, "wrist.R")]
		if d := sub(at, add(w.at, w.turn.turn(g.Offset))); dot3(d, d) > 1e-12 {
			t.Fatalf("pose %d: the tool is %v from the wrist's offset", k, d)
		}
		x, want := q.turn([3]float64{1, 0, 0}), w.turn.turn(quat(g.Turn).turn([3]float64{1, 0, 0}))
		if d := sub(x, want); dot3(d, d) > 1e-12 {
			t.Fatalf("pose %d: the tool points %v, not %v", k, x, want)
		}
	}
}
