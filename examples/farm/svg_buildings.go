package farm

import (
	"github.com/gabrielluizsf/sge/examples/farm/art"
	"github.com/gabrielluizsf/sge/sprite"
)

// Building is the stall, its art drawn at half its size so that it stands
// on one cell among the people and the trees, or a house, at its full
// size on its two by two cells.
func (a *svgArt) Building(name string) *sprite.Sprite {
	if name == "stall" {
		return &sprite.Sprite{Image: picture(spriteSVG(name), art.StallW/2, art.StallH/2), PivotX: art.StallPivot[0], PivotY: art.StallPivot[1]}
	}
	return &sprite.Sprite{Image: picture(spriteSVG(name), art.HouseW, art.HouseH), PivotX: art.HousePivot[0], PivotY: art.HousePivot[1]}
}

// Houses is how many kinds of house the assets draw.
func (a *svgArt) Houses() int { return art.Houses }
