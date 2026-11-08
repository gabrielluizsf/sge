package tree

import (
	"image"
	"image/color"
	"math"

	"github.com/gabrielluizsf/sge/gltf"
	"github.com/gabrielluizsf/sge/random"
)

// grower grows a tree: its random numbers, the metres it grows to (its
// height, or its width where the species is wide), the barks it may
// wear, and the model so far, one primitive to each material.
type grower struct {
	model *gltf.Model
	prims []gltf.Primitive
	luck  *random.Source
	size  float64
	stage Size
	barks Barks
}

// newGrower is a grower of nothing yet, of the tree spec says.
func newGrower(sp Spec, seed uint32) *grower {
	return &grower{model: &gltf.Model{Roots: []int{0}}, luck: random.New(seed), stage: sp.Size, barks: sp.Barks}
}

// between is a random number from lo to hi.
func (g *grower) between(lo, hi float64) float64 { return lo + (hi-lo)*g.luck.Float() }

// count is a random whole number from lo to hi, both included.
func (g *grower) count(lo, hi int) int { return lo + g.luck.Intn(hi-lo+1) }

// anywhere is a random unit direction.
func (g *grower) anywhere() vec {
	for {
		v := vec{g.between(-1, 1), g.between(-1, 1), g.between(-1, 1)}
		if l := length(v); l > 0.1 && l <= 1 {
			return scale(v, 1/l)
		}
	}
}

// material adds m, and the primitive that will carry it; its index.
func (g *grower) material(m gltf.Material) int {
	g.model.Materials = append(g.model.Materials, m)
	g.prims = append(g.prims, gltf.Primitive{Material: len(g.model.Materials) - 1})
	return len(g.prims) - 1
}

// prim is the primitive of material mat.
func (g *grower) prim(mat int) *gltf.Primitive { return &g.prims[mat] }

// image adds img to the model; its index, or -1 for none.
func (g *grower) image(img image.Image) int {
	if img == nil {
		return -1
	}
	g.model.Images = append(g.model.Images, img)
	return len(g.model.Images) - 1
}

// cutout is a material of texture img with its clear parts cut out, as
// leaves, needles and petals are, named name.
func (g *grower) cutout(name string, img image.Image) int {
	return g.material(gltf.Material{Name: name, BaseColor: [4]float64{1, 1, 1, 1}, Texture: g.image(img), MetalRough: -1, Normal: -1, Occlusion: -1,
		Roughness: 0.7, AlphaMode: "MASK", AlphaCutoff: 0.5, DoubleSided: true})
}

// solid is a material of one colour, in sRGB, and roughness rough.
func (g *grower) solid(c color.NRGBA, rough float64) int {
	return g.material(gltf.Material{Name: "paint", BaseColor: [4]float64{linear(c.R), linear(c.G), linear(c.B), 1}, Texture: -1, MetalRough: -1,
		Normal: -1, Occlusion: -1, Roughness: rough, AlphaMode: "OPAQUE", AlphaCutoff: 0.5, DoubleSided: true})
}

// textured is an opaque material of texture img, tiling, rough.
func (g *grower) textured(img image.Image, rough float64) int {
	return g.material(gltf.Material{Name: "skin", BaseColor: [4]float64{1, 1, 1, 1}, Texture: g.image(img), MetalRough: -1, Normal: -1, Occlusion: -1,
		Roughness: rough, AlphaMode: "OPAQUE", AlphaCutoff: 0.5, DoubleSided: true})
}

// bark is the material of bark b, its textures as Poly Haven's are, or
// plain brown where it has none.
func (g *grower) bark(b Bark) int {
	c := b.Color
	if c.A == 0 {
		c = color.NRGBA{110, 84, 62, 255}
	}
	arm := g.image(b.ARM)
	return g.material(gltf.Material{Name: "bark", BaseColor: [4]float64{linear(c.R), linear(c.G), linear(c.B), 1}, Texture: g.image(b.Albedo),
		MetalRough: arm, Occlusion: arm, OcclusionStrength: 1, Normal: g.image(b.Normal), NormalScale: 1,
		Roughness: 1, AlphaMode: "OPAQUE", AlphaCutoff: 0.5, DoubleSided: true})
}

// done is the tree as a model of one node.
func (g *grower) done() *gltf.Model {
	g.model.Meshes = []gltf.Mesh{{Name: "tree", Primitives: g.prims}}
	g.model.Nodes = []gltf.Node{{Name: "tree", Mesh: 0, Skin: -1, Rotation: [4]float64{0, 0, 0, 1}, Scale: [3]float64{1, 1, 1}}}
	return g.model
}

// linear is an 8-bit sRGB value in linear light, as glTF's factors are.
func linear(v uint8) float64 {
	f := float64(v) / 255
	if f <= 0.04045 {
		return f / 12.92
	}
	return math.Pow((f+0.055)/1.055, 2.4)
}
