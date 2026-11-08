package bake

import (
	"image"
	"sync"
)

// keepMost is how many textures read are kept for the next render: a
// model is baked in many poses and headings, its textures read once.
const keepMost = 8

// keptKey is an image, by its pointer, read for colour or as data.
type keptKey struct {
	img    *image.NRGBA
	colour bool
}

// kept are the textures last read, the oldest first. A model's images
// are taken not to change between renders.
var kept struct {
	sync.Mutex
	order []keptKey
	by    map[keptKey]*texture
}

// keptTexture is img read as textureOf reads it, kept for later renders
// if it is an NRGBA image, as decoders and packages make them.
func keptTexture(img image.Image, colour bool) *texture {
	n, ok := img.(*image.NRGBA)
	if !ok {
		return textureOf(img, colour)
	}
	k := keptKey{n, colour}
	kept.Lock()
	defer kept.Unlock()
	if t := kept.by[k]; t != nil {
		return t
	}
	if kept.by == nil {
		kept.by = map[keptKey]*texture{}
	}
	t := textureOf(img, colour)
	kept.by[k] = t
	kept.order = append(kept.order, k)
	if len(kept.order) > keepMost {
		delete(kept.by, kept.order[0])
		kept.order = kept.order[1:]
	}
	return t
}
