package building

import (
	"image"
	"image/color"
	"math"

	"github.com/gabrielluizsf/sge/gltf"
)

// paint is a material as a part is finished with it: tinted, its colour
// times the surface's, in sRGB, white, or none, leaving it as it is; or
// painted, coloured the tint alone, its texture's relief and roughness
// showing through.
type paint struct {
	of      Material
	tint    color.NRGBA
	painted bool
}

// builder is a building being built: its palette, and the model so far,
// one primitive for each paint used, the materials' images shared.
type builder struct {
	palette Palette
	model   *gltf.Model
	paints  []paint
	prims   []gltf.Primitive
	laid    []Surface // each paint's surface, as its texture is laid
}

// newBuilder is an empty building, to be made of palette's surfaces.
func newBuilder(palette Palette) *builder {
	return &builder{palette: palette, model: &gltf.Model{Roots: []int{0}}}
}

// primitive is the index of paint p's primitive, made, with its
// material, the first time p is used.
func (b *builder) primitive(p paint) int {
	for i, q := range b.paints {
		if q == p {
			return i
		}
	}
	s := b.palette.surface(p.of)
	c := s.Color
	if p.tint.A != 0 {
		c = color.NRGBA{uint8(int(c.R) * int(p.tint.R) / 255), uint8(int(c.G) * int(p.tint.G) / 255), uint8(int(c.B) * int(p.tint.B) / 255), 255}
	}
	albedo := s.Albedo
	if p.painted {
		c, albedo = p.tint, nil
	}
	arm := b.image(s.ARM)
	b.model.Materials = append(b.model.Materials, gltf.Material{
		Name:      materialNames[p.of],
		BaseColor: [4]float64{linear(c.R), linear(c.G), linear(c.B), 1},
		Texture:   b.image(albedo), MetalRough: arm, Occlusion: arm, OcclusionStrength: 1,
		Normal: b.image(s.Normal), NormalScale: 1,
		Roughness: s.Roughness, Metallic: s.Metallic, AlphaMode: "OPAQUE", AlphaCutoff: 0.5,
	})
	b.paints = append(b.paints, p)
	b.prims = append(b.prims, gltf.Primitive{Material: len(b.model.Materials) - 1})
	b.laid = append(b.laid, s)
	return len(b.prims) - 1
}

// image is img's index of the model's images, added the first time it is
// used; -1 for none.
func (b *builder) image(img image.Image) int {
	if img == nil {
		return -1
	}
	for i, have := range b.model.Images {
		if have == img {
			return i
		}
	}
	b.model.Images = append(b.model.Images, img)
	return len(b.model.Images) - 1
}

// done is the building, a model of one node.
func (b *builder) done() *gltf.Model {
	b.model.Meshes = []gltf.Mesh{{Name: "building", Primitives: b.prims}}
	b.model.Nodes = []gltf.Node{{Name: "building", Mesh: 0, Skin: -1, Rotation: [4]float64{0, 0, 0, 1}, Scale: [3]float64{1, 1, 1}}}
	return b.model
}

// linear is an 8-bit sRGB value in linear light, as glTF's factors are.
func linear(v uint8) float64 {
	f := float64(v) / 255
	if f <= 0.04045 {
		return f / 12.92
	}
	return math.Pow((f+0.055)/1.055, 2.4)
}

// materialNames name each material in the model.
var materialNames = map[Material]string{
	Plaster: "plaster", Brick: "brick", Stone: "stone", Planks: "planks", Timber: "timber",
	Tiles: "tiles", Thatch: "thatch", Slates: "slates", Iron: "iron", Glass: "glass",
	Water: "water", Soot: "soot", Canvas: "canvas", Flat: "paint",
}
