package pose

import (
	"math"
	"testing"

	"github.com/gabrielluizsf/sge/gltf"
)

// arm is two bones: a shoulder at the origin and an elbow one up its Y,
// the elbow turned a quarter round Z one second into "bend".
func arm() *gltf.Model {
	still := [4]float64{0, 0, 0, 1}
	one := [3]float64{1, 1, 1}
	return &gltf.Model{
		Roots: []int{0},
		Nodes: []gltf.Node{
			{Name: "shoulder", Children: []int{1}, Mesh: -1, Skin: -1, Rotation: still, Scale: one},
			{Name: "elbow", Mesh: -1, Skin: -1, Translation: [3]float64{0, 1, 0}, Rotation: still, Scale: one},
		},
		Skins: []gltf.Skin{{Joints: []int{0, 1}, InverseBind: [][16]float64{
			Identity, {0: 1, 5: 1, 10: 1, 13: -1, 15: 1},
		}}},
		Animations: []gltf.Animation{{Name: "bend", Channels: []gltf.Channel{{
			Node: 1, Path: "rotation", Interpolation: "LINEAR",
			Times: []float32{0, 2}, Values: [][4]float32{{0, 0, 0, 1}, {0, 0, 1, 0}},
		}}}},
	}
}

func near(a, b Mat) bool {
	for i := range a {
		if math.Abs(a[i]-b[i]) > 1e-6 {
			return false
		}
	}
	return true
}

// Half way through a turn of half a circle, the elbow has turned a
// quarter: its joint matrix turns a point at the elbow's bind place, one
// up, about the elbow, which a base lifting the model by 5 then moves.
func TestJointsMatchAPoseWorkedByHand(t *testing.T) {
	p := New(arm())
	anim := p.Animation("bend")
	if anim != 0 || p.Animation("run") != -1 {
		t.Fatalf("animations %d", anim)
	}
	base := TRS([3]float64{0, 0, 5}, [4]float64{0, 0, 0, 1}, [3]float64{1, 1, 1})
	p.At(anim, 1, base)
	got := p.Joints(0, nil)
	// shoulder: just the base. elbow: base, up 1, a quarter about Z, down 1.
	c, s := math.Cos(math.Pi/2), math.Sin(math.Pi/2)
	elbow := Mat{c, s, 0, 0, -s, c, 0, 0, 0, 0, 1, 0, 0, 1, 5, 1}.Mul(Mat{0: 1, 5: 1, 10: 1, 13: -1, 15: 1})
	if !near(got[0], base) || !near(got[1], elbow) {
		t.Fatalf("joints\n%v\n%v\nwant\n%v\n%v", got[0], got[1], base, elbow)
	}
	// The forearm's tip, two up in the bind pose, ends one to the left of
	// the elbow, at (-1, 1, 5).
	tip := got[1].Point([3]float64{0, 2, 0})
	if math.Abs(tip[0]+1) > 1e-9 || math.Abs(tip[1]-1) > 1e-9 || math.Abs(tip[2]-5) > 1e-9 {
		t.Errorf("tip at %v", tip)
	}
	// Looped: 3 seconds in is 1 second in.
	p.At(anim, 3, base)
	if again := p.Joints(0, nil); !near(again[1], elbow) {
		t.Errorf("not looped: %v", again[1])
	}
}

// Affine is Mul for matrices whose last row is 0, 0, 0, 1.
func TestAffineIsMul(t *testing.T) {
	a := TRS([3]float64{1, 2, 3}, [4]float64{0.1, 0.2, 0.3, 0.927}, [3]float64{2, 1, 0.5})
	b := TRS([3]float64{-4, 0, 1}, [4]float64{0.5, 0.5, 0.5, 0.5}, [3]float64{1, 3, 1})
	if !near(a.Affine(b), a.Mul(b)) {
		t.Errorf("%v != %v", a.Affine(b), a.Mul(b))
	}
}
