package human

import (
	"fmt"

	"github.com/gabrielluizsf/sge/gltf"
)

// Motion is a way of moving: its name, how long it takes, in seconds, and
// the poses it passes through, evenly spaced. A loop's last runs on into
// its first; a Once motion's is where it ends, and the first where it
// began.
type Motion struct {
	Name   string
	Length float64
	Poses  []Pose
	Once   bool
}

// Animate adds motions to m, a person Make made, as its animations, for
// bake to pose by name and moment.
func Animate(m *gltf.Model, motions ...Motion) error {
	nodes := map[string]int{}
	for i, n := range m.Nodes {
		nodes[n.Name] = i
	}
	a, err := restArms(m, nodes)
	if err != nil {
		return err
	}
	root, ok := nodes["root"]
	if !ok {
		return fmt.Errorf("human: a model with no root bone")
	}
	for _, mo := range motions {
		if len(mo.Poses) == 0 || mo.Length <= 0 {
			return fmt.Errorf("human: motion %s has no poses", mo.Name)
		}
		an := gltf.Animation{Name: mo.Name}
		chans := map[string]*gltf.Channel{}
		lift := gltf.Channel{Node: root, Path: "translation", Interpolation: "LINEAR"}
		keys, spans := len(mo.Poses)+1, len(mo.Poses)
		if mo.Once {
			keys, spans = len(mo.Poses), len(mo.Poses)-1
		}
		for k := range keys {
			p := mo.Poses[k%len(mo.Poses)]
			at := float32(mo.Length * float64(k) / float64(max(spans, 1)))
			for name, q := range p.turns(a) {
				node, ok := nodes[name]
				if !ok {
					continue
				}
				ch := chans[name]
				if ch == nil {
					ch = &gltf.Channel{Node: node, Path: "rotation", Interpolation: "LINEAR"}
					chans[name] = ch
				}
				if n := len(ch.Values); n > 0 && dot4(ch.Values[n-1], q) < 0 {
					q = quat{-q[0], -q[1], -q[2], -q[3]} // the same turn, the near way round
				}
				ch.Times = append(ch.Times, at)
				ch.Values = append(ch.Values, [4]float32{float32(q[0]), float32(q[1]), float32(q[2]), float32(q[3])})
			}
			t := m.Nodes[root].Translation
			lift.Times = append(lift.Times, at)
			lift.Values = append(lift.Values, [4]float32{float32(t[0] + p.Aside), float32(t[1] + p.Bob), float32(t[2] + p.Forward)})
		}
		an.Channels = append(an.Channels, lift)
		for _, name := range sortedKeys(chans) {
			an.Channels = append(an.Channels, *chans[name])
		}
		swayHair(m, &an, keys, mo.Once)
		m.Animations = append(m.Animations, an)
	}
	return nil
}

// restArms are how the arms of m lie at rest: the upper arm from the
// shoulder's head to the elbow's, the forearm from there to the wrist.
func restArms(m *gltf.Model, nodes map[string]int) (arms, error) {
	var a arms
	for s, side := range sides {
		var at [4][3]float64
		for i, name := range []string{"upperarm02", "lowerarm01", "lowerarm02", "wrist"} {
			n, ok := nodes[name+side.suffix]
			if !ok {
				return a, fmt.Errorf("human: a model with no %s", name+side.suffix)
			}
			at[i] = m.Nodes[n].Translation
		}
		along, fore := add(at[0], at[1]), add(at[2], at[3])
		a.hang[s] = between(along, [3]float64{side.out * 0.12, -1, 0})
		a.straighten[s] = between(fore, along)
		a.elbow[s] = normalize(cross(normalize(along), zAxis))
		a.plumb[s] = plumbLeg(m, nodes, side.suffix)
	}
	return a, nil
}

// plumbLeg is the turn that hangs a leg of m, whose bones' names end in
// suffix, straight down from the hip to the ankle: MakeHuman's rest
// stands with the legs apart. A model with no legs has none to turn.
func plumbLeg(m *gltf.Model, nodes map[string]int, suffix string) quat {
	var leg [3]float64
	for _, name := range []string{"upperleg02", "lowerleg01", "lowerleg02", "foot"} {
		n, ok := nodes[name+suffix]
		if !ok {
			return still
		}
		leg = add(leg, m.Nodes[n].Translation)
	}
	return between(leg, [3]float64{0, -1, 0})
}
