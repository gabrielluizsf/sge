package main

import (
	"path/filepath"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/farm/kit/bakes"
	"github.com/gabrielluizsf/sge/sprite"
)

// baked is a sprite read from the shelf: its picture and where it stands.
type baked struct {
	cv *canvas.Canvas
	at bakes.Sprite
}

// bake is the sprite called name, read once when it is first asked for;
// ok is false while the baker has not baked it.
func (a *realArt) bake(name string) (b *baked, ok bool) {
	if b, ok := a.seen[name]; ok {
		return b, true
	}
	at, ok := a.shelf.Sprite(name)
	if !ok {
		return nil, false
	}
	s, err := sprite.Load(filepath.Join(a.shelf.Dir, name+".png"))
	if err != nil {
		return nil, false
	}
	b = &baked{cv: s.Image, at: at}
	a.seen[name] = b
	return b, true
}

// picture is the image of the sprite called name, a blank frame while it
// is not baked.
func (a *realArt) picture(name string) *canvas.Canvas {
	if b, ok := a.bake(name); ok {
		return b.cv
	}
	return newFrame()
}

// turnOffset is how many facings the baker's first is turned from the
// farm's first, the right of the screen.
const turnOffset = 6

// first is the first of names that is baked.
func (a *realArt) first(names ...string) (*baked, bool) {
	for _, n := range names {
		if b, ok := a.bake(n); ok {
			return b, true
		}
	}
	return nil, false
}
