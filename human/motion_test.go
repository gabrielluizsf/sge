package human

import (
	"image/color"
	"math"
	"testing"
	"testing/fstest"

	"github.com/gabrielluizsf/sge/gltf"
)

func colorOf(v uint8) color.NRGBA { return color.NRGBA{v, v, v, 255} }

// Every motion loops: a channel's last key is its first, a turn of length
// one, each key on the near side of the last, as many keys a channel as
// the motion has poses and one.
func TestMotionsLoop(t *testing.T) {
	k := loaded(t, tinyKit(t))
	m, err := k.Make(Average, Look{Skin: "skins/plain/plain.mhmat"})
	if err != nil {
		t.Fatal(err)
	}
	if err := Animate(m, Motions...); err == nil {
		t.Fatal("animated a body with no arms")
	}
	m = armed(t)
	if err := Animate(m, Motions...); err != nil {
		t.Fatal(err)
	}
	if len(m.Animations) != len(Motions) {
		t.Fatalf("%d animations", len(m.Animations))
	}
	for i, a := range m.Animations {
		for _, ch := range a.Channels {
			if len(ch.Values) != len(Motions[i].Poses)+1 || ch.Times[len(ch.Times)-1] != float32(Motions[i].Length) {
				t.Errorf("%s: %d keys to %v", a.Name, len(ch.Values), ch.Times[len(ch.Times)-1])
			}
			first, last := ch.Values[0], ch.Values[len(ch.Values)-1]
			for c := range 4 {
				if math.Abs(float64(first[c]-last[c])) > 1e-5 && math.Abs(float64(first[c]+last[c])) > 1e-5 {
					t.Errorf("%s: node %d ends at %v, not %v", a.Name, ch.Node, last, first)
				}
			}
			for j := 1; ch.Path == "rotation" && j < len(ch.Values); j++ {
				q := quat{float64(ch.Values[j][0]), float64(ch.Values[j][1]), float64(ch.Values[j][2]), float64(ch.Values[j][3])}
				if dot4(ch.Values[j-1], q) < 0 {
					t.Errorf("%s: node %d turns the long way round at key %d", a.Name, ch.Node, j)
				}
				if l := dot4(ch.Values[j], q); math.Abs(l-1) > 1e-4 {
					t.Errorf("%s: a turn of length %v", a.Name, l)
				}
			}
		}
	}
}

// Standing, an arm let down from MakeHuman's rest hangs by the body, its
// forearm straight on, not reaching forward.
func TestArmsHangAtRest(t *testing.T) {
	m := armed(t)
	nodes := map[string]int{}
	for i, n := range m.Nodes {
		nodes[n.Name] = i
	}
	a, err := restArms(m, nodes)
	if err != nil {
		t.Fatal(err)
	}
	turns := Pose{}.turns(a)
	upper := turns["upperarm01.L"]
	along := add(m.Nodes[nodes["upperarm02.L"]].Translation, m.Nodes[nodes["lowerarm01.L"]].Translation)
	fore := add(m.Nodes[nodes["lowerarm02.L"]].Translation, m.Nodes[nodes["wrist.L"]].Translation)
	hand := upper.turn(turns["lowerarm01.L"].turn(fore))
	if d := normalize(upper.turn(along)); d[1] > -0.95 || math.Abs(d[2]) > 0.01 {
		t.Errorf("the upper arm hangs along %v", d)
	}
	if d := normalize(hand); d[1] > -0.95 || math.Abs(d[2]) > 0.01 {
		t.Errorf("the forearm hangs along %v", d)
	}
}

// armed is a model with MakeHuman's arm bones, as at its rest: out and
// down, the forearm bent forward.
func armed(t *testing.T) *gltf.Model {
	t.Helper()
	fsys := tinyKit(t)
	bones := `{"bones": {"root": {"head": "a", "tail": "b", "parent": null}`
	joints := `"a": [0], "b": [1]`
	for _, s := range []string{"L", "R"} {
		bones += `, "upperarm02.` + s + `": {"head": "a", "tail": "e` + s + `", "parent": "root"}` +
			`, "lowerarm01.` + s + `": {"head": "e` + s + `", "tail": "f` + s + `", "parent": "upperarm02.` + s + `"}` +
			`, "lowerarm02.` + s + `": {"head": "f` + s + `", "tail": "w` + s + `", "parent": "lowerarm01.` + s + `"}` +
			`, "wrist.` + s + `": {"head": "w` + s + `", "tail": "w` + s + `", "parent": "lowerarm02.` + s + `"}`
		joints += `, "e` + s + `": [` + map[string]string{"L": "2", "R": "3"}[s] + `], "f` + s + `": [` + map[string]string{"L": "4", "R": "5"}[s] + `], "w` + s + `": [6]`
	}
	fsys[baseSkeleton] = &fstest.MapFile{Data: []byte(bones + `}, "joints": {` + joints + `}}`)}
	fsys[baseWeights] = &fstest.MapFile{Data: []byte(`{"weights": {"root": [[0, 1]]}}`)}
	// Elbows out and down, forearms on forward.
	fsys[baseMesh] = &fstest.MapFile{Data: []byte("v 0 14 0\nv 0 15 0\nv 2 12 0\nv -2 12 0\nv 4 10 0\nv -4 10 0\nv 5 9 2\ng body\nf 1 2 3\n")}
	m, err := loaded(t, fsys).Make(Average, Look{Skin: "skins/plain/plain.mhmat"})
	if err != nil {
		t.Fatal(err)
	}
	return m
}
