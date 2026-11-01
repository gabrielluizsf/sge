package bakes

import (
	"fmt"
	"image"
	"image/color"

	"github.com/gabrielluizsf/sge/building"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
)

// Palette is the buildings' palette: each material the manifest pins,
// read from its glTF, its colour, metal-rough and normal maps, at the
// span Poly Haven measured. It is made once.
func (b *Baker) Palette() (building.Palette, error) {
	b.paletteOnce.Do(func() { b.palette, b.paletteErr = paletteOf(b.Manifest) })
	return b.palette, b.paletteErr
}

// paletteOf reads the palette m pins.
func paletteOf(m kit.Manifest) (building.Palette, error) {
	p := building.Plain()
	for _, a := range m.Assets {
		clad, ok := kit.Materials[a.ID]
		if a.Kind != "material" || !ok {
			continue
		}
		model, err := Load(a)
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
			return AsNRGBA(model.Images[i])
		}
		p[clad.Material] = building.Surface{Albedo: img(g.Texture), ARM: img(g.MetalRough), Normal: img(g.Normal), Span: a.Span, Across: clad.Across,
			Color: color.NRGBA{255, 255, 255, 255}, Roughness: 1, Metallic: 1}
	}
	return p, nil
}
