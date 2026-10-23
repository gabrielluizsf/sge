package game

import "github.com/gabrielluizsf/antui/canvas"

// New is an empty world on a black background.
func New() *Game {
	return &Game{Background: canvas.RGB(0, 0, 0)}
}
