// Package shapes makes the small glTF models render3d's tests draw: a
// textured cube, a skinned arm and a rig of many bones.
package shapes

import (
	"image"
	"image/color"
	"math"

	"github.com/gabrielluizsf/sge/gltf"
	"github.com/gabrielluizsf/sge/pose"
)

var (
	// Still is no rotation, One no scale.
	Still = [4]float64{0, 0, 0, 1}
	One   = [3]float64{1, 1, 1}
)

// Cuboid is a box round Center, Half its size each way, following Joint.
type Cuboid struct {
	Center, Half [3]float32
	Joint        uint16
}

// Box is b's six faces, in glTF's axes, added to p: four corners a face,
// counterclockwise from outside, its texture running right and down
// across it.
func Box(p *gltf.Primitive, b Cuboid) {
	center, half := b.Center, b.Half
	type face struct{ n, u [3]float32 }
	faces := []face{
		{[3]float32{0, 0, 1}, [3]float32{1, 0, 0}}, {[3]float32{0, 0, -1}, [3]float32{-1, 0, 0}},
		{[3]float32{1, 0, 0}, [3]float32{0, 0, -1}}, {[3]float32{-1, 0, 0}, [3]float32{0, 0, 1}},
		{[3]float32{0, 1, 0}, [3]float32{1, 0, 0}}, {[3]float32{0, -1, 0}, [3]float32{1, 0, 0}},
	}
	for _, f := range faces {
		v := [3]float32{f.n[1]*f.u[2] - f.n[2]*f.u[1], f.n[2]*f.u[0] - f.n[0]*f.u[2], f.n[0]*f.u[1] - f.n[1]*f.u[0]}
		base := uint32(len(p.Positions))
		for _, c := range [4][2]float32{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}} {
			var pos [3]float32
			for k := range 3 {
				pos[k] = center[k] + half[k]*(f.n[k]+c[0]*f.u[k]+c[1]*v[k])
			}
			p.Positions = append(p.Positions, pos)
			p.Normals = append(p.Normals, f.n)
			p.UVs = append(p.UVs, [2]float32{(c[0] + 1) / 2, (1 - c[1]) / 2})
			p.Joints = append(p.Joints, [4]uint16{b.Joint})
			p.Weights = append(p.Weights, [4]float32{1})
		}
		p.Indices = append(p.Indices, base, base+1, base+2, base, base+2, base+3)
	}
}

// Cube is a cube one across round its origin, its texture red on its
// left half and blue on its right.
func Cube() *gltf.Model {
	var p gltf.Primitive
	Box(&p, Cuboid{Half: [3]float32{0.5, 0.5, 0.5}})
	p.Joints, p.Weights = nil, nil
	img := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	for y := range 2 {
		img.SetNRGBA(0, y, color.NRGBA{255, 0, 0, 255})
		img.SetNRGBA(1, y, color.NRGBA{0, 0, 255, 255})
	}
	return &gltf.Model{
		Roots:     []int{0},
		Nodes:     []gltf.Node{{Mesh: 0, Skin: -1, Rotation: Still, Scale: One}},
		Meshes:    []gltf.Mesh{{Primitives: []gltf.Primitive{p}}},
		Materials: []gltf.Material{{BaseColor: [4]float64{1, 1, 1, 1}, Texture: 0, MetalRough: -1, Normal: -1, Occlusion: -1, Roughness: 1, AlphaMode: "OPAQUE"}},
		Images:    []image.Image{img},
	}
}

// Plain is Cube untextured, white.
func Plain() *gltf.Model {
	m := Cube()
	m.Materials[0].Texture, m.Images = -1, nil
	return m
}

