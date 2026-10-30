package gltf

import (
	"bytes"
	"image"
	_ "image/jpeg" // the image formats glTF allows
	_ "image/png"
)

// images are the model's images, decoded: each read from its buffer view
// or file, its size checked against the limit before any pixel is.
func images(doc *document, v views, src source) ([]image.Image, error) {
	var out []image.Image
	for i, im := range doc.Images {
		var data []byte
		var err error
		switch {
		case im.BufferView != nil:
			data, _, err = v.view(*im.BufferView)
		case im.URI != "":
			data, err = src.read(im.URI)
		default:
			err = malformed("image %d has no data", i)
		}
		if err != nil {
			return nil, err
		}
		img, err := DecodeImage(data, src.limits.Pixels)
		if err != nil {
			return nil, malformed("image %d: %v", i, err)
		}
		out = append(out, img)
	}
	return out, nil
}

// DecodeImage decodes a PNG or JPEG of at most pixels pixels, refusing a
// larger one from its header, before allocating for it.
func DecodeImage(data []byte, pixels int64) (image.Image, error) {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	if n, ok := mulInt(int64(cfg.Width), int64(cfg.Height)); !ok || n > pixels || cfg.Width <= 0 || cfg.Height <= 0 {
		return nil, malformed("image of %dx%d pixels", cfg.Width, cfg.Height)
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	return img, err
}
