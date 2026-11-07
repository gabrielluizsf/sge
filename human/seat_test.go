package human

import (
	"testing"

	"github.com/gabrielluizsf/sge/gltf"
)

// seatChair is a model of one node, to stand for a chair.
func seatChair() *gltf.Model {
	return &gltf.Model{Roots: []int{0}, Nodes: []gltf.Node{{Name: "chair", Mesh: -1, Skin: -1, Rotation: [4]float64{3: 1}, Scale: [3]float64{1, 1, 1}}}}
}

// seatPerson is a person of a root and two animations.
func seatPerson() *gltf.Model {
	return &gltf.Model{Roots: []int{0}, Nodes: []gltf.Node{{Name: "root", Mesh: -1, Skin: -1}}, Animations: []gltf.Animation{
		{Name: "stand", Channels: []gltf.Channel{{Node: 0, Times: []float32{0, 1}}}},
		{Name: "sit", Channels: []gltf.Channel{{Node: 0, Times: []float32{0, 1, 2}}}},
	}}
}

// The chair is out of sight but for the sitting, there at its size all
// through, in its place in the person's space and a child of their root.
func TestTheSeatShowsOnlyWhileSitting(t *testing.T) {
	m := seatPerson()
	err := Seat(m, seatChair(), SeatPlace{At: [3]float64{0, 0, 0.1}, Scale: 0.72, During: []string{"sit"}})
	if err != nil {
		t.Fatal(err)
	}
	at := m.Roots[0]
	kids := m.Nodes[at].Children
	if len(kids) != 1 {
		t.Fatalf("the root has %d children", len(kids))
	}
	n := m.Nodes[kids[0]]
	if n.Scale[0] != tiny || n.Translation[2] != 0.1 {
		t.Errorf("the seat at rest is %+v", n)
	}
	if len(m.Animations[0].Channels) != 1 {
		t.Error("the chair shows while standing")
	}
	sit := m.Animations[1].Channels
	if len(sit) != 2 || sit[1].Path != "scale" || sit[1].Node != kids[0] || len(sit[1].Values) != 3 {
		t.Fatalf("sitting: %+v", sit)
	}
	for _, v := range sit[1].Values {
		if v[0] != 0.72 || v[1] != 0.72 || v[2] != 0.72 {
			t.Errorf("size %v", v)
		}
	}
}

// Sitting in a motion the person does not have is an error, not a chair
// that never shows.
func TestASeatInAMissingMotionIsAnError(t *testing.T) {
	if err := Seat(seatPerson(), seatChair(), SeatPlace{Scale: 1, During: []string{"swim"}}); err == nil {
		t.Error("no error")
	}
}