// Arm is a bar two high standing on its origin, in glTF's axes, skinned
// to two bones: its lower half to a shoulder at the origin, its upper to
// an elbow one up, which "bend" turns a quarter round Z, toward -X, in
// its first second, and holds there for another.
func Arm() *gltf.Model {
	var p gltf.Primitive
	Box(&p, Cuboid{[3]float32{0, 0.5, 0}, [3]float32{0.1, 0.5, 0.1}, 0})
	Box(&p, Cuboid{[3]float32{0, 1.5, 0}, [3]float32{0.1, 0.5, 0.1}, 1})
	p.Material = -1
	q := float32(math.Sqrt(0.5))
	return &gltf.Model{
		Roots: []int{0, 1},
		Nodes: []gltf.Node{
			{Name: "body", Mesh: 0, Skin: 0, Rotation: Still, Scale: One},
			{Name: "shoulder", Children: []int{2}, Mesh: -1, Skin: -1, Rotation: Still, Scale: One},
			{Name: "elbow", Mesh: -1, Skin: -1, Translation: [3]float64{0, 1, 0}, Rotation: Still, Scale: One},
		},
		Meshes: []gltf.Mesh{{Primitives: []gltf.Primitive{p}}},
		Skins:  []gltf.Skin{{Joints: []int{1, 2}, InverseBind: [][16]float64{pose.Identity, {0: 1, 5: 1, 10: 1, 13: -1, 15: 1}}}},
		Animations: []gltf.Animation{{Name: "bend", Channels: []gltf.Channel{{
			Node: 2, Path: "rotation", Interpolation: "LINEAR",
			Times: []float32{0, 1, 2}, Values: [][4]float32{{0, 0, 0, 1}, {0, 0, q, q}, {0, 0, q, q}},
		}}}},
	}
}

// Rig is a character of bones bones, as MakeHuman's default skeleton has
// 163, in chains of eight from a root, each turned and moved by a "walk"
// of 30 keys a second, and a box skinned to every bone.
func Rig(bones int) *gltf.Model {
	m := &gltf.Model{Roots: []int{0, 1}, Nodes: []gltf.Node{{Name: "body", Mesh: 0, Skin: 0, Rotation: Still, Scale: One}}}
	var p gltf.Primitive
	var skin gltf.Skin
	walk := gltf.Animation{Name: "walk"}
	for b := range bones {
		node := len(m.Nodes)
		if b > 0 && b%8 != 1 {
			m.Nodes[node-1].Children = []int{node}
		} else if b > 0 {
			m.Nodes[1].Children = append(m.Nodes[1].Children, node)
		}
		m.Nodes = append(m.Nodes, gltf.Node{Mesh: -1, Skin: -1, Translation: [3]float64{0, 0.1, 0}, Rotation: Still, Scale: One})
		skin.Joints = append(skin.Joints, node)
		skin.InverseBind = append(skin.InverseBind, pose.Identity)
		Box(&p, Cuboid{[3]float32{0, float32(b) * 0.01, 0}, [3]float32{0.05, 0.05, 0.05}, uint16(b)})
		turn := gltf.Channel{Node: node, Path: "rotation", Interpolation: "LINEAR"}
		move := gltf.Channel{Node: node, Path: "translation", Interpolation: "LINEAR"}
		for k := range 30 {
			a := float64(k) / 30 * math.Pi
			turn.Times, move.Times = append(turn.Times, float32(k)/30), append(move.Times, float32(k)/30)
			turn.Values = append(turn.Values, [4]float32{float32(math.Sin(a / 4)), 0, 0, float32(math.Cos(a / 4))})
			move.Values = append(move.Values, [4]float32{0, 0.1, float32(a) * 0.01})
		}
		walk.Channels = append(walk.Channels, turn, move)
	}
	m.Meshes = []gltf.Mesh{{Primitives: []gltf.Primitive{p}}}
	m.Skins = []gltf.Skin{skin}
	m.Animations = []gltf.Animation{walk}
	return m
}

// Near is whether c is within 8 of each of rgb.
func Near(c color.NRGBA, rgb [3]int) bool {
	d := func(a uint8, b int) bool { return int(a)-b <= 8 && b-int(a) <= 8 }
	return d(c.R, rgb[0]) && d(c.G, rgb[1]) && d(c.B, rgb[2])
}
