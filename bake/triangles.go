package bake

import (
	"math"

	"github.com/gabrielluizsf/sge/gltf"
)

// vertex is a corner of a triangle in the game's world: where it is, the
// way its surface faces, and its place on its texture.
type vertex struct {
	p, n [3]float64
	t    [4]float64 // tangent, and the bitangent's sign
	uv   [2]float64
}

// triangle is three vertices of a material.
type triangle struct {
	v   [3]vertex
	mat *material
}

// material is a glTF material with its textures read for sampling: its
// colour, metal-rough, normal and occlusion maps, nil where it has none.
type material struct {
	gltf.Material
	tex, mr, normal, ao *texture
}

// toGame is the model's space turned into the game's: glTF's Y up becomes
// the game's Z, its Z toward the viewer the game's -Y, a rotation, so
// nothing is mirrored; then turned by yaw and scaled by unit. A model
// facing glTF's +Z, as models do, then faces the viewer, down the screen.
func toGame(yaw, unit float64) mat {
	// Facing the viewer is facing +X+Y; glTF's +Z lands on -Y, so a
	// further turn of 135 degrees brings it round.
	a := (yaw + 135) * math.Pi / 180
	c, s := math.Cos(a)*unit, math.Sin(a)*unit
	rot := mat{c, s, 0, 0, -s, c, 0, 0, 0, 0, unit, 0, 0, 0, 0, 1}
	axes := mat{1, 0, 0, 0, 0, 0, 1, 0, 0, -1, 0, 0, 0, 0, 0, 1}
	return rot.mul(axes)
}

// triangles are the model's triangles in the game's world, as shot: each
// mesh where its node puts it, or, skinned, where its joints bend it.
func triangles(m *gltf.Model, s Shot, unit float64) []triangle {
	w := worlds(m, s)
	game := toGame(s.Yaw, unit)
	mats := materials(m)
	var out []triangle
	shown := reachable(m, s)
	for ni, node := range m.Nodes {
		if node.Mesh < 0 || !shown[ni] {
			continue
		}
		joints := skinMatrices(m, node.Skin, w)
		for _, p := range m.Meshes[node.Mesh].Primitives {
			mt := &mats[0]
			if p.Material >= 0 {
				mt = &mats[p.Material+1]
			}
			verts := make([]vertex, len(p.Positions))
			for i := range verts {
				place := w[ni]
				if joints != nil && i < len(p.Joints) && i < len(p.Weights) {
					place = blend(joints, p.Joints[i], p.Weights[i])
				}
				full := game.mul(place)
				verts[i].p = full.point(f3(p.Positions[i]))
				if i < len(p.Normals) {
					verts[i].n = normalize(full.dir(f3(p.Normals[i])))
				}
				if i < len(p.UVs) {
					verts[i].uv = [2]float64{float64(p.UVs[i][0]), float64(p.UVs[i][1])}
				}
				if i < len(p.Tangents) {
					tg := p.Tangents[i]
					d := normalize(full.dir([3]float64{float64(tg[0]), float64(tg[1]), float64(tg[2])}))
					verts[i].t = [4]float64{d[0], d[1], d[2], float64(tg[3])}
				}
			}
			for t := 0; t+2 < len(p.Indices); t += 3 {
				tri := triangle{mat: mt}
				for k := range 3 {
					tri.v[k] = verts[p.Indices[t+k]]
				}
				if len(p.Normals) == 0 {
					n := normalize(cross(sub(tri.v[1].p, tri.v[0].p), sub(tri.v[2].p, tri.v[0].p)))
					tri.v[0].n, tri.v[1].n, tri.v[2].n = n, n, n
				}
				if len(p.Tangents) == 0 && mt.normal != nil {
					tangentOf(&tri)
				}
				out = append(out, tri)
			}
		}
	}
	return out
}

