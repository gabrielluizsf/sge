package building

import (
	"math"
	"testing"

	"github.com/gabrielluizsf/sge/gltf"
)

// bounds are the least and greatest corners of m, or of its primitive of
// material name alone, if name is not "".
func bounds(m *gltf.Model, name string) (vec, vec) {
	lo, hi := vec{math.Inf(1), math.Inf(1), math.Inf(1)}, vec{math.Inf(-1), math.Inf(-1), math.Inf(-1)}
	for _, p := range m.Meshes[0].Primitives {
		if name != "" && m.Materials[p.Material].Name != name {
			continue
		}
		for _, q := range p.Positions {
			for i := range 3 {
				lo[i], hi[i] = min(lo[i], float64(q[i])), max(hi[i], float64(q[i]))
			}
		}
	}
	return lo, hi
}

// Each building is of a real size: houses 6 to 8 m wide, their roofs
// overhanging a little, a storey under 3 m; things the size they are.
func TestBuildingsAreLifeSize(t *testing.T) {
	sizes := map[string][2]float64{ // the least and most height
		"farmhouse": {5, 9}, "cottage": {5, 9}, "cabin": {4, 8}, "townhouse": {8, 13}, "brickhouse": {8, 13}, "sobrado": {7, 11},
		"stall": {2.4, 2.8}, "smithy": {4.5, 5.2}, "anvil": {0.75, 0.9}, "well": {2.4, 3},
		"fence_straight": {1, 1.1}, "fence_corner": {1, 1.1}, "fence_end": {1, 1.1}, "bench": {0.4, 0.5}, "wheelbarrow": {0.6, 0.85},
	}
	for _, k := range All() {
		lo, hi := bounds(k.Model(Plain()), "")
		want, ok := sizes[k.Name]
		if !ok {
			t.Errorf("%s: no size to check", k.Name)
			continue
		}
		if h := hi[1] - lo[1]; lo[1] < -1e-6 || h < want[0] || h > want[1] {
			t.Errorf("%s: from %.2f, %.2f m tall, not %v", k.Name, lo[1], h, want)
		}
		w, d := hi[0]-lo[0], hi[2]-lo[2]
		if w < k.Width-0.01 || w > k.Width+2.5 || d < k.Depth-0.01 || d > k.Depth+5 {
			t.Errorf("%s: %.2f by %.2f m, its footprint %v by %v", k.Name, w, d, k.Width, k.Depth)
		}
	}
	for _, k := range Houses {
		if k.Width < 6 || k.Width > 8 {
			t.Errorf("%s is %v m wide", k.Name, k.Width)
		}
	}
}

// A door is a door's size: 2.1 m tall, under a metre wide, its foot at
// the floor; a window 1.2 m tall, its sill 0.9 m above the floor.
func TestDoorsAndWindowsAreLifeSize(t *testing.T) {
	b := newBuilder(Plain())
	r := &raising{b: b, h: House{Width: 8, Depth: 6, Storeys: 1, Storey: 3, Wall: Plaster, Door: white}, eaves: 3.35}
	front := sides(8, 6)[0]
	r.door(front, 0)
	lo, hi := bounds(b.done(), "planks")
	if h, w := hi[1]-lo[1], hi[0]-lo[0]; math.Abs(h-2.1) > 1e-3 || w < 0.8 || w > 1 || math.Abs(lo[1]-floor) > 1e-3 {
		t.Errorf("the door is %.2f by %.2f m, its foot at %.2f", w, h, lo[1])
	}
	b = newBuilder(Plain())
	r.b = b
	r.window(front, vec{0, floor + sillAt})
	lo, hi = bounds(b.done(), "glass")
	if h := hi[1] - lo[1]; math.Abs(h-1.2) > 1e-3 || math.Abs(lo[1]-floor-0.9) > 1e-3 {
		t.Errorf("the window is %.2f m tall, its sill at %.2f", h, lo[1])
	}
}

// A texture repeats every span metres of a face, whatever its size.
func TestTexturesAreTheirSize(t *testing.T) {
	p := textured()
	b := newBuilder(p)
	b.block(box{lo: vec{-2, 0, -1}, hi: vec{2, 3, 1}, paint: paint{of: Brick}})
	prim := b.done().Meshes[0].Primitives[0]
	span := p[Brick].Span
	for k := 0; k < len(prim.Positions); k += 4 { // a face's four corners
		for i := 1; i < 4; i++ {
			dp := length(sub(f64(prim.Positions[k+i]), f64(prim.Positions[k])))
			du := math.Hypot(float64(prim.UVs[k+i][0]-prim.UVs[k][0]), float64(prim.UVs[k+i][1]-prim.UVs[k][1]))
			if math.Abs(du*span-dp) > 1e-4 {
				t.Fatalf("%.2f m of face is %.3f repeats of a %.2f m texture", dp, du, span)
			}
		}
	}
}

// A building's materials are its palette's: textured with its images,
// shared; a painted part its paint's colour, its relief kept; plain, the
// plain colours, but for the awning's painted stripes.
func TestMaterialsAreThePalettes(t *testing.T) {
	p := textured()
	m := Farmhouse.Model(p)
	for _, mt := range m.Materials {
		for _, i := range []int{mt.Texture, mt.MetalRough, mt.Occlusion, mt.Normal} {
			if i >= len(m.Images) {
				t.Fatalf("%s: image %d of %d", mt.Name, i, len(m.Images))
			}
		}
		if mt.MetalRough < 0 || mt.Normal < 0 {
			t.Errorf("%s has no maps", mt.Name)
		}
	}
	if len(m.Images) > 3*len(m.Materials) {
		t.Errorf("%d images for %d materials", len(m.Images), len(m.Materials))
	}
	painted := 0
	for _, mt := range m.Materials {
		if mt.Texture < 0 {
			painted++
		}
	}
	if painted == 0 {
		t.Error("the farmhouse's blue shutters are not painted")
	}
	plain := Stall.Model(Plain())
	for _, mt := range plain.Materials {
		if (mt.Texture >= 0) != (mt.Name == "canvas") {
			t.Errorf("plain %s: texture %d", mt.Name, mt.Texture)
		}
	}
}

// A texture whose grain runs across its image lays it along a beam's
// length; one whose grain runs down it, down it.
func TestGrainRunsAlongBeams(t *testing.T) {
	for _, across := range []bool{false, true} {
		p := textured()
		s := p[Timber]
		s.Across = across
		p[Timber] = s
		b := newBuilder(p)
		b.lay(beam{a: vec{0, 1, 0}, b: vec{3, 1, 0}, w: 0.1, h: 0.1, paint: paint{of: Timber}})
		prim := b.done().Meshes[0].Primitives[0]
		var du, dv float64
		for k := range prim.Positions {
			du = max(du, float64(prim.UVs[k][0]))
			dv = max(dv, float64(prim.UVs[k][1]))
		}
		if along := 3 / s.Span; (across && du < along*0.99) || (!across && dv < along*0.99) {
			t.Errorf("across %v: UVs reach %v, %v", across, du, dv)
		}
	}
}

// Every building has a name of its own.
func TestNamesAreUnique(t *testing.T) {
	seen := map[string]bool{}
	for _, k := range All() {
		if k.Name == "" || seen[k.Name] {
			t.Errorf("the name %q again", k.Name)
		}
		seen[k.Name] = true
	}
}
