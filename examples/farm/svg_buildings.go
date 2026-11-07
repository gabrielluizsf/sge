package farm

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/farm/art"
	"github.com/gabrielluizsf/sge/sprite"
)

// Building is the stall, its art drawn at half its size so that it stands
// on one cell among the people and the trees, or a house, at its full
// size on its two by two cells.
func (a *svgArt) Building(name string) *sprite.Sprite {
	if name == "well" { // this look has none: a blank picture, for a look that sizes one
		blank, _ := canvas.NewCanvas(1, 1)
		return &sprite.Sprite{Image: blank}
	}
	if art.IsShed(name) {
		return &sprite.Sprite{Image: picture(spriteSVG(name), art.StallW/2, art.StallH/2), PivotX: art.StallPivot[0], PivotY: art.StallPivot[1]}
	}
	if name == "stall" {
		return &sprite.Sprite{Image: picture(spriteSVG(name), art.StallW/2, art.StallH/2), PivotX: art.StallPivot[0], PivotY: art.StallPivot[1]}
	}
	return &sprite.Sprite{Image: picture(spriteSVG(name), art.HouseW, art.HouseH), PivotX: art.HousePivot[0], PivotY: art.HousePivot[1]}
}

// Houses is how many kinds of house the assets draw.
func (a *svgArt) Houses() int { return art.Houses }

// Footprint is one cell for the stall, none for a well, which this look
// does not draw, and two by two for a house.
func (a *svgArt) Footprint(name string) (cols, rows int) {
	switch {
	case art.IsShed(name):
		return 1, 1
	}
	switch name {
	case "stall":
		return 1, 1
	case "well":
		return 0, 0
	}
	return 2, 2
}

// Reach is nothing: every picture is the tile's size.
func (a *svgArt) Reach() int { return 0 }