// skinMatrices are skin's joints' matrices, posed by w, or nil for none.
func skinMatrices(m *gltf.Model, skin int, w []mat) []mat {
	if skin < 0 {
		return nil
	}
	sk := m.Skins[skin]
	out := make([]mat, len(sk.Joints))
	for j, node := range sk.Joints {
		out[j] = w[node].mul(mat(sk.InverseBind[j]))
	}
	return out
}

// blend is the joints' matrices mixed by the weights; a joint past the
// skin's is left out.
func blend(joints []mat, js [4]uint16, ws [4]float32) mat {
	var m mat
	var total float64
	for k := range 4 {
		if int(js[k]) >= len(joints) || ws[k] == 0 {
			continue
		}
		total += float64(ws[k])
		for i := range m {
			m[i] += joints[js[k]][i] * float64(ws[k])
		}
	}
	if total == 0 {
		return identity
	}
	for i := range m {
		m[i] /= total
	}
	return m
}

// materials are the model's materials, after a default white one, with
// their textures read: colour as sRGB, the maps of data as they are.
func materials(m *gltf.Model) []material {
	cache := map[[2]int]*texture{}
	read := func(img int, colour bool) *texture {
		if img < 0 {
			return nil
		}
		k := [2]int{img, map[bool]int{true: 1}[colour]}
		if cache[k] == nil {
			cache[k] = textureOf(m.Images[img], colour)
		}
		return cache[k]
	}
	out := []material{{Material: gltf.Material{BaseColor: [4]float64{1, 1, 1, 1}, Texture: -1, MetalRough: -1, Normal: -1, Occlusion: -1,
		Roughness: 1, AlphaMode: "OPAQUE", AlphaCutoff: 0.5}}}
	for _, g := range m.Materials {
		out = append(out, material{Material: g, tex: read(g.Texture, true), mr: read(g.MetalRough, false),
			normal: read(g.Normal, false), ao: read(g.Occlusion, false)})
	}
	return out
}

// tangentOf gives a triangle the tangent its texture runs along, worked
// out from how its UVs change across it, for a normal map to bend by.
func tangentOf(t *triangle) {
	e1, e2 := sub(t.v[1].p, t.v[0].p), sub(t.v[2].p, t.v[0].p)
	du1, dv1 := t.v[1].uv[0]-t.v[0].uv[0], t.v[1].uv[1]-t.v[0].uv[1]
	du2, dv2 := t.v[2].uv[0]-t.v[0].uv[0], t.v[2].uv[1]-t.v[0].uv[1]
	det := du1*dv2 - du2*dv1
	if det == 0 {
		return
	}
	r := 1 / det
	tg := normalize([3]float64{(e1[0]*dv2 - e2[0]*dv1) * r, (e1[1]*dv2 - e2[1]*dv1) * r, (e1[2]*dv2 - e2[2]*dv1) * r})
	bt := [3]float64{(e2[0]*du1 - e1[0]*du2) * r, (e2[1]*du1 - e1[1]*du2) * r, (e2[2]*du1 - e1[2]*du2) * r}
	for k := range t.v {
		w := 1.0
		if dot(cross(t.v[k].n, tg), bt) < 0 {
			w = -1
		}
		t.v[k].t = [4]float64{tg[0], tg[1], tg[2], w}
	}
}

func f3(v [3]float32) [3]float64 { return [3]float64{float64(v[0]), float64(v[1]), float64(v[2])} }

func sub(a, b [3]float64) [3]float64 { return [3]float64{a[0] - b[0], a[1] - b[1], a[2] - b[2]} }

// reachable is which nodes the shot shows: every one, or those under its
// one root.
func reachable(m *gltf.Model, s Shot) []bool {
	shown := make([]bool, len(m.Nodes))
	var walk func(i int)
	walk = func(i int) {
		shown[i] = true
		for _, c := range m.Nodes[i].Children {
			walk(c)
		}
	}
	for i, r := range m.Roots {
		if !s.Only || i == s.Root {
			walk(r)
		}
	}
	return shown
}
