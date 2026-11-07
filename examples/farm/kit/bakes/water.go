package bakes

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/gabrielluizsf/sge/bake"
	"github.com/gabrielluizsf/sge/terrain"
	waves "github.com/gabrielluizsf/sge/water"
)

// waterFrames is how many frames each shore's ripple loops in.
const waterFrames = 4

// shores are the kinds of water baked and the ground each shore lies on:
// a lake's on grass, the sea's on sand.
var shores = []struct {
	preset func() waves.Preset
	ground string
}{
	{waves.Lake, "sparse_grass"},
	{waves.Sea, "coast_sand_01"},
}

// shoreTile is what a tile of shore is made of: the kind of water, the
// normals of its frame, the ground, and the sprite's pivot.
type shoreTile struct {
	preset  waves.Preset
	normals *image.NRGBA
	ground  *image.NRGBA
	pivot   image.Point
}

// waterTasks bake, for each kind of water, each frame of its ripple.
func (b *Baker) waterTasks() []Task {
	var out []Task
	for _, k := range shores {
		for f := range waterFrames {
			p := k.preset()
			out = append(out, Task{Name: fmt.Sprintf("water %s %d", p.Name, f), Kind: "water", ID: p.Name, Core: true, Run: func() error { return b.water(p.Name, k.ground, f) }})
		}
	}
	return out
}

// water bakes the shore of every mask of the kind named kind, in frame f
// and in each of the variants a world of two cells' repeat needs, over
// its ground: water_<kind>_shore<mask>_<col><row>_<frame>.
func (b *Baker) water(kind, ground string, f int) error {
	img, pivot, err := b.groundTile(ground)
	if err != nil {
		return err
	}
	var p waves.Preset
	for _, k := range shores {
		if q := k.preset(); q.Name == kind {
			p = q
		}
	}
	n := waves.Normals(waves.Config{Size: 128, Span: waves.Span, Seed: 1, Steepness: p.Steepness, Frame: f, Frames: waterFrames})
	st := shoreTile{preset: p, normals: n, ground: img, pivot: pivot}
	for _, mask := range terrain.Masks() {
		for v := range waves.Span * waves.Span {
			shore := waves.Shore{Mask: mask, Col: v % waves.Span, Row: v / waves.Span, Phase: float64(f) / waterFrames}
			name := fmt.Sprintf("water_%s_shore%d_%d%d_%d", kind, mask, shore.Col, shore.Row, f)
			if err := b.shore(name, st, shore); err != nil {
				return err
			}
		}
	}
	return nil
}

// shore makes the tile of shore and puts it on the shelf as name.
func (b *Baker) shore(name string, st shoreTile, shore waves.Shore) error {
	shore.Surface, shore.Normals, shore.Ground = &st.preset.Surface, st.normals, st.ground
	return b.Shelf.Put(name, bake.Sprite{Image: waves.Tile(shore), Pivot: st.pivot})
}

// groundTile is the first tile of ground id, as the shelf holds it, baked
// now if it does not hold it yet.
func (b *Baker) groundTile(id string) (*image.NRGBA, image.Point, error) {
	if !b.Shelf.Has(id + "_g0") {
		for _, a := range b.Manifest.Assets {
			if a.ID == id && a.Kind == "ground" {
				if err := b.ground(a); err != nil {
					return nil, image.Point{}, err
				}
			}
		}
	}
	sp, ok := b.Shelf.Sprite(id + "_g0")
	if !ok {
		return nil, image.Point{}, fmt.Errorf("no ground %q in the manifest", id)
	}
	img, err := readNRGBA(filepath.Join(b.Shelf.Dir, id+"_g0.png"))
	return img, image.Pt(sp.PivotX, sp.PivotY), err
}

// readNRGBA reads the PNG at name as an NRGBA image.
func readNRGBA(name string) (*image.NRGBA, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}
	return AsNRGBA(img).(*image.NRGBA), nil
}
