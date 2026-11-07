package human

import (
	"math"
	"reflect"
	"testing"

	"github.com/gabrielluizsf/sge/gltf"
)

// womanWithHair is a person with long hair, animated in the motions.
func womanWithHair(t *testing.T, hair string) *gltf.Model {
	t.Helper()
	k := realKit(t)
	b := baselines[1]
	l := Look{Eyes: b.look.Eyes, Wear: []string{eyesProxy, hair}}
	m, err := k.Make(b.body, l)
	if err != nil {
		t.Fatal(err)
	}
	if err := Animate(m, Stand, Walk, Chop); err != nil {
		t.Fatal(err)
	}
	return m
}

func swayChannel(t *testing.T, m *gltf.Model, motion string) gltf.Channel {
	t.Helper()
	node := -1
	for i, n := range m.Nodes {
		if n.Name == swayBone {
			node = i
		}
	}
	if node < 0 {
		t.Fatal("no sway bone")
	}
	for _, a := range m.Animations {
		if a.Name != motion {
			continue
		}
		for _, c := range a.Channels {
			if c.Node == node {
				return c
			}
		}
	}
	t.Fatalf("no sway in %s", motion)
	return gltf.Channel{}
}

// The roots of the hair stay with the head, the ends hang free, and
// nothing but hair hangs from the sway bone.
func TestHairRootsAreFixedAndEndsFree(t *testing.T) {
	m := womanWithHair(t, "hair/long01/long01.mhclo")
	joint := len(m.Skins[0].Joints) - 1
	top := math.Inf(-1)
	hair := hairPrimitives(m)
	for _, p := range hair {
		for _, v := range p.Positions {
			top = max(top, float64(v[1]))
		}
	}
	var roots, ends int
	weight := func(p *gltf.Primitive, i int) float32 {
		var w float32
		for k, j := range p.Joints[i] {
			if int(j) == joint {
				w += p.Weights[i][k]
			}
		}
		return w
	}
	for _, p := range hair {
		for i, v := range p.Positions {
			depth := top - float64(v[1])
			if depth < 0.01 {
				roots++
				if weight(p, i) > 0.01 {
					t.Fatalf("a root %.3f m from the top hangs %.2f from the sway bone", depth, weight(p, i))
				}
			}
			if depth > hairLength {
				ends++
				if weight(p, i) < 0.5 {
					t.Fatalf("an end %.3f m from the top hangs only %.2f from the sway bone", depth, weight(p, i))
				}
			}
		}
	}
	if roots == 0 || ends == 0 {
		t.Fatalf("%d roots and %d ends: the hair has no length", roots, ends)
	}
	for _, p := range m.Meshes[0].Primitives {
		if hasFolder(m.Materials[p.Material].Name, "hair") {
			continue
		}
		for i := range p.Positions {
			if weight(&p, i) != 0 {
				t.Fatalf("a vertex of %s hangs from the sway bone", m.Materials[p.Material].Name)
			}
		}
	}
}

// Weights stay a whole: each vertex's add to one.
func TestHairWeightsStayWhole(t *testing.T) {
	m := womanWithHair(t, "hair/long01/long01.mhclo")
	for _, p := range hairPrimitives(m) {
		for i := range p.Weights {
			var sum float32
			for _, w := range p.Weights[i] {
				sum += w
			}
			if math.Abs(float64(sum)-1) > 1e-4 {
				t.Fatalf("vertex %d weighs %v", i, sum)
			}
		}
	}
}

// The hair's swing comes back to where it began, so a baked loop has no
// jump; it is small, moves with the person and is the same every time.
func TestHairSwayLoopsAndIsSmall(t *testing.T) {
	m := womanWithHair(t, "hair/long01/long01.mhclo")
	for _, name := range []string{"stand", "walk", "chop"} {
		c := swayChannel(t, m, name)
		first, last := c.Values[0], c.Values[len(c.Values)-1]
		if name != "chop" && first != last {
			t.Errorf("%s: the sway begins at %v and ends at %v", name, first, last)
		}
		moved := 0.0
		for _, v := range c.Values {
			angle := 2 * math.Acos(min(float64(v[3]), 1))
			if angle > 2*swayLimit*1.05 {
				t.Errorf("%s: the hair swings %.3f rad, past its limit", name, angle)
			}
			moved = max(moved, angle)
		}
		if moved < 0.005 {
			t.Errorf("%s: the hair does not move at all", name)
		}
	}
	// Walking moves the hair more than standing does.
	maxAngle := func(c gltf.Channel) float64 {
		a := 0.0
		for _, v := range c.Values {
			a = max(a, 2*math.Acos(min(float64(v[3]), 1)))
		}
		return a
	}
	if s, w := maxAngle(swayChannel(t, m, "stand")), maxAngle(swayChannel(t, m, "walk")); w <= s {
		t.Errorf("walking swings the hair %.3f, standing %.3f", w, s)
	}
	again := womanWithHair(t, "hair/long01/long01.mhclo")
	if !reflect.DeepEqual(swayChannel(t, m, "walk"), swayChannel(t, again, "walk")) {
		t.Error("the sway is not the same twice")
	}
}

// A model with no hair has no sway bone, and nothing is changed.
func TestABaldPersonHasNoSwayBone(t *testing.T) {
	k := realKit(t)
	b := baselines[1]
	m, err := k.Make(b.body, Look{Eyes: b.look.Eyes, Wear: []string{eyesProxy}})
	if err != nil {
		t.Fatal(err)
	}
	n := len(m.Nodes)
	if err := Animate(m, Walk); err != nil {
		t.Fatal(err)
	}
	if len(m.Nodes) != n {
		t.Error("a bone was added to a bald person")
	}
}
