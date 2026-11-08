package bake

import (
	"image"
	"testing"
)

// A person is baked many times, in many poses: their textures, 1024 by
// 1024, are read for each.
func BenchmarkRenderTextured(b *testing.B) {
	img := image.NewNRGBA(image.Rect(0, 0, 1024, 1024))
	for i := range img.Pix {
		img.Pix[i] = uint8(i * 7)
	}
	m := box(1, 1, [3]float64{1, 1, 1})
	m.Images = []image.Image{img}
	m.Materials[0].Texture = 0
	p := &m.Meshes[0].Primitives[0]
	p.UVs = make([][2]float32, len(p.Positions))
	o := options()
	b.ResetTimer()
	for range b.N {
		Render(m, Shot{}, o)
	}
}

// A texture is read once for many renders, and the oldest let go when
// more than keepMost are kept; one read for colour is not one read as
// data.
func TestTexturesAreKeptBetweenRenders(t *testing.T) {
	imgs := make([]*image.NRGBA, keepMost+1)
	for i := range imgs {
		imgs[i] = image.NewNRGBA(image.Rect(0, 0, 2, 2))
	}
	first := keptTexture(imgs[0], true)
	if keptTexture(imgs[0], true) != first {
		t.Error("read again")
	}
	if keptTexture(imgs[0], false) == first {
		t.Error("its colour read as its data")
	}
	for _, img := range imgs[1:] {
		keptTexture(img, true)
	}
	if keptTexture(imgs[0], true) == first {
		t.Error("kept more than keepMost")
	}
	// Two pixels of sRGB grey 188 are half the light, as the curve says.
	g := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	g.Pix = []uint8{188, 188, 188, 255}
	if c := textureOf(g, true).px[0]; c[0] < 0.49 || c[0] > 0.51 || c[3] != 1 {
		t.Errorf("sRGB 188 is %v", c)
	}
}
