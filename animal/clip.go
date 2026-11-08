package animal

import "github.com/gabrielluizsf/sge/gltf"

// Pose is how the bones are held at a moment: world turns of bones, which
// their bones below follow unless turned too, and a shift of the root.
type Pose struct {
	turns map[int]Quat
	Shift Vec
}

// NewPose is the pose of a rig at rest.
func NewPose() *Pose { return &Pose{turns: map[int]Quat{}} }

// Turn turns bone b, and what hangs from it, by q in the scene's axes,
// about the bone's own head.
func (p *Pose) Turn(b int, q Quat) { p.turns[b] = q }

// Clip writes the animation called name onto the rig: the pose of each
// moment of times, held between as the bone turns between them.
func (r *Rig) Clip(name string, times []float64, pose func(i int) *Pose) {
	joints := r.M.Skins[0].Joints
	rots := make([][][4]float32, len(r.M.Nodes))
	var moves [][4]float32
	root := r.Root()
	for i := range times {
		p := pose(i)
		world := make([]Quat, len(r.M.Nodes)) // each node's turn as posed
		turned := make([]Quat, len(r.M.Nodes))
		var walk func(n int, inherited Quat)
		walk = func(n int, inherited Quat) {
			d := inherited
			if q, ok := p.turns[n]; ok && r.joint[n] {
				d = q
			}
			if !r.joint[n] {
				d = identity
			}
			turned[n] = d
			world[n] = mul(d, r.rest[n].rot)
			for _, c := range r.M.Nodes[n].Children {
				walk(c, d)
			}
		}
		for _, rt := range r.M.Roots {
			walk(rt, identity)
		}
		for _, j := range joints {
			up := identity
			if par := r.parent[j]; par >= 0 {
				up = world[par]
			}
			q := mul(conj(up), world[j])
			if k := len(rots[j]); k > 0 {
				q = sameSide(q, Quat(floats(rots[j][k-1])))
			}
			rots[j] = append(rots[j], [4]float32{float32(q[0]), float32(q[1]), float32(q[2]), float32(q[3])})
		}
		par := r.parent[root]
		up, size := identity, 1.0
		if par >= 0 {
			up, size = r.rest[par].rot, r.rest[par].size
		}
		t := turn(conj(up), scale(p.Shift, 1/size))
		base := Vec(r.M.Nodes[root].Translation)
		moves = append(moves, [4]float32{float32(base[0] + t[0]), float32(base[1] + t[1]), float32(base[2] + t[2])})
	}
	at := make([]float32, len(times))
	for i, t := range times {
		at[i] = float32(t)
	}
	a := gltf.Animation{Name: name}
	for _, j := range joints {
		a.Channels = append(a.Channels, gltf.Channel{Node: j, Path: "rotation", Interpolation: "LINEAR", Times: at, Values: rots[j]})
	}
	a.Channels = append(a.Channels, gltf.Channel{Node: root, Path: "translation", Interpolation: "LINEAR", Times: at, Values: moves})
	r.M.Animations = append(r.M.Animations, a)
}

func floats(v [4]float32) [4]float64 {
	return [4]float64{float64(v[0]), float64(v[1]), float64(v[2]), float64(v[3])}
}
