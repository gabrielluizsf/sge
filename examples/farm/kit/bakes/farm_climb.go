package bakes

import (
	"image"

	"github.com/gabrielluizsf/sge/bake"
	"github.com/gabrielluizsf/sge/gltf"
	"github.com/gabrielluizsf/sge/human"
)

// FramesOf is how many frames of a motion are baked: getting on a horse and
// down are long, and move far, so they have more.
func FramesOf(motion string) int {
	if climbing(motion) {
		return len(human.Mount.Poses)
	}
	return Frames
}

// climbing reports whether a motion is getting on a horse or down from it.
func climbing(motion string) bool {
	return motion == human.Mount.Name || motion == human.Dismount.Name
}

// climbShade is how strong the climber's own shadow is, u of the way through
// the motion: whole on the ground, none in the saddle, where the horse's is.
func climbShade(motion string, u float64) float64 {
	if motion == human.Dismount.Name {
		return span(u, 0.6, 0.86)
	}
	return 1 - span(u, 0.5, 0.85)
}

// span is where u is from a to b, held from 0 to 1.
func span(u, a, b float64) float64 { return min(max((u-a)/(b-a), 0), 1) }

// climber is the picture of someone getting on or down a horse as shot: a
// rider's, hidden where the horse's body would hide them, over their own
// shadow of a share of its strength.
func climber(m *gltf.Model, sh bake.Shot, shade float64) bake.Sprite {
	o := Options()
	o.Shadow *= shade
	body := only(m, false)
	return overShadow(rider(m, sh), shadowOf(bake.Render(body, sh, o), body, sh))
}

// shadowOf is what of shot is not the body: its shadow on the ground.
func shadowOf(shot bake.Sprite, body *gltf.Model, sh bake.Shot) bake.Sprite {
	o := Options()
	o.Shadow = 0
	alone := bake.Render(body, sh, o)
	out := *shot.Image
	out.Pix = append([]uint8(nil), shot.Image.Pix...)
	b := out.Rect
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if bareAt(alone, shot, x, y).A > 0 {
				out.SetNRGBA(x, y, colorNRGBA{})
			}
		}
	}
	return bake.Sprite{Image: &out, Pivot: shot.Pivot}
}

// overShadow is a picture drawn over a shadow, on one spot.
func overShadow(top, shadow bake.Sprite) bake.Sprite {
	box := top.Image.Rect.Sub(top.Pivot).Union(shadow.Image.Rect.Sub(shadow.Pivot))
	out := image.NewNRGBA(image.Rect(0, 0, box.Dx(), box.Dy()))
	pivot := image.Pt(-box.Min.X, -box.Min.Y)
	for _, layer := range []bake.Sprite{shadow, top} {
		b := layer.Image.Rect
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				if c := layer.Image.NRGBAAt(x, y); c.A > 0 {
					at := image.Pt(x, y).Sub(layer.Pivot).Add(pivot)
					out.SetNRGBA(at.X, at.Y, c)
				}
			}
		}
	}
	return bake.Sprite{Image: out, Pivot: pivot}
}
