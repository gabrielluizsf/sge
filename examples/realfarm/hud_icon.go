package main

import (
	"path/filepath"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/farm"
	"github.com/gabrielluizsf/sge/sprite"
)

// iconPx is the side of an icon on screen: 16 pixels of art.
const iconPx = 16 * farm.Scale

// Icon is the HUD's icon name, drawn from what was baked from the models
// of the things: the tools, the wood, the coin and the crops, picked and
// as seeds; or, for what no model is of, the HUD's own: a heart, and the
// marker over someone who wants something. While a model is not baked its
// place is a faint ring, never another game's picture.
func (a *realArt) Icon(name string) *canvas.Canvas {
	key := "i" + name
	if cv, ok := a.hud[key]; ok {
		return cv
	}
	cv := a.iconOf(name)
	if cv == nil {
		return a.hudPicture(drawing{"unbaked", unbakedSVG, iconPx, iconPx}) // not kept: the baker may yet bake it
	}
	a.hud[key] = cv
	return cv
}

// iconOf is the picture of icon name, or nil if it is a baked one not
// baked yet.
func (a *realArt) iconOf(name string) *canvas.Canvas {
	switch name {
	case "icon_heart":
		return a.hudPicture(drawing{name, heartSVG(true), iconPx, iconPx})
	case "icon_heart_empty":
		return a.hudPicture(drawing{name, heartSVG(false), iconPx, iconPx})
	case "icon_food":
		return a.hudPicture(drawing{name, foodSVG, iconPx, iconPx})
	case "icon_bubble":
		return a.hudPicture(drawing{name, markerSVG, iconPx, iconPx})
	}
	if !a.shelf.Has(name) {
		return nil
	}
	s, err := sprite.Load(filepath.Join(a.shelf.Dir, name+".png"))
	if err != nil {
		return nil
	}
	return shrink(s.Image, iconPx, iconPx)
}

// drawing is an SVG src to render w by h pixels, kept as name.
type drawing struct {
	name, src string
	w, h      int
}

// hudPicture is drawing d, rendered once.
func (a *realArt) hudPicture(d drawing) *canvas.Canvas {
	key := "svg" + d.name
	if cv, ok := a.hud[key]; ok {
		return cv
	}
	s, err := sprite.FromSVG(d.src, d.w, d.h)
	if err != nil {
		panic(err)
	}
	a.hud[key] = s.Image
	return s.Image
}
