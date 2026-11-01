package bakes

import (
	"image"
	"image/draw"
	"os"
	"path"

	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
	"github.com/gabrielluizsf/sge/gltf"
)

// Load fetches asset a's files, checked, and reads its model.
func Load(a kit.Asset) (*gltf.Model, error) {
	fe := kit.Fetcher()
	for _, f := range a.Files {
		if _, err := fe.Get(f); err != nil {
			return nil, err
		}
	}
	root, err := os.OpenRoot(fe.Cache)
	if err != nil {
		return nil, err
	}
	defer root.Close()
	return gltf.Open(root, path.Clean(a.Model), gltf.DefaultLimits)
}

// AsNRGBA is img as an NRGBA image, as bake keeps between renders.
func AsNRGBA(img image.Image) image.Image {
	if n, ok := img.(*image.NRGBA); ok {
		return n
	}
	n := image.NewNRGBA(img.Bounds())
	draw.Draw(n, n.Bounds(), img, img.Bounds().Min, draw.Src)
	return n
}
