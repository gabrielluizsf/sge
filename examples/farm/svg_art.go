package farm

import "github.com/gabrielluizsf/antui/canvas"

// svgArt is the farm's stylized look, every picture drawn from SVG: the
// tiles, icons, buildings and the player from the embedded assets, and
// the villagers and the HUD's text and panels as package art writes them.
type svgArt struct {
	// hud holds the HUD's text and panels by what they show, so the SVG
	// of each is written and rendered once, not every frame.
	hud map[string]*canvas.Canvas
}

// SVGArt is the farm's stylized look, drawn from SVG, which the farm
// command plays in.
func SVGArt() Art { return &svgArt{hud: map[string]*canvas.Canvas{}} }

// Tile is the tile name's picture at the game's scale: a cell's diamond
// with room above it.
func (a *svgArt) Tile(name string) *canvas.Canvas {
	return picture(spriteSVG(name), isoW/scale, isoH/scale)
}

// Icon is the 16 by 16 sprite name at the game's scale.
func (a *svgArt) Icon(name string) *canvas.Canvas {
	return picture(spriteSVG(name), 16, 16)
}

// Cursor is the outline of a cell's diamond.
func (a *svgArt) Cursor() *canvas.Canvas {
	return picture(spriteSVG("cursor"), isoW/scale, tile/scale)
}
