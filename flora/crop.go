package flora

import (
	"image/color"
	"math"

	"github.com/gabrielluizsf/sge/gltf"
	"github.com/gabrielluizsf/sge/random"
)

// Crop is a kind of crop: its name, how many stages it grows through,
// from a sprout at 0 to ripe at the last, and how it grows.
type Crop struct {
	Name   string
	Stages int
	leaf   Leaf
	grow   func(g *grower, stage int)
}

// Grow is the crop at stage, a plant of its own as seed makes it: no two
// seeds grow quite alike.
func (c Crop) Grow(stage int, seed uint32) *gltf.Model {
	stage = max(0, min(stage, c.Stages-1))
	g := &grower{luck: random.New(seed*2654435761 + uint32(stage)), grown: float64(stage+1) / float64(c.Stages)}
	g.p = newPlant([]gltf.Material{
		{Name: "leaf", BaseColor: [4]float64{1, 1, 1, 1}, Texture: 0, MetalRough: -1, Normal: -1, Occlusion: -1,
			Roughness: 0.6, AlphaMode: "MASK", AlphaCutoff: 0.5, DoubleSided: true},
	}, []imageOf{c.leaf.texture(128, seed)})
	g.stemMat = g.solid(color.NRGBA{88, 120, 52, 255}, 0.7)
	c.grow(g, stage)
	return g.p.done()
}

// grower grows a plant: its random numbers, how far grown it is, from 0
// to 1, and the materials made so far.
type grower struct {
	p       *plant
	luck    *random.Source
	grown   float64
	stemMat int
}

// solid is a new material of one colour, given in sRGB, and roughness.
func (g *grower) solid(c color.NRGBA, rough float64) int {
	lin := func(v uint8) float64 {
		f := float64(v) / 255
		if f <= 0.04045 {
			return f / 12.92
		}
		return math.Pow((f+0.055)/1.055, 2.4)
	}
	g.p.model.Materials = append(g.p.model.Materials, gltf.Material{BaseColor: [4]float64{lin(c.R), lin(c.G), lin(c.B), 1},
		Texture: -1, MetalRough: -1, Normal: -1, Occlusion: -1, Roughness: rough, AlphaMode: "OPAQUE", AlphaCutoff: 0.5, DoubleSided: true})
	return len(g.p.model.Materials) - 1
}

func (g *grower) add(p gltf.Primitive) { g.p.prims = append(g.p.prims, p) }

// between is a random number from lo to hi.
func (g *grower) between(lo, hi float64) float64 { return lo + (hi-lo)*g.luck.Float() }

// whorl is leaves round a point, at: how many, how long, how wide for
// their length, and the pitches they rise between, outer to inner.
type whorl struct {
	n           int
	at          vec
	size, width float64
	lo, hi      float64
}

// rosette grows a whorl's leaves, each turned round by the golden angle
// as plants place them.
func (g *grower) rosette(w whorl) {
	n, at, size, width, lo, hi := w.n, w.at, w.size, w.width, w.lo, w.hi
	turn := g.between(0, 2*math.Pi)
	for i := range n {
		inner := float64(i) / math.Max(float64(n-1), 1)
		l := size * g.between(0.8, 1.15) * (1 - 0.3*inner)
		g.add(leaf(LeafShape{Base: at, Yaw: turn + float64(i)*2.39996, Pitch: lo + (hi-lo)*inner + g.between(-0.1, 0.1),
			Length: l, Width: l * width, Droop: l * g.between(0.1, 0.3) * (1 - inner)}, 0))
	}
}
