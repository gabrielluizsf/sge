package bake

import (
	"image/color"
	"testing"

	"github.com/gabrielluizsf/sge/gltf"
	"github.com/gabrielluizsf/sge/iso"
)

// box is a model of a box, size across and high, standing on its origin,
// in glTF's axes (Y up), each face a quad with its own normal, in colour.
func box(size, high float32, rgb [3]float64) *gltf.Model {
	h := size / 2
	type face struct {
		n       [3]float32
		corners [4][3]float32
	}
	faces := []face{
		{[3]float32{0, 1, 0}, [4][3]float32{{-h, high, -h}, {-h, high, h}, {h, high, h}, {h, high, -h}}},
		{[3]float32{1, 0, 0}, [4][3]float32{{h, 0, -h}, {h, high, -h}, {h, high, h}, {h, 0, h}}},
		{[3]float32{-1, 0, 0}, [4][3]float32{{-h, 0, h}, {-h, high, h}, {-h, high, -h}, {-h, 0, -h}}},
		{[3]float32{0, 0, 1}, [4][3]float32{{h, 0, h}, {h, high, h}, {-h, high, h}, {-h, 0, h}}},
		{[3]float32{0, 0, -1}, [4][3]float32{{-h, 0, -h}, {-h, high, -h}, {h, high, -h}, {h, 0, -h}}},
	}
	var p gltf.Primitive
	for _, f := range faces {
		base := uint32(len(p.Positions))
		for _, c := range f.corners {
			p.Positions = append(p.Positions, c)
			p.Normals = append(p.Normals, f.n)
		}
		p.Indices = append(p.Indices, base, base+1, base+2, base, base+2, base+3)
	}
	p.Material = 0
	return &gltf.Model{
		Roots:     []int{0},
		Nodes:     []gltf.Node{{Mesh: 0, Skin: -1, Rotation: [4]float64{0, 0, 0, 1}, Scale: [3]float64{1, 1, 1}}},
		Meshes:    []gltf.Mesh{{Primitives: []gltf.Primitive{p}}},
		Materials: []gltf.Material{{BaseColor: [4]float64{rgb[0], rgb[1], rgb[2], 1}, Texture: -1, MetalRough: -1, Normal: -1, Occlusion: -1, Roughness: 0.9, AlphaMode: "OPAQUE", DoubleSided: true}},
	}
}

func options() Options { return DefaultOptions(iso.New(32, 64, 32)) }

func luma(c color.NRGBA) int { return int(c.R) + int(c.G) + int(c.B) }

// A box in the sun, lit from the upper left of the screen: its top
// brightest, its left side lit, its right side in shade; its shadow on the ground beyond it, black and part
// see-through; air in the picture's corners; and the point it stands on
// at the pivot.
func TestABoxIsLitAndCastsAShadow(t *testing.T) {
	s := Render(box(1, 1, [3]float64{0.8, 0.8, 0.8}), Shot{Yaw: 45}, options())
	img := s.Image
	at := func(x, y int) color.NRGBA { return img.NRGBAAt(s.Pivot.X+x, s.Pivot.Y+y) }
	// One unit is a cell of 32: the box's top is 32 pixels above its foot,
	// its sides 32 to either side of its middle.
	top, left, right := at(0, -32), at(-16, -8), at(16, -8)
	if top.A != 255 || left.A != 255 || right.A != 255 {
		t.Fatalf("the box is not opaque: top %v left %v right %v", top, left, right)
	}
	if luma(top) <= luma(left) || luma(left) <= luma(right) {
		t.Errorf("not lit from the upper left: top %v, left %v, right %v", top, left, right)
	}
	if c := img.NRGBAAt(0, 0); c.A != 0 {
		t.Errorf("the corner is %v", c)
	}
	shadow := 0
	for y := s.Pivot.Y; y < img.Bounds().Dy(); y++ {
		for x := range img.Bounds().Dx() {
			c := img.NRGBAAt(x, y)
			if c.A > 40 && c.A < 160 && luma(c) == 0 {
				shadow++
			}
		}
	}
	if shadow < 50 {
		t.Errorf("%d pixels of shadow on the ground", shadow)
	}
}

// A one-sided wall facing glTF's +Z faces the viewer unturned, and turned
// half round shows its back, which is not drawn.
func TestYawTurnsAModel(t *testing.T) {
	wall := box(1, 1, [3]float64{1, 0.2, 0.2})
	front := wall.Meshes[0].Primitives[0]
	// Keep only the +Z face: its quad is the fourth.
	front.Indices = front.Indices[18:24]
	wall.Meshes[0].Primitives[0] = front
	wall.Materials[0].DoubleSided = false
	o := options()
	o.Shadow = 0
	count := func(s Sprite) int {
		n := 0
		b := s.Image.Bounds()
		for y := range b.Dy() {
			for x := range b.Dx() {
				if s.Image.NRGBAAt(x, y).A == 255 {
					n++
				}
			}
		}
		return n
	}
	if n := count(Render(wall, Shot{}, o)); n < 200 {
		t.Errorf("facing the viewer, %d pixels drawn", n)
	}
	if n := count(Render(wall, Shot{Yaw: 180}, o)); n != 0 {
		t.Errorf("turned away, %d pixels drawn", n)
	}
}

// Only one root of a model of variants is baked, stood on its own foot.
func TestOneVariantIsBakedOnItsOwnFoot(t *testing.T) {
	m := box(1, 1, [3]float64{0.5, 0.5, 0.5})
	m.Nodes = append(m.Nodes, gltf.Node{Mesh: 0, Skin: -1, Translation: [3]float64{5, 0, 0}, Rotation: [4]float64{0, 0, 0, 1}, Scale: [3]float64{1, 1, 1}})
	m.Roots = []int{0, 1}
	o := options()
	o.Shadow = 0
	both := Render(m, Shot{}, o)
	second := Render(m, Shot{Only: true, Root: 1}, o)
	first := Render(m, Shot{Only: true, Root: 0}, o)
	if second.Image.Bounds() != first.Image.Bounds() || second.Pivot != first.Pivot {
		t.Errorf("the second box, alone, is %v at %v, not as the first, %v at %v", second.Image.Bounds(), second.Pivot, first.Image.Bounds(), first.Pivot)
	}
	if both.Image.Bounds().Dx() <= first.Image.Bounds().Dx() {
		t.Error("both boxes are no wider than one")
	}
}
