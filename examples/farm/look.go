package main

import "github.com/gabrielluizsf/antui/canvas"

// look is the 16 by 16 sprite name at the game's scale.
func look(name string) *canvas.Canvas {
	return picture(spriteSVG(name), 16, 16)
}
