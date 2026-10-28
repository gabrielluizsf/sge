package sprite

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
)

// Sheet cuts a sprite sheet into its frames, frameW by frameH pixels each,
// left to right and top to bottom.
func Sheet(img *canvas.Canvas, frameW, frameH int) ([]*Sprite, error) {
	if frameW <= 0 || frameH <= 0 || img.Width < frameW || img.Height < frameH {
		return nil, fmt.Errorf("sprite: a %dx%d sheet has no %dx%d frames", img.Width, img.Height, frameW, frameH)
	}
	var frames []*Sprite
	for y := 0; y+frameH <= img.Height; y += frameH {
		for x := 0; x+frameW <= img.Width; x += frameW {
			f, err := canvas.NewCanvas(frameW, frameH)
			if err != nil {
				return nil, err
			}
			for row := 0; row < frameH; row++ {
				copy(f.Pixels[row*f.Stride:row*f.Stride+frameW], img.Pixels[(y+row)*img.Stride+x:])
			}
			frames = append(frames, New(f))
		}
	}
	return frames, nil
}
