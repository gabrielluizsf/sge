package installer

import (
	"bytes"
	"errors"
	"image"
	"image/png"
)

// maxIconSide bounds the icon read: a PNG's header is read first, so a
// tiny file claiming a huge picture is refused before it is decoded.
const maxIconSide = 4096

// iconImage is the PNG data's picture, checked to be a sane size.
func iconImage(data []byte) (image.Image, error) {
	cfg, err := png.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	if cfg.Width < 1 || cfg.Height < 1 || cfg.Width > maxIconSide || cfg.Height > maxIconSide {
		return nil, errors.New("installer: the icon must be a PNG of at most 4096 pixels a side")
	}
	return png.Decode(bytes.NewReader(data))
}

// iconPNG is the icon as a square PNG of the largest of sides (largest
// first) no greater than its own, or of the smallest, the picture
// centred on a clear square and its colours averaged as it shrinks.
func iconPNG(data []byte, sides []int) ([]byte, int, error) {
	img, err := iconImage(data)
	if err != nil {
		return nil, 0, err
	}
	b := img.Bounds()
	own := max(b.Dx(), b.Dy())
	side := sides[len(sides)-1]
	for _, s := range sides {
		if s <= own {
			side = s
			break
		}
	}
	var out bytes.Buffer
	if err := png.Encode(&out, squared(img, side)); err != nil {
		return nil, 0, err
	}
	return out.Bytes(), side, nil
}
