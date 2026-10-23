package sprite

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
)

// Load is a sprite of the PNG file at path.
func Load(path string) (*Sprite, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("sprite: %s: %w", path, err)
	}
	return FromImage(img)
}
