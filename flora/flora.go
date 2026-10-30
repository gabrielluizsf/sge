// Package flora grows plants in code, as 3D models to bake: crops from
// seed to ripe, each plant a little unlike the next. Leaves are drawn
// too, their outline, veins and greens painted into textures, so nothing
// in a plant is anyone's but the game's.
package flora

import "github.com/gabrielluizsf/sge/gltf"

// plant is a model being grown: its primitives, a material each.
type plant struct {
	model *gltf.Model
	prims []gltf.Primitive
}

// newPlant is an empty plant with the materials it may use.
func newPlant(mats []gltf.Material, images []imageOf) *plant {
	m := &gltf.Model{Roots: []int{0}, Materials: mats}
	for _, im := range images {
		m.Images = append(m.Images, im)
	}
	return &plant{model: m}
}

// done is the plant as a model of one node.
func (p *plant) done() *gltf.Model {
	p.model.Meshes = []gltf.Mesh{{Name: "plant", Primitives: p.prims}}
	p.model.Nodes = []gltf.Node{{Name: "plant", Mesh: 0, Skin: -1, Rotation: [4]float64{0, 0, 0, 1}, Scale: [3]float64{1, 1, 1}}}
	return p.model
}
