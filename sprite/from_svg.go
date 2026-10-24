package sprite

import (
	"fmt"

	"github.com/gabrielluizsf/antui/svg"
)

// FromSVG is a sprite of the SVG drawing src rendered width by height
// pixels, pivoting on its centre: pixel art drawn as SVG, with
// shape-rendering="crispEdges", stays sharp at any whole multiple of its
// size. The drawing is antui's, which loads the pictures it names, local
// files and URLs alike: give it only SVG you trust, such as your own art.
func FromSVG(src string, width, height int) (*Sprite, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("sprite: an SVG drawn %dx%d", width, height)
	}
	img, err := svg.Parse(src)
	if err != nil {
		return nil, fmt.Errorf("sprite: %w", err)
	}
	return New(img.Render(width, height)), nil
}
