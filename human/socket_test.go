package human

import (
	"testing"

	"github.com/gabrielluizsf/sge/gltf"
)

// Every socket and every part of a ragdoll is on a bone a person has, each
// part hanging from one before it; a thing attached is seen in its animation.
func TestSocketsAndRagdollAreOnTheBones(t *testing.T) {
	k := realKit(t)
	m, err := k.Make(Body{Gender: 1, Age: 0.5, Muscle: 0.5, Weight: 0.5, Height: 0.5, Proportions: 0.5, European: 1}, Look{})
	if err != nil {
		t.Fatal(err)
	}
	if err := Animate(m, Stand, Sit); err != nil {
		t.Fatal(err)
	}
	for _, s := range Sockets {
		if boneOf(m, s.Bone) < 0 {
			t.Errorf("no bone %s", s.Bone)
		}
	}
	seen := map[string]bool{}
	for i, l := range Ragdoll {
		if boneOf(m, l.Bone) < 0 {
			t.Errorf("no bone %s", l.Bone)
		}
		if i > 0 && !seen[l.Parent] {
			t.Errorf("%s hangs from %q, not yet made", l.Bone, l.Parent)
		}
		seen[l.Bone] = true
	}
	nodes := len(m.Nodes)
	if err := Attach(m, sockThing(), SocketHead, Sit.Name); err != nil {
		t.Fatal(err)
	}
	if len(m.Nodes) <= nodes {
		t.Error("nothing was attached")
	}
	if err := Attach(m, sockThing(), Socket{Bone: "no such bone"}, Sit.Name); err == nil {
		t.Error("a socket on no bone was let by")
	}
}

// sockThing is a small thing to attach.
func sockThing() *gltf.Model {
	return &gltf.Model{Nodes: []gltf.Node{{Name: "thing", Mesh: -1, Skin: -1, Rotation: [4]float64{3: 1}, Scale: [3]float64{1, 1, 1}}}, Roots: []int{0}}
}

// Every named face is made of units the kit has, and turns the face.
func TestTheNamedFacesAreMadeOfRealUnits(t *testing.T) {
	k := realKit(t)
	m, err := k.Make(Body{Gender: 0, Age: 0.5, Muscle: 0.5, Weight: 0.5, Height: 0.5, Proportions: 0.5, African: 1}, Look{})
	if err != nil {
		t.Fatal(err)
	}
	for name, face := range map[string]func(float64) Expression{"smile": Smile, "sad": Sad, "angry": Angry, "surprised": Surprised, "pain": Pain, "wink": Wink} {
		e, err := k.Express(m, face(1))
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		moved := false
		for i := range e.Nodes {
			if e.Nodes[i].Rotation != m.Nodes[i].Rotation {
				moved = true
			}
		}
		if !moved {
			t.Errorf("%s moves nothing", name)
		}
	}
}

// Faces blend by the most of each unit.
func TestFacesBlend(t *testing.T) {
	b := Blend(Smile(1), Wink(1), Smile(0.2))
	if b["MouthLeftPullUp"] != 1 || b["LeftUpperLidClosed"] != 1 {
		t.Errorf("%v", b)
	}
	if len(Blend()) != 0 {
		t.Error("nothing blended is not nothing")
	}
}
