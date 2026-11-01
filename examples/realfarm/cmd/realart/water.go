package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/gabrielluizsf/sge/bake"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
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

// bakeWater bakes the shore of every mask, in every frame and in each of
// the four variants a world of two cells' repeat needs, for each kind of
// water, over its ground, as j says; each is a tile named water_<kind>_
// shore<mask>_<col><row>_<frame>, the variant's cell parities.
func bakeWater(j batch, sprites map[string]sprite) {
	for _, k := range shores {
		ground, pivot, err := groundTile(j, sprites, k.ground)
		if err != nil {
			fail(err)
		}
		p := k.preset()
		for f := range waterFrames {
			n := waves.Normals(waves.Config{Size: 128, Span: waves.Span, Seed: 1, Steepness: p.Steepness, Frame: f, Frames: waterFrames})
			st := shoreTile{preset: p, normals: n, ground: ground, pivot: pivot}
			for _, mask := range terrain.Masks() {
				for v := range waves.Span * waves.Span {
					shore := waves.Shore{Mask: mask, Col: v % waves.Span, Row: v / waves.Span, Phase: float64(f) / waterFrames}
					name := fmt.Sprintf("water_%s_shore%d_%d%d_%d", p.Name, mask, shore.Col, shore.Row, f)
					sprites[name] = writeShore(tileFile{j.out, name}, st, shore)
				}
			}
		}
		fmt.Println(p.Name, "flows")
	}
}

// tileFile is where a tile is written: the folder and its name there.
type tileFile struct{ out, name string }

// writeShore makes the tile of shore, writes it to f, and is where its
// sprite stands.
func writeShore(f tileFile, st shoreTile, shore waves.Shore) sprite {
	out, name := f.out, f.name
	shore.Surface, shore.Normals, shore.Ground = &st.preset.Surface, st.normals, st.ground
	tile := waves.Tile(shore)
	if err := writePNG(filepath.Join(out, name+".png"), bake.Sprite{Image: tile, Pivot: st.pivot}); err != nil {
		fail(err)
	}
	b := tile.Bounds()
	return sprite{W: b.Dx(), H: b.Dy(), PivotX: st.pivot.X, PivotY: st.pivot.Y}
}

// groundTile is the first tile of ground id, as realart baked it into the
// output directory, or baked now from the manifest if it has not.
func groundTile(j batch, sprites map[string]sprite, id string) (*image.NRGBA, image.Point, error) {
	if img, ok := readTile(filepath.Join(j.out, id+"_g0.png")); ok {
		s, found := sprites[id+"_g0"]
		if !found {
			s = sprite{PivotX: img.Bounds().Dx() / 2}
		}
		return img, image.Pt(s.PivotX, s.PivotY), nil
	}
	for _, a := range j.manifest.Assets {
		if a.ID != id || a.Kind != "ground" {
			continue
		}
		model, err := load(kit.Fetcher().Cache, a)
		if err != nil {
			return nil, image.Point{}, fmt.Errorf("%s: %w", id, err)
		}
		s := bake.Ground(model, patches[0], options())
		return asNRGBA(s.Image).(*image.NRGBA), s.Pivot, nil
	}
	return nil, image.Point{}, fmt.Errorf("no ground %q in the manifest", id)
}

// readTile reads the PNG at name as an NRGBA image.
func readTile(name string) (*image.NRGBA, bool) {
	f, err := os.Open(name)
	if err != nil {
		return nil, false
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		return nil, false
	}
	return asNRGBA(img).(*image.NRGBA), true
}
