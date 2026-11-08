package water

import (
	"image"
	"image/color"
	"math"
)

// Config is how a normal map of waves is made.
type Config struct {
	// Size is the map's side in pixels.
	Size int
	// Seed makes one sea unlike another.
	Seed uint32
	// Span is how many cells of water the map is across, 1 if left 0; the
	// waves are the same size whatever it is, so a map of more cells is
	// less alike from one cell to the next.
	Span int
	// Steepness is how steep the waves are together; 1 is a breeze's
	// ripples.
	Steepness float64
	// Frame is how far into its loop the map is, of Frames: after Frames
	// frames each wave is where it began, so Frame Frames is Frame 0.
	Frame, Frames int
}

// Normals is a map of the normals of small waves, Span cells of water
// across: the surface's normal at every pixel, X along u, Y along v and
// Z up, each stored as 0 to 255 from -1 to 1. Its left edge follows its
// right, and its top its bottom, so copies of it meet without a seam, and
// its frames loop.
func Normals(c Config) *image.NRGBA {
	size := max(c.Size, 1)
	t := 0.0
	if c.Frames > 0 {
		t = float64(c.Frame%c.Frames) / float64(c.Frames)
	}
	ws := waves(c.Seed, c.Steepness, max(c.Span, 1))
	img := image.NewNRGBA(image.Rect(0, 0, size, size))
	for y := range size {
		for x := range size {
			sx, sy := slopeAt(ws, [2]float64{float64(x) / float64(size), float64(y) / float64(size)}, t)
			l := math.Sqrt(sx*sx + sy*sy + 1)
			img.SetNRGBA(x, y, color.NRGBA{enc(-sx / l), enc(-sy / l), enc(1 / l), 255})
		}
	}
	return img
}

// enc is a component from -1 to 1 as a byte.
func enc(c float64) uint8 { return uint8(math.Round((c*0.5 + 0.5) * 255)) }
