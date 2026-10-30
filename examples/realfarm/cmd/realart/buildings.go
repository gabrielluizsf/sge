package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/gabrielluizsf/sge/bake"
	"github.com/gabrielluizsf/sge/building"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
)

// batch is what to bake: the manifest, the directory to write to, and
// the things to bake, by name, commas between, or "" for all.
type batch struct {
	manifest  kit.Manifest
	out, only string
}

// bakeBuildings bakes every house and thing package building builds,
// clad in the manifest's materials, at every heading, as j says.
func bakeBuildings(j batch, sprites map[string]sprite) {
	palette, err := paletteOf(j.manifest)
	if err != nil {
		fail(err)
	}
	out := j.out
	for _, k := range building.All() {
		if j.only != "" && !slices.Contains(strings.Split(j.only, ","), k.Name) {
			continue
		}
		model := k.Model(palette)
		// Its headings are baked at once, a core each.
		var wg sync.WaitGroup
		var mu sync.Mutex
		for i, yaw := range headings {
			wg.Add(1)
			go func() {
				defer wg.Done()
				s := bake.Render(model, bake.Shot{Yaw: yaw}, options())
				name := fmt.Sprintf("building_%s_%d", k.Name, i)
				if err := writePNG(filepath.Join(out, name+".png"), s); err != nil {
					fail(err)
				}
				b := s.Image.Bounds()
				mu.Lock()
				sprites[name] = sprite{W: b.Dx(), H: b.Dy(), PivotX: s.Pivot.X, PivotY: s.Pivot.Y}
				mu.Unlock()
			}()
		}
		wg.Wait()
		fmt.Println(k.Name, "built")
	}
}

// paletteOf is the buildings' palette: each material the manifest pins,
// read from its glTF, its colour, metal-rough and normal maps, at the
// span Poly Haven measured.
func paletteOf(m kit.Manifest) (building.Palette, error) {
	p := building.Plain()
	fe := kit.Fetcher()
	for _, a := range m.Assets {
		clad, ok := kit.Materials[a.ID]
		if a.Kind != "material" || !ok {
			continue
		}
		model, err := load(fe.Cache, a)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", a.ID, err)
		}
		if len(model.Materials) == 0 {
			return nil, fmt.Errorf("%s has no material", a.ID)
		}
		g := model.Materials[0]
		img := func(i int) image.Image {
			if i < 0 || i >= len(model.Images) {
				return nil
			}
			return asNRGBA(model.Images[i])
		}
		p[clad.Material] = building.Surface{Albedo: img(g.Texture), ARM: img(g.MetalRough), Normal: img(g.Normal), Span: a.Span, Across: clad.Across,
			Color: color.NRGBA{255, 255, 255, 255}, Roughness: 1, Metallic: 1}
	}
	return p, nil
}

// asNRGBA is img as an NRGBA image, as bake keeps between renders.
func asNRGBA(img image.Image) image.Image {
	if n, ok := img.(*image.NRGBA); ok {
		return n
	}
	n := image.NewNRGBA(img.Bounds())
	draw.Draw(n, n.Bounds(), img, img.Bounds().Min, draw.Src)
	return n
}
